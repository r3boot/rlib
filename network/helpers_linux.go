package network

import (
    "errors"
    "strconv"
)

func mkUseAf(af byte) (result string, err error) {
    if af == AF_INET {
        result = "-4"
    } else if af == AF_INET6 {
        result = "-6"
    } else {
        err = errors.New("Unknown address family: " + strconv.Itoa(int(af)))
    }
    return
}
