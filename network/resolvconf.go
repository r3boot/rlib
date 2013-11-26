package network

import (
    "bytes"
    "errors"
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

func (rc *ResolvConf) AddConfig () (err error) {
    r := *rc

    var stdin []byte
    stdin = append(stdin, "search " + r.Search + "\n"...)

    for _, nameserver := range r.Nameservers {
        stdin = append(stdin, "nameserver " + nameserver.String() + "\n"...)
    }

    _, _, err = sys.RunWithInput(bytes.NewReader(stdin), r.CmdResolvconf, "-a", r.Interface)
    if err != nil {
        err = errors.New("Failed to add config for " + r.Interface + ": " + err.Error())
        return
    }

    if err = r.UpdateResolvConf(); err != nil {
        err = errors.New("Failed to update /etc/resolv.conf: " + err.Error())
        return
    }

    *rc = r
    return
}

func (rc *ResolvConf) RemoveConfig () {
    r := *rc

    _, _, err := sys.Run(r.CmdResolvconf, "-d", r.Interface)
    if err != nil {
        err = errors.New("Failed to remove config for " + r.Interface + ": " + err.Error())
    }

    if err = r.UpdateResolvConf(); err != nil {
        err = errors.New("Failed to update /etc/resolv.conf: " + err.Error())
        return
    }

    *rc = r
}

func ResolvConfFactory (intf string) (r ResolvConf, err error) {
    resolvconf, err := sys.BinaryPrefix("resolvconf")
    if err != nil {
        return
    }

    r = ResolvConf{intf, "", *new([]net.IP), resolvconf}

    return
}
