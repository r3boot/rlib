package network

import (
    "bytes"
    "net"
    "github.com/r3boot/rlib/sys"
)

type ResolvConf struct {
    Interface       string
    Search          string
    Nameservers     []net.IP
    CmdResolvconf   string
}

func (rc *ResolvConf) UpdateResolvConf () (err error) {
    r := *rc
    _, _, err = sys.Run(r.CmdResolvconf, "-u")
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

    _, _, err := sys.RunWithInput(bytes.NewReader(stdin), r.CmdResolvconf, "-a", r.Interface)
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

    _, _, err := sys.Run(r.CmdResolvconf, "-d", r.Interface)
    if err != nil {
        Log.Warning(myname, "Failed to remove config for " + r.Interface + ": " + err.Error())
    }

    if err = r.UpdateResolvConf(); err != nil {
        Log.Warning(myname, "Failed to update /etc/resolv.conf: " + err.Error())
        return
    }

    *rc = r
}

func ResolvConfFactory (intf string) (r ResolvConf, err error) {
    resolvconv, err := sys.BinaryPrefix("resolvconf")
    if err != nil {
        return
    }

    r = ResolvConf(intf, nil, make([]string), resolvconf)

    return
}
