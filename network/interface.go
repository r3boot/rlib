package network

import (
    "net"
)

type Interface struct {
    net.Interface
    Type byte
    WpaSupplicant WpaSupplicant
    Dhcpcd Dhcpcd
    Link
    RA
    Ip
    RIB
}

func InterfaceFactory (intf net.Interface) Interface {
    return Interface{
        intf,
        INTF_TYPE_UNKNOWN,
        WpaSupplicant{Interface: intf.Name},
        Dhcpcd{Interface: intf.Name},
        Link{Interface: intf},
        RA{Interface: intf.Name},
        Ip{Interface: intf.Name},
        RIB{Interface: intf.Name},
    }
}
