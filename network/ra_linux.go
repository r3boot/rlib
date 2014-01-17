package network

import (
    "github.com/r3boot/rlib/sys"
)

func (r RA) AcceptsRA () (result bool, err error) {
    value, err := sys.GetSysctl("net.ipv6.conf." + r.Interface + ".accept_ra")
    if err != nil {
        return
    }

    result = value[0] == sys.SYSCTL_ONE
    return
}

func (r RA) EnableRA () {
    ra_enabled, _ := r.AcceptsRA()

    if ! ra_enabled {
        sys.SetSysctl("net.ipv6.conf." + r.Interface + ".accept_ra", "1")
    }
}

func (r RA) DisableRA () {
   ra_enabled, _ := r.AcceptsRA()
   if ra_enabled {
        sys.SetSysctl("net.ipv6.conf." + r.Interface + ".accept_ra", "0")
    }
}

func RAFactory (intf string) (r RA, err error) {
    r = RA{intf, ""}

    return
}
