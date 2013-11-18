package network

import (
    "bytes"
    "net"
    "github.com/r3boot/rlib/sys"
)

type ResolvConf struct {
    Interface string
    Nameservers []net.IP
    Search string
}

func (rc *ResolvConf) UpdateResolvConf () (err error) {
    _, _, err = sys.Run("/sbin/resolvconf", "-u")
    return
}

func (rc *ResolvConf) AddConfig () {
    myname := "ResolvConf.AddConfig"
    r := *rc

    var stdin []byte
    stdin = append(stdin, "search " + r.Search + "\n"...)

    for _, nameserver := range r.Nameservers {
        stdin = append(stdin, "nameserver " + nameserver.String() + "\n"...)
    }

    _, _, err := sys.RunWithInput(bytes.NewReader(stdin), "/sbin/resolvconf", "-a", r.Interface)
    if err != nil {
        Log.Warning(myname, "Failed to add config for " + r.Interface + ": " + err.Error())
        return
    }

    if err = r.UpdateResolvConf(); err != nil {
        Log.Warning(myname, "Failed to update /etc/resolv.conf: " + err.Error())
        return
    }

    *rc = r
}

func (rc *ResolvConf) RemoveConfig () {
    myname := "ResolvConf.RemoveConfig"
    r := *rc

    _, _, err := sys.Run("/sbin/resolvconf", "-d", r.Interface)
    if err != nil {
        Log.Warning(myname, "Failed to remove config for " + r.Interface + ": " + err.Error())
    }

    if err = r.UpdateResolvConf(); err != nil {
        Log.Warning(myname, "Failed to update /etc/resolv.conf: " + err.Error())
        return
    }

    *rc = r
}

func ResolvConfFactory (intf, search string, servers []net.IP) (r ResolvConf) {
    r.Interface = intf
    r.Search = search
    for _, server := range servers {
        r.Nameservers = append(r.Nameservers, server)
    }

    return
}
