package vpn

import (
    "net"
    "os"
    "time"
    "io/ioutil"
    "strconv"
    "strings"
    "github.com/r3boot/rlib/sys"
    "github.com/r3boot/rlib/network"
)

type OpenVPN struct {
    Name string
    InterfaceName string
    Interface network.Interface
    ConfigFile string
    PidFile string
    StatusFile string
    Remote net.IP
    Port int
}

func (ovpn *OpenVPN) Start () {
    myname := "OpenVPN.Start"
    o := *ovpn

    Log.Debug(myname, "Starting OpenVPN for " + o.Name)

    _, _, err := sys.Run("/usr/sbin/openvpn", "--daemon", "--config", o.ConfigFile, "--writepid", o.PidFile)
    if err != nil {
        Log.Warning(myname, "Failed to start OpenVPN")
        return
    }

    var timeout_counter = 0
    for {
        if timeout_counter >= OVPN_STARTUP_MAXWAIT {
            Log.Warning(myname, "Connection establishment took too long")
            o.Stop()
            break
        }

        if o.IsConnected() {
            break
        }

        time.Sleep(1 * time.Second)
        timeout_counter += 1
    }

    raw_intf, err := net.InterfaceByName(o.InterfaceName)
    intf := network.InterfaceFactory(*raw_intf)
    o.Interface = intf

    *ovpn = o
    Log.Debug(myname, "Connected to OpenVPN tunnel " + o.Name)
}

func (ovpn *OpenVPN) Stop () {
    myname := "OpenVPN.Stop"

    o := *ovpn
    Log.Debug(myname, "Stopping OpenVPN")

    if ! sys.FileExists(o.PidFile) {
        Log.Warning(myname, "Pidfile " + o.PidFile + " does not exist")
        return
    }

    pid := o.GetPid()
    proc, err := os.FindProcess(pid)
    if err != nil {
        Log.Warning(myname, "Failed to find running process for pid " + strconv.Itoa(pid))
        return
    }

    proc.Kill()
    os.Remove(o.PidFile)
    os.Remove(o.StatusFile)
    o.Interface = network.Interface{}

    *ovpn = o
}

func (ovpn *OpenVPN) GetPid () (pid int) {
    myname := "OpenVPN.GetPid"
    o := *ovpn

    if ! sys.FileExists(o.PidFile) {
        Log.Warning(myname, "Pidfile " + o.PidFile + " does not exist")
        return
    }

    content, err := ioutil.ReadFile(o.PidFile)
    if err != nil {
        Log.Warning(myname, "Cannot read pidfile " + o.PidFile)
        return
    }

    pid, err = strconv.Atoi(strings.Split(string(content), "\n")[0])
    if err != nil {
        Log.Warning(myname, "Failed to convert pid to int: " + string(content))
        pid = 0
        return
    }

    return
}

func (ovpn *OpenVPN) IsRunning () bool {
    o := *ovpn
    pid := o.GetPid()

    proc_file := "/proc/" + strconv.Itoa(pid) + "/cmdline"
    content, err := ioutil.ReadFile(proc_file)
    if err != nil {
        return false
    }
    ps := string(content)

    return strings.Contains(ps, "openvpn") &&
           strings.Contains(ps, o.Name)
}

func (ovpn *OpenVPN) IsConnected () bool {
    myname := "OpenVPN.IsConnected"

    o := *ovpn

    if ! sys.FileExists(o.StatusFile) {
        Log.Warning(myname, "Status file " + o.StatusFile + " does not exist")
        return false
    }

    content, err := ioutil.ReadFile(o.StatusFile)
    if err != nil {
        Log.Warning(myname, "Failed to read status file " + o.StatusFile)
        return false
    }

    var t_lastupdate time.Time
    for _, line := range strings.Split(string(content), "\n") {
        t := strings.Split(line, ",")

        if t[0] == OVPN_STATUS_UPDATE {
            t_lastupdate, err = time.Parse("Mon Jan 02 15:04:05 2006", t[1])
            if err != nil {
                Log.Warning(myname, "Failed to parse date")
                return false
            }
            break
        }
    }

    return time.Now().Sub(t_lastupdate) < (61 * time.Second)
}

func (ovpn *OpenVPN) ReadConfig () {
    myname := "OpenVPN.ReadConfig"
    o := *ovpn

    if ! sys.FileExists(o.ConfigFile) {
        Log.Fatal(myname, "Cannot find" + o.ConfigFile)
        return
    }

    content, err := ioutil.ReadFile(o.ConfigFile)
    if err != nil {
        Log.Warning(myname, "Cannot read " + o.ConfigFile)
        return
    }
    for _, line := range strings.Split(string(content), "\n") {
        t := strings.Split(line, " ")
        switch t[0] {
        case OVPN_CFG_REMOTE: {
            o.Remote, _, err = net.ParseCIDR(t[1] + "/32")
            if err != nil {
                Log.Warning(myname, "Failed to parse remote")
                continue
            }
        }
        case OVPN_CFG_PORT: {
            o.Port, err = strconv.Atoi(t[1])
            if err != nil {
                Log.Warning(myname, "Failed to parse port")
                continue
            }
        }
        case OVPN_CFG_DEVICE: {
            o.InterfaceName = t[1]
        }
        }
    }

    *ovpn = o
}

func OpenVPNFactory (name string) (o OpenVPN) {
    o = OpenVPN{
        name,
        "",
        network.Interface{},
        "/etc/openvpn/" + name + ".conf",
        "/var/run/openvpn-" + name + ".pid",
        "/var/run/openvpn-" + name + ".status",
        net.IP{},
        1900,
    }
    o.ReadConfig()

    return
}
