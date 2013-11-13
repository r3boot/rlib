package network

const WPA_STATE_COMPLETED = "COMPLETED"
const WPA_SCAN_INTERVAL = 5
const WPA_CONNECT_INTERVAL = 500
const WPA_CONNECT_TIMEOUT = 6

const LINK_UP byte = 0x31
const LINK_DOWN byte = 0x30

const INTF_TYPE_UNKNOWN  byte = 0x0
const INTF_TYPE_LOOPBACK byte = 0x1
const INTF_TYPE_ETHERNET byte = 0x2
const INTF_TYPE_WIRELESS byte = 0x3

const LINK_LOOPBACK string = "772"
const LINK_ETHERNET string = "0x020000"
const LINK_WIRELESS string = "0x028000"

const AF_INET byte = 4
const AF_INET6 byte= 6
