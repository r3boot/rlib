package network

import (
    "encoding/hex"
    "net"
    "strconv"
)

func HexToIpv4 (hex_ip string) (ip net.IP, err error) {
    ip = make(net.IP, net.IPv4len)
    a, err := hex.DecodeString(hex_ip)
    if err != nil { return }

    for i := 0; i <= 3; i++ {
        ip[i] = a[3-i]
    }

    return
}

func HexToIpv4Mask (hex_mask string) (mask net.IPMask, err error) {
    a, err := hex.DecodeString(hex_mask)
    if err != nil { return}

    mask = net.IPv4Mask(a[0], a[1], a[2], a[3])
    return
}

func mkUseAf(af byte) string {
    myname := "network.mkUseAf"
    if af == AF_INET {
        return "-4"
    } else if af == AF_INET6 {
        return "-6"
    } else {
        Log.Fatal(myname, "Unknown AF: " + strconv.Itoa(int(af)))
    }
    return ""
}
