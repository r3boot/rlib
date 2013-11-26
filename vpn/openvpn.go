package vpn

import (
    "errors"
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
    Name            string
    InterfaceName   string
    Interface       network.Interface
    ConfigFile      string
    PidFile         string
    StatusFile      string
    Remote          net.IP
    Port            int
}

func (ovpn *OpenVPN) Start () (err error) {
    o := *ovpn

    _, _, err = sys.Run("/usr/sbin/openvpn", "--daemon", "--config", o.ConfigFile, "--writepid", o.PidFile)
    if err != nil {
        err = errors.New("Failed to start OpenVPN: " + err.Error())
        return
    }

    var timeout_counter = 0
    for {
        if timeout_counter >= OVPN_STARTUP_MAXWAIT {
            err = errors.New("Connection establishment took too long")
            o.Stop()
            break
        }

        is_connected, e := o.IsConnected()
        if e != nil {
            err = e
            return
        }

        if is_connected {
            break
        }

        time.Sleep(1 * time.Second)
        timeout_counter += 1
    }

    raw_intf, err := net.InterfaceByName(o.InterfaceName)
    intf, err := network.InterfaceFactory(*raw_intf)
    if err != nil {
        return
    }

    o.Interface = intf

    *ovpn = o

    return
}

func (ovpn *OpenVPN) Stop () (err error) {
    o := *ovpn

    if ! sys.FileExists(o.PidFile) {
        return
    }

    pid, err := o.GetPid()
    if err != nil {
        return
    }

    proc, err := os.FindProcess(pid)
    if err != nil {
        err = errors.New("Failed to find running process for pid " + strconv.Itoa(pid) + ": " + err.Error())
        return
    }

    proc.Kill()
    os.Remove(o.PidFile)
    os.Remove(o.StatusFile)
    o.Interface = network.Interface{}

    *ovpn = o

    return
}

func (ovpn *OpenVPN) GetPid () (pid int, err error) {
    o := *ovpn

    if ! sys.FileExists(o.PidFile) {
        return
    }

    content, err := ioutil.ReadFile(o.PidFile)
    if err != nil {
        err = errors.New("Cannot read pidfile " + o.PidFile + ": " + err.Error())
        return
    }

    pid, err = strconv.Atoi(strings.Split(string(content), "\n")[0])
    if err != nil {
        err = errors.New("Failed to convert pid to int: " + string(content) + ": " + err.Error())
        pid = 0
        return
    }

    return
}

func (ovpn *OpenVPN) IsRunning () (result bool, err error) {
    o := *ovpn
    pid, err := o.GetPid()
    if err != nil {
        return
    }

    proc_file := "/proc/" + strconv.Itoa(pid) + "/cmdline"
    content, err := ioutil.ReadFile(proc_file)
    if err != nil {
        return
    }
    ps := string(content)

    result = strings.Contains(ps, "openvpn") &&
           strings.Contains(ps, o.Name)

    return
}

func (ovpn *OpenVPN) IsConnected () (result bool, err error) {
    o := *ovpn

    if ! sys.FileExists(o.StatusFile) {
        err = errors.New("Status file " + o.StatusFile + " does not exist")
        return
    }

    content, err := ioutil.ReadFile(o.StatusFile)
    if err != nil {
        err = errors.New("Failed to read status file " + o.StatusFile + ": " + err.Error())
        return
    }

    var t_lastupdate time.Time
    for _, line := range strings.Split(string(content), "\n") {
        t := strings.Split(line, ",")

        if t[0] == OVPN_STATUS_UPDATE {
            t_lastupdate, err = time.Parse("Mon Jan 02 15:04:05 2006", t[1])
            if err != nil {
                err = errors.New("Failed to parse date: " + err.Error())
                return
            }
            break
        }
    }

    result = time.Now().Sub(t_lastupdate) < (61 * time.Second)
    return
}

func (ovpn *OpenVPN) ReadConfig () (err error) {
    o := *ovpn

    if ! sys.FileExists(o.ConfigFile) {
        err = errors.New("Cannot find" + o.ConfigFile)
        return
    }

    content, err := ioutil.ReadFile(o.ConfigFile)
    if err != nil {
        err = errors.New("Cannot read " + o.ConfigFile + ": " + err.Error())
        return
    }
    for _, line := range strings.Split(string(content), "\n") {
        t := strings.Split(line, " ")
        switch t[0] {
        case OVPN_CFG_REMOTE: {
            o.Remote, _, err = net.ParseCIDR(t[1] + "/32")
            if err != nil {
                err = errors.New("Failed to parse remote: " + err.Error())
                return
            }
        }
        case OVPN_CFG_PORT: {
            o.Port, err = strconv.Atoi(t[1])
            if err != nil {
                err = errors.New("Failed to parse port: " + err.Error())
                return
            }
        }
        case OVPN_CFG_DEVICE: {
            o.InterfaceName = t[1]
        }
        }
    }

    *ovpn = o

    return
}

func OpenVPNFactory (name string) (o OpenVPN, err error) {
    o = *new(OpenVPN)

    cfg_file, err := sys.ConfigPrefix("openvpn/" + name + ".conf")
    if err != nil {
        return
    }

    o.Name = name
    o.InterfaceName = ""
    o.Interface = network.Interface{}
    o.ConfigFile = cfg_file
    o.PidFile = "/var/run/openvpn-" + name + ".pid"
    o.StatusFile = "/var/run/openvpn-" + name + ".status"
    o.Remote = net.IP{}
    o.Port = 1900

    o.ReadConfig()

    return
}
