package network

import (
    "errors"
    "net"
    "strconv"
)

func (r RIB) GetDefaultRoute (af byte) (route Route, err error) {
    var n string
    if af == AF_INET {
        n = "0.0.0.0/0"
    } else if af == AF_INET6 {
        n = "::/0"
    } else {
        err = errors.New("Unsupported address family: " + strconv.Itoa(int(af)))
        return
    }

    _, network, err := net.ParseCIDR(n)
    if err != nil {
        return
    }

    route, err = r.GetRoute(*network)

    return
}

func (r RIB) GetDefaultGateway (af byte) (gateway net.IP, err error) {
    route, err := r.GetDefaultRoute(af)
    if err != nil {
        return
    }
    gateway = route.Gateway
    return
}

func RIBFactory () (rib RIB) {
    rib = RIB{}
    return
}

