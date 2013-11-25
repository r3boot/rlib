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
    Resolver ResolvConf
}

func InterfaceFactory (intf net.Interface) (i Interface, err error) {
    wpa_supplicant, err := WpaSupplicantFactory(intf.Name)
    if err != nil {
        return
    }

    i = Interface{
        intf,
        INTF_TYPE_UNKNOWN,
        Link{Interface: intf},
        RA{Interface: intf.Name},
        Ip{Interface: intf.Name},
        wpa_supplicant,
        Dhcpcd{Interface: intf.Name},
        ResolvConf{Interface: intf.Name},
    }

    i.Type = i.GetType()

    return
}
