package ssm

import (
    "github.com/r3boot/rlib/sys"
)

type Systemd struct {
    CmdSystemctl    string
}

func (s Systemd) Start (name string) (err error) {
    _, _, err = sys.Run(s.CmdSystemctl, "start", name)
    return
}

func (s Systemd) Stop (name string) (err error) {
    _, _, err = sys.Run(s.CmdSystemctl, "stop", name)
    return
}

func (s Systemd) Restart (name string) (err error) {
    _, _, err = sys.Run(s.CmdSystemctl, "restart", name)
    return
}

func (s Systemd) IsRunning (name string) (result bool, err error) {
    err = errors.New("Feature not implemented")
    return
}

func (s Systemd) Enable (name string) (err error) {
    _, _, err = sys.Run(s.CmdSystemctl, "enable", name)
    return
}

func (s Systemd) Disable (name string) (err error) {
    _, _, err = sys.Run(s.CmdSystemctl, "disable", name)
    return
}

func (s Systemd) IsEnabled (name string) (result bool, err error) {
    err = errors.New("Feature not implemented") 
    return
}

func SystemdFactory () (s Systemd, err error) {
    systemctl, err := sys.BinaryPrefix("systemctl")
    if err != nil {
        return
    }

    s = Systemd{systemctl}

    return
}
