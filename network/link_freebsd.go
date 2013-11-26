package network

import (
    "errors"
    "strconv"
    "strings"
    "github.com/r3boot/rlib/sys"
)

func (l Link) HasCarrier () (result bool, err error) {
    stdout, _, err := sys.Run(l.CmdIfconfig, l.Interface.Name)
    if err != nil {
        return
    }

    for _, line := range stdout {
        if strings.Contains(line, IFCONFIG_STATUS) {
            status := strings.Split(line, " ")[1]
            result = status == (IFCONFIG_CARRIER_ACTIVE || IFCONFIG_CARRIER_ASSOCIATED)
            return
        }
    }

    err = errors.New("Failed to determine carrier status")
    return
}

func (l Link) SetLinkStatus (link_status byte) (err error) {
    carrier, err := l.HasCarrier()
    if err != nil {
        return
    }

    var status string
    if link_status == LINK_UP {
        status = "up"
    } else if link_status == LINK_DOWN {
        status = "down"
    } else {
        err = errors.New("Unknown link status: " + + strconv.Itoa(int(link_status)))
        return
    }

    _, _, err = sys.Run(l.CmdIfconfig, l.Interface.Name, status)
    return
}

func (l Link) GetType () (intf_type byte, err error) {
    has_link, err := l.HasLink()
    if err != nil {
        return
    }

    if ! has_link {
        if err = l.SetLinkStatus(LINK_UP); err != nil {
            return
        }
    }

    stdout, _, err := sys.Run(l.CmdIfconfig, l.Interface.Name)
    if err != nil {
        return
    }

    for _, line := range stdout {
        if strings.Contains(line, IFCONFIG_MEDIA_ETHERNET) {
            intf_type = INTF_TYPE_ETHERNET
            return
        } else if strings.Contains(line, IFCONFIG_MEDIA_WIRELESS) {
            intf_type = INTF_TYPE_WIRELESS
            return
        }
    }

    err = errors.New("Unknown interface type")

    return
}

func (l Link) GetMtu () (mtu int, err error) {
    stdout, _, err := sys.Run(l.CmdIfconfig, l.Interface.Name)
    if err != nil {
        return
    }

    raw_mtu := strings.Fields(stdout[0])[5]
    mtu, err = strconv.Itoa(raw_mtu)
    return
}

func (l Link) SetMtu (mtu int) (err error) {
    _, _, err := sys.Run(l.CmdIfconfig, l.Interface.Name, "mtu", strconv.Itoa(mtu))
    return
}

func LinkFactory(intf net.Interface) (l Link, err error) {
    ifconfig, err := sys.BinaryPrefix("ifconfig")
    if err != nil {
        return
    }

    l.Interface = intf
    l.CmdIfconfig = ifconfig

    return
}
