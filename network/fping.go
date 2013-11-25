package network

import (
    "errors"
    "net"
    "strings"
    "strconv"
    "github.com/r3boot/rlib/sys"
)

/*
 * Send count icmp/ipv6-icmp packet(s) to ipaddr using fping. Return true if
 * the return code of fping is zero, false otherwise.
 */
func Fping(ipaddr net.IP, count int) (up bool, latency float64) {
    var cmd string

    if ipaddr == nil {
        return
    }

    ip_len := len(ipaddr)
    if ip_len == net.IPv4len {
        cmd = "fping"
    } else if ip_len == net.IPv6len {
        cmd = "fping6"
    } else  {
        err = errors.New("Unknown address length: " + strconv.Itoa(ip_len))
        return
    }

    fping, err := sys.BinaryPrefix(cmd)
    if err != nil {
        return
    }

    _, stderr, err := sys.Run(fping, "-q", "-c", strconv.Itoa(int(count)), ipaddr.String())
    if err == nil {
        up = true

        latency, err = strconv.ParseFloat(strings.Split(stderr[0], "/")[7], 64)
        if err != nil {
            Log.Warning(myname, "Error parsing float: " + strings.Split(stderr[0], "/")[7])
            up = false
            return
        }
    }

    return
}

/*
 * Send three ping packets to ipaddr using Ping and return the results
 */
func IsReachable (ipaddr net.IP) (up bool, latency float64) {
    return Fping(ipaddr, 3)
}
