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
func Arping(ipaddr net.IP, intf net.Interface, count int) (up bool, latency float64, err error) {
    if ipaddr == nil {
        return
    }

    arping, err := sys.BinaryPrefix("arping")
    if err != nil {
        return
    }

    stdout, _, err := sys.Run(arping, "-I", intf.Name, "-w", "3000", "-c", strconv.Itoa(count), ipaddr.String())

    var tot_latency float64 = 0

    if err == nil {
        up = true

        for _, line := range stdout {
            if ! strings.Contains(line, "bytes from") {
                continue
            }

            raw_latency := strings.Fields(line)[6]
            raw_latency = strings.Split(raw_latency, "=")[1]
            l, err := strconv.ParseFloat(raw_latency, 64)
            if err != nil {
                continue
            }

            tot_latency += (l / 1000)
        }
    }

    latency = (tot_latency / float64(count))

    return
}
