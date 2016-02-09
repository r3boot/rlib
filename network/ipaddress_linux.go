package network

import (
    "github.com/r3boot/rlib/sys"
)

func (i Ip) FlushAddresses (af byte) (err error) {
    use_af, err := mkUseAf(af)
    if err != nil {
        return
    }

    if len(use_af) == 0 {
        return
    }

    sys.Run("/sbin/ip", use_af, "addr", "flush", "dev", i.Interface, "scope", "global")

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

func (i Ip) AddAddress (ip string, af byte) (err error) {
    use_af, err := mkUseAf(af)
    if err != nil {
        return
    }

    sys.Run("/sbin/ip", use_af, "addr", "add", ip, "dev", i.Interface)

    return
}
