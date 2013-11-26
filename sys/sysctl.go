package sys

import (
    "strings"
)

func (s Sysctl) Get (key string) (value []byte, err error) {
    stdout, _, err := Run(s.CmdSysctl, key)
    value = []byte(strings.Fields(stdout)[1])
    return
}
