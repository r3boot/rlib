package network

import (
    "net"
    "strconv"
    "github.com/r3boot/rlib/sys"
)

/*
 * Send count icmp/ipv6-icmp packet(s) to ipaddr using fping. Return true if
 * the return code of fping is zero, false otherwise.
 */
func Ping(ipaddr net.IP, count int) bool {
    myname := "network.Ping"
    var fping string
    if ipaddr == nil {
        return false
    }

    ip_len := len(ipaddr)
    if ip_len == net.IPv4len {
        fping = "/usr/sbin/fping"
    } else if ip_len == net.IPv6len {
        fping = "/usr/sbin/fping6"
    } else  {
        Log.Warning(myname, "Unknown address length: " + strconv.Itoa(ip_len))
        return false
    }

    _, _, err := sys.Run(fping, "-q", "-c", strconv.Itoa(int(count)), ipaddr.String())
    return err == nil
}

/*
 * Send three ping packets to ipaddr using Ping and return the results
 */
func IsReachable (ipaddr net.IP) bool {
    return Ping(ipaddr, 3)
}
