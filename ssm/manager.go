package ssm

import (
    "github.com/r3boot/rlib/sys"
)

type Manager interface {
    Start (name string) (err error)
    Stop (name string) (err error)
    Restart (name string) (err error)
    Enable (name string) (err error)
    Disable (name string) (err error)
    IsRunning (name string) (result bool)
    IsEnabled(name string) (result bool)
}

func ManagerFactory () (m Manager) {
    myname := "ssm.ManagerFactory"
    lsb := sys.LSBFactory()
    if lsb.Id == sys.DISTRO_ARCHLINUX {
        m = Manager(Systemd{})
    } else {
        Log.Warning(myname, "Unsupported distro")
    }
    return
}
