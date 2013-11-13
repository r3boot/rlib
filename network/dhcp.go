package network

import (
    "errors"
    "net"
    "io/ioutil"
    "os"
    "strings"
    "github.com/r3boot/rlib/sys"
)

type Dhcpcd struct {
    Interface string
}

/*
 * Start dhcpcd on intf. Print an error if anything goes wrong
 */
func (d Dhcpcd) Start () {
    myname := "Dhcpcd.Start"
    _, _, err := sys.Run("/usr/sbin/dhcpcd", "-q", d.Interface)
    if err != nil {
        Log.Warning(myname, "Failed to start dhcpcd")
    }
}

/*
 * Stop dhcpcd on intf. Print an error if anything goes wrong
 */
func (d Dhcpcd) Stop () {
    myname := "Dhcpcd.Stop"
    _, _, err := sys.Run("/usr/sbin/dhcpcd", "-x", d.Interface)
    if err != nil {
        Log.Warning(myname, "Failed to stop dhcpcd")
    }
}

/*
 * See if dhcpcd is running. First, check for the existence of
 * /run/dhcpcd-<nic>.pid. If this file does not exist, return false. Then,
 * read the content of /proc/<pid>/cmdline and do a string match on both
 * "dhcpcd" and intf.Name. If both match, return true. All other results
 * will return false.
 */
func (d Dhcpcd) IsRunning () bool {

    pid_file := "/run/dhcpcd-" + d.Interface + ".pid"
    _, err := os.Stat(pid_file)
    if err != nil {
        return false
    }

    content, err := ioutil.ReadFile(pid_file)
    if err != nil {
        return false
    }
    pid := string(content)

    proc_file := "/proc/" + pid + "/cmdline"
    content, err = ioutil.ReadFile(proc_file)
    if err != nil {
        return false
    }
    ps := string(content)

    return strings.Contains(ps, "dhcpcd") &&
           strings.Contains(ps, d.Interface)
}

func (d Dhcpcd) GetOffer () (ip net.IP, network net.IPNet, err error) {
    stdout, _, err := sys.Run("/usr/sbin/dhcpcd", "-4", "-T", d.Interface)
    if err != nil {
        return
    }

    var raw_ip, cidr_mask string
    for _, line := range stdout {
        t := strings.Split(line, "=")
        if t[0] == "new_ip_address" {
            raw_ip = t[1]
        } else if t[0] == "new_subnet_cidr" {
            cidr_mask = t[1]
        }
    }

    if len(raw_ip) == 0 {
        err = errors.New("Failed to obtain ip address")
        return
    } else if len(cidr_mask) == 0 {
        err = errors.New("Failed to obtain cidr mask")
        return
    }

    i, n, err := net.ParseCIDR(raw_ip + "/" + cidr_mask)
    if err == nil {
        ip = i
        network = *n
    }

    return
}
