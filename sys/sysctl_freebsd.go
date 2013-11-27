package sys

import (
    "strings"
)

func (sysctl *Sysctl) Get (key string) (value []byte, err error) {
    s := *sysctl
    stdout, _, err := Run(s.CmdSysctl, key)
    if err != nil {
        return
    }

    value = []byte(strings.Split(stdout[0], " ")[1])
    return
}
