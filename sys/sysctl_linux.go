package sys

import (
    "io/ioutil"
    "strings"
)

func (s Sysctl) Get (key string) (value []byte, err error) {
    sysctl_file := "/proc/sys/" + strings.Replace(key, ".", "/", -1)

    value , err = ioutil.ReadFile(sysctl_file)
    return
}
