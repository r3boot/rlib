package sys

import (
    "errors"
    "io/ioutil"
    "strconv"
    "strings"
)

type LSB struct {
    Id          string
    Description string
    Release     string
    LsbVersion  float64
}

func LSBFactory () (l LSB, err error) {
    lsb_file := "/etc/lsb-release"

    if ! FileExists(lsb_file) {
        err = errors.New("Failed to locate LSB file " + lsb_file)
        return
    }

    content, err := ioutil.ReadFile(lsb_file)
    if err != nil {
        err = errors.New("Failed to read LSB file " + lsb_file + ": " + err.Error())
    }

    for _, line := range strings.Split(string(content), "\n") {
        t := strings.Split(line, "=")
        switch t[0] {
            case LSB_VERSION: {
                l.LsbVersion, err = strconv.ParseFloat(t[1], 64)
                if err != nil {
                    err = errors.New("Failed to parse float: " + err.Error())
                    l = LSB{}
                    return
                }
            }
            case LSB_D_ID: {
                l.Id = strings.Join(t[1:], " ")
            }
            case LSB_D_RELEASE: {
                l.Release = strings.Join(t[1:], " ")
            }
            case LSB_D_DESCRIPTION: {
                l.Description = strings.Replace(strings.Join(t[1:], " "), "\"", "", -1)
            }
        }
    }

    return
}
