package network

import (
    "errors"
    "net"
    "strconv"
    "strings"
    "github.com/r3boot/rlib/sys"
)

func (r RIB) GetRoute (network net.IPNet) (result Route, err error) {
    af, use_af, err := getAfDetails(network)
    if err != nil {
        return
    }

    stdout, _, err := sys.Run("/usr/bin/netstat", "-rn", "-f", use_af)
    if err != nil {
        return
    }

    for _, line := range stdout {
        if strings.HasPrefix(line, "Routing") {
            continue
        } else if line == "" {
            continue
        } else if strings.HasPrefix(line, "Internet") {
            continue
        } else if strings.HasPrefix(line, "Destination") {
            continue
        }

        route, e := parseRoutingLine(af, line)
        if e != nil {
            continue
        }

        if (network.IP.Equal(route.Destination.IP)) &&
            (network.Mask.String() == route.Destination.Mask.String()) {
                result = route
                return
            }
    }

    err = errors.New("Failed to find any routes")

    return
}

func (r RIB) AddRoute (network net.IPNet, gateway net.IP) (err error) {
    _, use_af, err := getAfDetails(network)
    if err != nil {
        return
    }

    _, _, err = sys.Run("/sbin/route", "add", use_af, network.String(), gateway.String())

    return
}

func parseRoutingLine (af byte, line string) (route Route, err error) {
    if af == AF_INET {
        route, err = parseIpv4RoutingLine(line)
    } else if af == AF_INET6 {
        route, err = parseIpv6RoutingLine(line)
    } else {
        err = errors.New("Unsupported address family: " + strconv.Itoa(int(af)))
    }
    return
}

func parseIpv4RoutingLine (line string) (route Route, err error) {
    t := strings.Fields(line)

    if strings.HasPrefix(t[1], "link") {
        err = errors.New("Link-local route")
        return
    }

    raw_dest := t[0]
    if raw_dest == "default" {
        raw_dest = "0.0.0.0/0"
    }

    _, destination, err := net.ParseCIDR(raw_dest)
    if err != nil {
        return
    }

    raw_gateway := t[1]
    gateway, _, err := net.ParseCIDR(raw_gateway + "/32")
    if err != nil {
        return
    }

    route.Destination = *destination
    route.Gateway = gateway

    return
}

func parseIpv6RoutingLine (line string) (route Route, err error) {
    t := strings.Fields(line)

    if strings.HasPrefix(t[2], "link") {
        err = errors.New("Link-local route")
        return
    } else if t[2] == "::1" {
        err = errors.New("Loopback route")
        return
    }

    raw_dest := t[0]
    if raw_dest == "default" {
        raw_dest = "::/0"
    }

    _, destination, err := net.ParseCIDR(raw_dest)
    if err != nil {
        return
    }

    raw_gateway := t[1]
    gateway, _, err := net.ParseCIDR(raw_gateway + "/128")
    if err != nil {
        return
    }

    route.Destination = *destination
    route.Gateway = gateway

    return
}

func getAfDetails (network net.IPNet) (af byte, use_af string, err error) {
    if len(network.IP) == net.IPv4len {
        af = AF_INET
        use_af = IFCONFIG_USE_AF_INET
    } else if len(network.IP) == net.IPv6len {
        af = AF_INET6
        use_af = IFCONFIG_USE_AF_INET6
    } else {
        err = errors.New("Failed to determine address family")
    }
    return
}
