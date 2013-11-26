package network

import (
    "encoding/hex"
    "errors"
    "net"
)

func HexToIp (hex_ip string) (ip net.IP, err error) {
    hex_ip_len := len(hex_ip)
    if hex_ip_len == AF_INET_STR_LEN {
        ip = make(net.IP, net.IPv4len)
        a, e := hex.DecodeString(hex_ip)
        if e != nil {
            err = e
            return
        }

        for i := 0; i <= 3; i++ {
            ip[i] = a[3-i]
        }
    } else if hex_ip_len == AF_INET6_STR_LEN {
        ip = make(net.IP, net.IPv6len)

        raw_ip := hex_ip[0:4] + ":" + hex_ip[4:8] + ":" + hex_ip[8:12] +
            ":" + hex_ip[12:16] + ":" + hex_ip[16:20] + ":" + hex_ip[20:24] +
            ":" + hex_ip[24:28] + ":" + hex_ip[28:32]

        ip = net.ParseIP(raw_ip)
        if err != nil {
            return
        }

    }

    return
}

func HexToIpv4Mask (hex_mask string) (mask net.IPMask, err error) {
    a, err := hex.DecodeString(hex_mask)
    if err != nil { return}

    mask = net.IPv4Mask(a[0], a[1], a[2], a[3])
    return
}

func HexToBytes (hex_mask string) (mask net.IPMask, err error) {
    a, err := hex.DecodeString(hex_mask)
    if err != nil {
        err = errors.New("Failed to decode hex: " + hex_mask)
        return
    }

    mask = net.CIDRMask(int(a[0]), 128)
    return
}
