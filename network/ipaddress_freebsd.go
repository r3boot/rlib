package network

import (
    "errors"
    "net"
    "github.com/r3boot/rlib/sys"
)

func (i Ip) FlushAddresses (af byte) (err error) {
    use_af, err := mkUseAf(af)
    if err != nil {
        return
    }

    intf, err := net.InterfaceByName(i.Interface)
    if err != nil {
        err = errors.New("Failed to find interface " + i.Interface + ": " + err.Error())
        return
    }

    addrs, err := intf.Addrs()
    if err != nil {
        err = errors.New("Failed to find addresses for " + i.Interface + ": " + err.Error())
        return
    }

    for _, addr := range addrs {
        sys.Run("/sbin/ifconfig", i.Interface, use_af, addr.String(), "delete")
    }

    return
}

func (i Ip) FlushAllAddresses () (err error) {
    if err = i.FlushAddresses(AF_INET); err != nil {
        return
    }

    if err = i.FlushAddresses(AF_INET6); err != nil {
        return
    }
    return
}

func (i Ip) AddAddress (ip net.IPNet, af byte) (err error) {
    use_af, err := mkUseAf(af)
    if err != nil {
        return
    }

    sys.Run("/sbin/ifconfig", i.Interface, use_af, "alias", ip.String())

    return
}
