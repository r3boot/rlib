package sys

import (
    "errors"
    "strings"
    "github.com/r3boot/rlib/sys"
)

func Uname () (u UnameInfo, err error) {
    var uname string
    uname, err = sys.BinaryPrefix("uname")
    if err != nil {
        return
    }

    stdout, _, err := sys.Run(uname, "-srpn")
    if err != nil {
        return
    }

    t := strings.Fields(stdout[0])
    u.Ident = t[0]
    u.Hostname = t[1]
    u.Release = t[2]
    u.Platform = t[3]

    stdout, _, err = sys.Run(uname, "-v")
    if err != nil {
        err = errors.New("Failed to find version info")
        return
    }

    u.Version = stdout[0]

    return
}
