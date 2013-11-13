package network

import (
    "io/ioutil"
    "log"
    "net"
    "strings"
    "github.com/r3boot/rlib/sys"
)

type Route struct {
    Destination net.IPNet
    Gateway net.IP
    Interface net.Interface
}

type RIB struct {
    Interface string
}

func (r RIB) GetIpv4Route (network net.IPNet) (result Route) {
    routing_file := "/proc/net/route"
    content, err := ioutil.ReadFile(routing_file)
    if err != nil { return }

    for _, line := range strings.Split(string(content), "\n") {
        if strings.HasPrefix(line, "Iface") { continue }
        if line == "" { continue }

        route, err := r.parseIpv4RoutingLine(line)
        if err != nil {
            log.Print("Failed to parse routing line: " + line)
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

func (r RIB) GetIpv4DefaultRoute () (route Route) {
    _, network, err := net.ParseCIDR("0.0.0.0/0")
    if err != nil { return }

    route = r.GetIpv4Route(*network)

    return
}

func (r RIB) AddRoute (network net.IPNet, gateway net.IP, af byte) {
    use_af := mkUseAf(af)
    if len(use_af) == 0 {
        return
    }

    _, _, err := sys.Run("/sbin/ip", use_af, "route", "add", network.String(),
        "via", gateway.String(), "dev", r.Interface)
    if err != nil {
        return
    }
}

func (r RIB) parseIpv4RoutingLine (line string) (route Route, err error) {
    t := strings.Split(line, "\t")

    intf, err := net.InterfaceByName(t[0])
    if err != nil { return }

    dest, err := HexToIpv4(t[1])
    if err != nil { return }

    mask, err := HexToIpv4Mask(t[7])
    if err != nil { return }

    gateway, err := HexToIpv4(t[2])
    if err != nil { return }

    route.Interface = *intf
    route.Destination.IP = dest
    route.Destination.Mask = mask
    route.Gateway = gateway

    return
}
