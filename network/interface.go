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
    link, err := LinkFactory(intf)
    if err != nil {
        return
    }

    ip, err := IpFactory(intf.Name)
    if err != nil {
        return
    }

    ra, err := RAFactory(intf.Name)
    if err != nil {
        return
    }

    wpa_supplicant, err := WpaSupplicantFactory(intf.Name)
    if err != nil {
        return
    }

    dhcpcd, err := DhcpcdFactory(intf.Name)
    if err != nil {
        return
    }

    resolvconf, err := ResolvConfFactory(intf.Name)
    if err != nil {
        return
    }

    i = Interface{
        intf,
        INTF_TYPE_UNKNOWN,
        link,
        ra,
        ip,
        wpa_supplicant,
        dhcpcd,
        resolvconf,
    }

    i.Type = i.GetType()

    return
}
