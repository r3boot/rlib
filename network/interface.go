package network

import (
    "log"
    "errors"
    "net"
)

type Interface struct {
    Name            string
    Type            byte
    Intf            net.Interface
    Link            Link
    RA              RA
    Ip              Ip
    WpaSupplicant   WpaSupplicant
    Dhcpcd          Dhcpcd
    Resolvconf      ResolvConf
}

func InterfaceFactory (intf net.Interface) (i Interface, err error) {
    i = *new(Interface)
    i.Name = intf.Name
    i.Intf = intf

    link, err := LinkFactory(intf)
    if err != nil {
        err = errors.New("Failed to initialize Link: " + err.Error())
        return
    }
    i.Link = link

    ip, err := IpFactory(intf.Name)
    if err != nil {
        err = errors.New("Failed to initialize Ip: " + err.Error())
        return
    }
    i.Ip = ip

    ra, err := RAFactory(intf.Name)
    if err != nil {
        err = errors.New("Failed to initialize RA: " + err.Error())
        return
    }
    i.RA = ra


    i.Type, err = link.GetType()
    if err != nil {
        err = errors.New("Failed to get interface type: " + err.Error())
        return
    }

    var wpa_supplicant WpaSupplicant
    if i.Type == INTF_TYPE_WIRELESS {
        wpa_supplicant, err = WpaSupplicantFactory(intf.Name)
        if err != nil {
            err = errors.New("Failed to initialize wpa_supplicant: " + err.Error())
            return
        }
    }
    i.WpaSupplicant = wpa_supplicant

    dhcpcd, err := DhcpcdFactory(intf.Name)
    if err != nil {
        err = errors.New("Failed to initialize dhcpcd: " + err.Error())
        return
    }
    i.Dhcpcd = dhcpcd

    resolvconf, err := ResolvConfFactory(intf.Name)
    if err != nil {
        err = errors.New("Failed to initialize resolvconf: " + err.Error())
        return
    }
    i.Resolvconf = resolvconf

    log.Print(i)

    return
}
