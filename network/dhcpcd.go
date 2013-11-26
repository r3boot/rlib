package network

import (
    "errors"
    "net"
    "io/ioutil"
    "os"
    "strings"
)

type Dhcpcd struct {
    Interface       string
    SetMtu          bool
    UseResolvconf   bool
    CmdDhcpcd       string
}

/*
 * Start dhcpcd on intf. Print an error if anything goes wrong
 */
func (d Dhcpcd) Start () (err error) {
    var args []string
    args = append(args, "-q")

    if ! d.UseResolvconf {
        args = append(args, "-C")
        args = append(args, "resolv.conf")
    }

    if ! d.SetMtu {
        args = append(args, "-C")
        args = append(args, "mtu")
    }

    args = append(args, d.Interface)

    _, _, err := sys.Run(d.CmdDhcpcd, args...)
    if err != nil {
        err = errors.New("Failed to start dhcpcd: " + err.Error())
    }

    return
}

/*
 * Stop dhcpcd on intf. Print an error if anything goes wrong
 */
func (d Dhcpcd) Stop () {
    myname := "Dhcpcd.Stop"
    _, _, err := sys.Run(d.CmdDhcpcd, "-k", d.Interface)
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
    stdout, _, err := sys.Run(d.CmdDhcpcd, "-4", "-T", d.Interface)
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

func DhcpcdFactory (intf string) (d Dhcpcd, err error) {
    dhcpcd, err := sys.BinaryPrefix("dhcpcd")
    if err != nil {
        return
    }

    d = Dhcpcd{intf, false, false, dhcpcd}

    return
}
