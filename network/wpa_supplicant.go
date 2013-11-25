package network

import (
    "errors"
    "log"
    "strings"
    "strconv"
    "time"
    "github.com/r3boot/rlib/sys"
)

type WirelessNetwork struct {
    Id          int
    Ssid        string
    Frequency   int
    Signal      int
    Connected   bool
}

type WpaSupplicant struct {
    Interface           string
    CmdWpaSupplicant    string
    CmdWpaCli           string
    CfgFile             string
    Driver              string
}

/*
 * Run a command through wpa_cli. Return an error if the return-code from
 * wpa_cli is non-zero.
 */
func (w WpaSupplicant) Run (command string) (stdout, stderr []string, err error) {
    stdout, stderr, err = sys.Run(w.CmdWpaCli, "-i" + w.Interface, command)
    return
}

/*
 * Start wpa_supplicant. Check for the existence of the configuration file. If
 * it does not exist, return an error. If the returncode of the wpa_supplicant
 * is non-zero, return an error.
 */
func (w WpaSupplicant) Start () (err error) {
    _, _, err = sys.Run(w.CmdWpaSupplicant, "-B", "-D" + w.Driver, "-i" +
        w.Interface, "-c", w.CfgFile)

    if err != nil {
        log.Print(err)
    }

    return
}

/*
 * Stop wpa_supplicant. Throw an error if the returncode from wpa_cli is
 * non-zero.
 */
func (w WpaSupplicant) Stop () (err error) {
    _, _, err = w.Run("terminate")
    return
}

/*
 * Ping the wpa_supplicant processs to see if it's reachable. Return true if
 * it is, false otherwise.
 */
func (w WpaSupplicant) IsRunning () bool {
    _, _, err := w.Run("ping")
    return err == nil
}

/*
 * Check to see if wpa_supplicant is connected to a network. If the return-code
 * of wpa_cli is non-zero, or the wpa_state is not equal to "COMPLETED", return
 * false, else return true.
 */
func (w WpaSupplicant) IsConnected () bool {
    stdout, _, err := w.Run("status")
    if err != nil {
        return false
    }

    for _, line := range stdout {
        t := strings.Split(line, "=")
        if (t[0] == "wpa_state") && (t[1] == WPA_STATE_COMPLETED) {
            return true
        }
    }
    return false
}

func (w WpaSupplicant) Scan () {
    _, _, err := w.Run("scan")
    if err != nil {
        return
    }
    time.Sleep(WPA_SCAN_INTERVAL * time.Second)
}

/*
 * Process the scan results from wpa_cli scan and return an array
 * of WirelessNetwork
 */
func (w WpaSupplicant) AvailableNetworks () (nets []WirelessNetwork) {
    myname := "WpaSupplicant.AvailableNetworks"
    Log.Debug(myname, "Scanning for available wireless networks")

    stdout, _, err := w.Run("scan_results")
    if err != nil {
        Log.Fatal(myname, err.Error())
    }

    for _, line := range stdout {
        if strings.HasPrefix(line, "bssid") { continue }
        if len(line) == 0 { continue }

        t := strings.Split(line, "\t")

        var network = new(WirelessNetwork)
        network.Ssid = t[4]
        network.Frequency, err = strconv.Atoi(t[1])
        if err != nil {
            Log.Warning(myname, "Failed to parse frequency for " + t[4])
            continue
        }
        network.Signal, err = strconv.Atoi(t[2])
        if err != nil {
            Log.Warning(myname, "Failed to parse signal for " + t[4])
            continue
        }

        Log.Debug(myname, "Found wireless network (ssid: " + t[4] + ", freq: " + t[1] + "KHz, strengh: " + t[2] + "dB)")
        nets = append(nets, *network)
    }

    return
}

/*
 * Parse the results from wpa_cli list_networks and return a list of
 * configured wireless networks as WirelessNetwork. If the network is
 * currently established, set the Connected flag of the respective network
 * to true.
 */
func (w WpaSupplicant) ConfiguredNetworks () (nets []WirelessNetwork) {
    myname := "WpaSupplicant.ConfiguredNetworks"

    Log.Debug(myname, "Looking up configured wireless networks")

    stdout, _, err := w.Run("list_networks")
    if err != nil {
        log.Fatal(err)
    }

    for _, line := range stdout {
        if strings.HasPrefix(line, "network id") { continue }
        if len(line) == 0 { continue }

        t := strings.Split(line, "\t")

        var network = new(WirelessNetwork)
        network.Id, err = strconv.Atoi(t[0])
        if err != nil {
            Log.Warning(myname, "Failed to parse network id for " + t[1] + ", skipping")
            continue
        }
        network.Ssid = t[1]
        if t[3] == "[CURRENT]" {
            network.Connected = true
        }

        Log.Debug(myname, "Found configuration for " + network.Ssid)

        nets = append(nets, *network)
    }

    return
}

/*
 * Look at the cross-section between the available and configured wireless
 * networks, and return a list of WirelessNetwork of usable networks.
 */
func (w WpaSupplicant) GetUsableNetworks () (nets []WirelessNetwork) {
    myname := "WpaSupplicant.GetUsableNetworks"
    Log.Debug(myname, "Creating list of usable wireless networks")

    available_networks := w.AvailableNetworks()
    configured_networks := w.ConfiguredNetworks()

    for _, available := range available_networks {
        for _, configured := range configured_networks {
            if available.Ssid == configured.Ssid {
                available.Id = configured.Id
                available.Connected = configured.Connected
                nets = append(nets, available)
                break
            }
        }
    }

    return
}

/*
 * Connect to an ssid. Return true if wireless is established, false otherwise
 */
func (w WpaSupplicant) Connect(id int) bool {
    myname := "WpaSupplicant.Connect"

    Log.Debug(myname, "Connecting to wireless network")
    var timeout_counter int = 0

    if w.IsConnected() {
        return true
    }

    _, _, err := w.Run("select_network " + strconv.Itoa(id))
    if err != nil {
        return false
    }

    for {
        if timeout_counter >= WPA_CONNECT_TIMEOUT {
            return false
        }

        if w.IsConnected() {
            return true
        }

        timeout_counter += 1
        time.Sleep(WPA_CONNECT_INTERVAL * time.Millisecond)
    }

    return false
}

/*
 * Connect to the first usable wireless network. Return true if wireless
 * is established, false otherwise.
 */
func (w WpaSupplicant) ConnectAny () bool {
    if w.IsConnected() {
        return true
    }

    wireless_networks := w.GetUsableNetworks()
    for _, network := range wireless_networks {
        if w.Connect(network.Id) {
            return true
        }
    }

    return false
}

func WpaSupplicantFactory (intf string) (w WpaSupplicant, err error) {
    wpa_supplicant, err := sys.BinaryPrefix("wpa_supplicant")
    if err != nil {
        return
    }

    wpa_cli, err := sys.BinaryPrefix("wpa_cli")
    if err != nil {
        return
    }

    cfg_file, err := sys.ConfigPrefix("wpa_supplicant-" + intf + ".conf")
    if err != nil {
        return
    }

    uname, err := sys.Uname()
    if err != nil {
        return
    }

    var driver string
    if uname.Ident == "Linux" {
        driver = "wext"
    } else if uname.Ident == "FreeBSD" {
        driver = "bsd"
    }

    w = WpaSupplicant{intf, wpa_supplicant, wpa_cli, cfg_file, driver}

    return
}
