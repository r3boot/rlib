package network

import (
    "net"
    "github.com/r3boot/rlib/sys"
)

type Ip struct {
    Interface string
}

func (i Ip) FlushAddresses (af byte) {
    use_af := mkUseAf(af)
    if len(use_af) == 0 {
        return
    }

    sys.Run("/sbin/ip", use_af, "addr", "flush", "dev", i.Interface, "scope", "global")
}

func (i Ip) FlushAllAddresses () {
    i.FlushAddresses(AF_INET)
    i.FlushAddresses(AF_INET6)
}

func (i Ip) AddAddress (ip net.IPNet, af byte) {
    use_af := mkUseAf(af)
    if len(use_af) == 0 {
        return
    }

    sys.Run("/sbin/ip", use_af, "addr", "add", ip.String(), "dev", i.Interface)
}
