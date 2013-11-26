package network

import (
    "strings"
    "github.com/r3boot/rlib/sys"
)

func (r RA) AcceptsRA () (result bool, err error) {
    stdout, _, err := sys.Run(r.CmdIfconfig, r.Interface, IFCONFIG_USE_AF_INET6)
    if err != nil {
        return
    }

    for _, line := range stdout {
        if strings.Contains(line, IFCONFIG_ND6_OPTIONS) {
            result = strings.Contains(line, IFCONFIG_ACCEPT_RTADV)
            return
        }
    }

    return
}

func (r RA) EnableRA () (err error) {
    enabled , err := r.AcceptsRA()
    if err != nil {
        return
    }

    if ! enabled {
        _, _, err = sys.Run(r.CmdIfconfig, r.Interface, IFCONFIG_USE_AF_INET6, IFCONFIG_ENABLE_RTADV)
    }

    return
}

func (r RA) DisableRA () (err error) {
    enabled, err := r.AcceptsRA()
    if err != nil {
        return
    }

    if enabled {
        _, _, err = sys.Run(r.CmdIfconfig, r.Interface, IFCONFIG_USE_AF_INET6, "-" + IFCONFIG_ENABLE_RTADV)
    }

    return
}

func RAFactory (intf string) (r RA, err error) {
    ifconfig, err := sys.BinaryPrefix("ifconfig")
    if err != nil {
        return
    }

    r = RA{intf, ifconfig}
    return
}
