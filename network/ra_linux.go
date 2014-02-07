package network

import (
    "github.com/r3boot/rlib/sys"
)

func (r RA) AcceptsRA () (result bool, err error) {
    s, err := sys.SysctlFactory()
    if err != nil {
        return
    }

    value, err := s.Get("net.ipv6.conf." + r.Interface + ".accept_ra")
    if err != nil {
        return
    }

    result = value[0] == sys.SYSCTL_ONE
    return
}

func (r RA) EnableRA () (err error) {
    s, err := sys.SysctlFactory()
    if err != nil {
        return
    }

    ra_enabled, _ := r.AcceptsRA()

    if ! ra_enabled {
        s.Set("net.ipv6.conf." + r.Interface + ".accept_ra", "1")
    }

    return
}

func (r RA) DisableRA () (err error){
    s, err := sys.SysctlFactory()
    if err != nil {
        return
    }

    ra_enabled, _ := r.AcceptsRA()
    if ra_enabled {
        s.Set("net.ipv6.conf." + r.Interface + ".accept_ra", "0")
    }

    return
}

func RAFactory (intf string) (r RA, err error) {
    r = RA{intf, ""}

    return
}
