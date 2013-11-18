package network

import (
    "errors"
    "io/ioutil"
    "net"
    "strconv"
    "strings"
    "github.com/r3boot/rlib/sys"
)

type Route struct {
    Destination net.IPNet
    Gateway net.IP
}

type RIB struct {
}

func (r RIB) GetRoute (network net.IPNet) (result Route) {
    myname := "RIB.GetRoute"

    af, proc_file, _ := getAfDetails(network)

    content, err := ioutil.ReadFile(proc_file)
    if err != nil { return }

    for _, line := range strings.Split(string(content), "\n") {
        if strings.HasPrefix(line, "Iface") { continue }
        if line == "" { continue }

        route, err := parseRoutingLine(af, line)
        if err != nil {
            Log.Warning(myname, "Failed to parse routing line: " + line)
            continue
        }

        if (network.IP.Equal(route.Destination.IP)) &&
            (network.Mask.String() == route.Destination.Mask.String()) {
            result = route
            return
        }
    }

    return
}

func (r RIB) GetDefaultRoute (af byte) (route Route) {
    myname := "RIB.GetDefaultRoute"
    var n string
    if af == AF_INET {
        n = "0.0.0.0/0"
    } else if af == AF_INET6 {
        n = "::/0"
    } else {
        Log.Fatal(myname, "Unsupported address family: " + strconv.Itoa(int(af)))
        return
    }

    _, network, err := net.ParseCIDR(n)
    if err != nil {
        return
    }
    route = r.GetRoute(*network)

    return
}

func (r RIB) GetDefaultGateway (af byte) (gateway net.IP) {
    return r.GetDefaultRoute(af).Gateway
}

func (r RIB) AddRoute (network net.IPNet, gateway net.IP) {
    _, _, use_af := getAfDetails(network)
    _, _, err := sys.Run("/sbin/ip", use_af, "route", "add", network.String(),
        "via", gateway.String())
    if err != nil {
        return
    }
}

func RIBFactory () (rib RIB) {
    rib = RIB{ }
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
    myname := "network.parseIpv6RoutingLine"
    t := strings.Split(line, " ")

    // Skip routes on localhost
    if len(t) == LOOPBACK_LINE_LENGTH {
        return
    }

    dest, err := HexToIp(t[0])
    if err != nil {
        Log.Warning(myname, err.Error())
        return
    }

    mask, err := HexToBytes(t[1])
    if err != nil {
        Log.Warning(myname, err.Error())
        return
    }

    gateway, err := HexToIp(t[4])
    if err != nil {
        Log.Warning(myname, err.Error())
        return
    }

    route.Destination.IP = dest
    route.Destination.Mask = mask
    route.Gateway = gateway

    return
}

func getAfDetails (network net.IPNet) (af byte, proc_file, use_af string) {
    myname := "network.getAfDetails"

    if len(network.IP) == net.IPv4len {
        af = AF_INET
        proc_file = RIB_PROC_FILE_AF_INET
        use_af = IP_USE_AF_INET
    } else if len(network.IP) == net.IPv6len {
        af = AF_INET6
        proc_file = RIB_PROC_FILE_AF_INET6
        use_af = IP_USE_AF_INET6
    } else {
        Log.Warning(myname, "Address family unknown")
    }

    return
}
