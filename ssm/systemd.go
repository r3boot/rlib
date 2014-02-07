package ssm

import (
    "errors"
    "strings"
    "github.com/r3boot/rlib/sys"
)

var CmdSystemd string

type Systemd struct {
}

func (s Systemd) Start (name string) (err error) {
    _, _, err = sys.Run(CmdSystemd, "start", name)
    return
}

func (s Systemd) Stop (name string) (err error) {
    _, _, err = sys.Run(CmdSystemd, "stop", name)
    return
}

func (s Systemd) Restart (name string) (err error) {
    _, _, err = sys.Run(CmdSystemd, "restart", name)
    return
}

func (s Systemd) IsRunning (name string) (result bool, err error) {
    stdout, _, err := sys.Run(CmdSystemd, "status", name)
    if err != nil {
        result = false
        err = nil
        return
    }

    for _, line := range stdout {
        if strings.Contains(line, SYSTEMD_SERVICE_RUNNING) {
            result = true
            break
        } else if strings.Contains(line, SYSTEMD_SERVICE_NOT_RUNNING) {
            result = false
            break
        } else if strings.Contains(line, SYSTEMD_SERVICE_NOT_ENABLED) {
            err = errors.New("Service " + name + " is not enabled")
        } else {
            err = errors.New("Unknown error trying to determine service status of " + name)
        }
    }
    return
}

func (s Systemd) Enable (name string) (err error) {
    _, _, err = sys.Run(CmdSystemd, "enable", name)
    return
}

func (s Systemd) Disable (name string) (err error) {
    _, _, err = sys.Run(CmdSystemd, "disable", name)
    return
}

func (s Systemd) IsEnabled (name string) (result bool, err error) {
    err = errors.New("Feature not implemented")
    return
}

func SystemdSetup () (err error) {
    systemctl, err := sys.BinaryPrefix("systemctl")
    if err != nil {
        return
    }

    CmdSystemd = systemctl

    return
}
