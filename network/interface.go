package network

import (
    "net"
)

type Interface struct {
    net.Interface
    Type byte
    Link
    RA
    Ip
    WpaSupplicant WpaSupplicant
    Dhcpcd Dhcpcd
    ResolvConf ResolvConf
}

func InterfaceFactory (intf net.Interface) Interface {
    var i = Interface{
        intf,
        INTF_TYPE_UNKNOWN,
        Link{Interface: intf},
        RA{Interface: intf.Name},
        Ip{Interface: intf.Name},
        WpaSupplicant{Interface: intf.Name},
        Dhcpcd{Interface: intf.Name},
        ResolvConf{Interface: intf.Name},
    }

    i.Type = i.GetType()

    return i
}
