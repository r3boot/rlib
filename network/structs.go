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

const AF_INET_STR_LEN int = 8
const AF_INET6_STR_LEN int = 32

const LOOPBACK_LINE_LENGTH int = 16

const RIB_PROC_FILE_AF_INET string = "/proc/net/route"
const RIB_PROC_FILE_AF_INET6 string = "/proc/net/ipv6_route"

const IP_USE_AF_INET string = "-4"
const IP_USE_AF_INET6 string = "-6"

const IFCONFIG_USE_AF_INET string = "inet"
const IFCONFIG_USE_AF_INET6 string = "inet6"
const IFCONFIG_ND6_OPTIONS string = "nd6 options"
const IFCONFIG_ACCEPT_RTADV string = "ACCEPT_RTADV"
const IFCONFIG_ENABLE_RTADV string = "accept_rtadv"
const IFCONFIG_CARRIER_ACTIVE string = "active"
const IFCONFIG_CARRIER_ASSOCIATED string = "associated"
const IFCONFIG_STATUS string = "status"
const IFCONFIG_MEDIA_ETHERNET string = "media: Ethernet"
const IFCONFIG_MEDIA_WIRELESS string = "media: IEEE 802.11"
