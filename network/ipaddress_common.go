package network

import (
    "net"
)

type Ip struct {
    Interface   string
}

func IpFactory (intf string) (i Ip, err error) {
    i = IP{intf}

    return
}
