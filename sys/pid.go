package sys

import (
    "errors"
    "strconv"
    "strings"
    "io/ioutil"
)

func PidOf (cmd string) (pid int, err error) {
    var raw_max_pid []byte

    sysctl, err := SysctlFactory()
    if err != nil {
        return
    }

    uname, err := Uname()
    if err != nil {
        return
    }

    if uname.Ident == UNAME_LINUX {
        raw_max_pid, err = sysctl.Get("kernel.pid_max")
        if err != nil {
            err = errors.New("Failed to retrieve kernel.pid_max: " + err.Error())
            return
        }
    } else if uname.Ident == UNAME_FREEBSD {
        raw_max_pid, err = sysctl.Get("kern.pid_max")
        if err != nil {
            err = errors.New("Failed to retrieve kern.pid_max: " + err.Error())
            return
        }
    } else {
        err = errors.New("Unsupported UNIX release")
    }

    max_pid, err := strconv.Atoi(strings.Split(string(raw_max_pid), "\n")[0])
    if err != nil {
        err = errors.New("Failed to convert raw pid to int: " + err.Error())
        return
    }

    proc_cmd := strings.Replace(cmd, " ", "", -1)

    for pid = 0; pid < max_pid; pid ++ {
        cmdline_file := "/proc/" + strconv.Itoa(pid) + "/cmdline"
        if ! FileExists(cmdline_file) {
            continue
        }

        raw_cmdline, e := ioutil.ReadFile(cmdline_file)
        cmdline := strings.Split(string(raw_cmdline), "\n")[0]
        if e != nil {
            err = errors.New("Failed to read " + cmdline_file + ": " + err.Error())
            pid = 0
            return
        }

        cmdline = strings.Replace(cmdline, "\000", "", -1)
        cmdline = strings.Replace(cmdline, " ", "", -1)

        if cmdline == proc_cmd {
            return
        }

    }

    pid = 0
    err = errors.New("Failed to find pid for " + cmd)

    return
}
