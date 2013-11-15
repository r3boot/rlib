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
    Interface net.Interface
}

type RIB struct {
    Af byte
    procFile string
    ipPrefix string
}

func (r RIB) GetRoute (network net.IPNet) (result Route) {
    myname := "RIB.GetRoute"

    content, err := ioutil.ReadFile(r.procFile)
    if err != nil { return }

    for _, line := range strings.Split(string(content), "\n") {
        if strings.HasPrefix(line, "Iface") { continue }
        if line == "" { continue }

        route, err := parseRoutingLine(r.Af, line)
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

func (r RIB) GetDefaultRoute () (route Route) {
    myname := "RIB.GetDefaultRoute"
    var n string
    if r.Af == AF_INET {
        n = "0.0.0.0/0"
    } else if r.Af == AF_INET6 {
        n = "::/0"
    } else {
        Log.Fatal(myname, "Unsupported address family: " + strconv.Itoa(int(r.Af)))
        return
    }

    _, network, err := net.ParseCIDR(n)
    if err != nil {
        return
    }
    route = r.GetRoute(*network)

    return
}

func (r RIB) AddRoute (network net.IPNet, gateway net.IP, af byte) {
    _, _, err := sys.Run("/sbin/ip", r.ipPrefix, "route", "add", network.String(),
        "via", gateway.String())
    if err != nil {
        return
    }
}

func RIBFactory (af byte) (rib RIB) {
    myname := "network.RIBFactory"

    var proc_file, ip_prefix string

    if af == AF_INET {
        proc_file = "/proc/net/route"
        ip_prefix = "-4"
    } else if af == AF_INET6 {
        proc_file = "/proc/net/ipv6_route"
        ip_prefix = "-6"
    } else {
        Log.Fatal(myname, "Unsupported addres family: " + strconv.Itoa(int(af)))
        return
    }

    rib = RIB{
        af,
        proc_file,
        ip_prefix,
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

    intf, err := net.InterfaceByName(t[0])
    if err != nil { return }

    dest, err := HexToIp(t[1])
    if err != nil { return }

    mask, err := HexToIpv4Mask(t[7])
    if err != nil { return }

    gateway, err := HexToIp(t[2])
    if err != nil { return }

    route.Interface = *intf
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

    intf, err := net.InterfaceByName(t[11])
    if err != nil {
        Log.Warning(myname, err.Error())
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

    route.Interface = *intf
    route.Destination.IP = dest
    route.Destination.Mask = mask
    route.Gateway = gateway

    return
}
