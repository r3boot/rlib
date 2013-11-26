package network

type Ip struct {
    Interface   string
}

func IpFactory (intf string) (i Ip, err error) {
    i = Ip{intf}

    return
}
