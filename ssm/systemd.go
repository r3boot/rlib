package ssm

import (
    "github.com/r3boot/rlib/sys"
)

type Systemd struct {
}

func (s Systemd) Start (name string) (err error) {
    _, _, err = sys.Run(CMD_SYSTEMCTL, "start", name)
    return
}

func (s Systemd) Stop (name string) (err error) {
    _, _, err = sys.Run(CMD_SYSTEMCTL, "stop", name)
    return
}

func (s Systemd) Restart (name string) (err error) {
    _, _, err = sys.Run(CMD_SYSTEMCTL, "restart", name)
    return
}

func (s Systemd) IsRunning (name string) (result bool) {
    return
}

func (s Systemd) Enable (name string) (err error) {
    _, _, err = sys.Run(CMD_SYSTEMCTL, "enable", name)
    return
}

func (s Systemd) Disable (name string) (err error) {
    _, _, err = sys.Run(CMD_SYSTEMCTL, "disable", name)
    return
}

func (s Systemd) IsEnabled (name string) (result bool) {
    return
}
