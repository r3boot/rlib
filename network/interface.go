package network

import (
    "net"
)

type Interface struct {
    net.Interface
    Type byte
    WpaSupplicant WpaSupplicant
    Dhcpcd Dhcpcd
    Link
    RA
    Ip
    RIB
}
