package network

import (
    "net"
    "strings"
    "strconv"
    "github.com/r3boot/rlib/sys"
)

/*
 * Send count ARP Request packet(s) to ipaddr using arping. Return true if
 * the return code of arping is zero, false otherwise.
 */
func Arping(ipaddr net.IP, intf net.Interface, count int) (up bool, latency float64) {
    myname := "network.Arping"
    if ipaddr == nil {
        return
    }

    stdout, _, err := sys.Run("/usr/sbin/arping", "-I", intf.Name, "-c", strconv.Itoa(count), "-w", "3", ipaddr.String())

    var tot_latency float64 = 0

    if err == nil {
        up = true

        for _, line := range stdout {
            if ! strings.HasPrefix(line, "Unicast reply from") {
                continue
            }

            raw_latency := strings.Replace(strings.Split(line, " ")[6], "ms", "", -1)
            l, err := strconv.ParseFloat(raw_latency, 64)
            if err != nil {
                Log.Warning(myname, "Failed to parse float")
                continue
            }

            tot_latency += l
        }
    }

    latency = tot_latency / float64(count)

    return
}
