package network

import (
    "errors"
    "strconv"
)

func mkUseAf(af byte) (result string, err error) {
    myname := "network.mkUseAf"
    if af == AF_INET {
        result = "inet"
    } else if af == AF_INET6 {
        result = "inet6"
    } else {
        err = errors.New("Unknown address family: " + strconv.Itoa(int(af)))
    }
    return
}
