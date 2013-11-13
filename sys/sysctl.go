package sys

import (
    "io/ioutil"
    "strings"
)

func SetSysctl (key string, value string) {
    Run("/sbin/sysctl", "-w", key, "=", value)
}

func GetSysctl (key string) (value []byte, err error) {
    sysctl_file := "/proc/sys/" + strings.Replace(key, ".", "/", -1)

    value , err = ioutil.ReadFile(sysctl_file)
    return
}
