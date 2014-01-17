package network

import (
    "errors"
    "io/ioutil"
    "net"
    "strconv"
    "strings"
    "github.com/r3boot/rlib/sys"
)

func (r RIB) GetRoute (network net.IPNet) (result Route, err error) {
    af, proc_file, _, err := getAfDetails(network)
    if err != nil {
        return
    }

    content, err := ioutil.ReadFile(proc_file)
    if err != nil { return }

    for _, line := range strings.Split(string(content), "\n") {
        if strings.HasPrefix(line, "Iface") {
            continue
        } else if line == "" {
            continue
        }

        route, err := parseRoutingLine(af, line)
        if err != nil {
            continue
        }

        if (network.IP.Equal(route.Destination.IP)) &&
            (network.Mask.String() == route.Destination.Mask.String()) {
            result = route
            return result, err
        }
    }

    err = errors.New("Failed to find any routes")

    return
}

func (r RIB) AddRoute (network net.IPNet, gateway net.IP) (err error) {
    _, _, use_af, err := getAfDetails(network)
    if err != nil {
        return
    }

    _, _, err = sys.Run("/sbin/ip", use_af, "route", "add", network.String(),
        "via", gateway.String())
    if err != nil {
        return
    }

    return
}

func (r RIB) RemoveRoute (network net.IPNet) (err error) {
    _, _, use_af, err := getAfDetails(network)
    if err != nil {
        return
    }

    _, _, err = sys.Run("/sbin/ip", use_af, "route", "del", network.String())
    if err != nil {
        return
    }

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
    t := strings.Split(line, "\t")

    dest, err := HexToIp(t[1])
    if err != nil { return }

    mask, err := HexToIpv4Mask(t[7])
    if err != nil { return }

    gateway, err := HexToIp(t[2])
    if err != nil { return }

    route.Destination.IP = dest
    route.Destination.Mask = mask
    route.Gateway = gateway

    return
}

func parseIpv6RoutingLine (line string) (route Route, err error) {
    t := strings.Split(line, " ")

    // Skip routes on localhost
    if len(t) == LOOPBACK_LINE_LENGTH {
        err = errors.New("Loopback route")
        return
    }

    dest, err := HexToIp(t[0])
    if err != nil {
        return
    }

    mask, err := HexToBytes(t[1])
    if err != nil {
        return
    }

    gateway, err := HexToIp(t[4])
    if err != nil {
        return
    }

    route.Destination.IP = dest
    route.Destination.Mask = mask
    route.Gateway = gateway

    return
}

func getAfDetails (network net.IPNet) (af byte, proc_file, use_af string, err error) {

    if len(network.IP) == net.IPv4len {
        af = AF_INET
        proc_file = RIB_PROC_FILE_AF_INET
        use_af = IP_USE_AF_INET
    } else if len(network.IP) == net.IPv6len {
        af = AF_INET6
        proc_file = RIB_PROC_FILE_AF_INET6
        use_af = IP_USE_AF_INET6
    } else {
        err = errors.New("Failed to determine address family")
    }

    return
}
