package ssm

import (
    "errors"
    "github.com/r3boot/rlib/sys"
)

type Systemd struct {
    CmdService string
}

func (s Systemd) Start (name string) (err error) {
    _, _, err = sys.Run(s.CmdService, "start", name)
    return
}

func (s Systemd) Stop (name string) (err error) {
    _, _, err = sys.Run(s.CmdService, "stop", name)
    return
}

func (s Systemd) Restart (name string) (err error) {
    _, _, err = sys.Run(s.CmdService, "restart", name)
    return
}

func (s Systemd) IsRunning (name string) (result bool, err error) {
    err = errors.New("Feature not implemented")
    return
}

func (s Systemd) Enable (name string) (err error) {
    _, _, err = sys.Run(s.CmdService, "enable", name)
    return
}

func (s Systemd) Disable (name string) (err error) {
    _, _, err = sys.Run(s.CmdService, "disable", name)
    return
}

func (s Systemd) IsEnabled (name string) (result bool, err error) {
    err = errors.New("Feature not implemented") 
    return
}

func (s Systemd) Setup () (err error) {
    systemctl, err := sys.BinaryPrefix("systemctl")
    if err != nil {
        return
    }

    s.CmdService = systemctl

    return
}
