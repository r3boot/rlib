package network

import (
    "github.com/r3boot/rlib/sys"
)

type RA struct {
    Interface string
}

func (r RA) AcceptsRA () bool {
    value, err := sys.GetSysctl("net.ipv6.conf." + r.Interface + ".accept_ra")
    if err != nil {
        return false
    }

    return value[0] == sys.SYSCTL_ONE
}

func (r RA) EnableRA () {
    if r.AcceptsRA() {
        sys.SetSysctl("net.ipv6.conf." + r.Interface + ".accept_ra", "1")
    }
}

func (r RA) DisableRA () {
   if ! r.AcceptsRA() {
        sys.SetSysctl("net.ipv6.conf." + r.Interface + ".accept_ra", "0")
    }
}

