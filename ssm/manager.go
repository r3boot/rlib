package ssm

import (
    "errors"
    "github.com/r3boot/rlib/sys"
)

type Manager interface {
    Start (name string) (err error)
    Stop (name string) (err error)
    Restart (name string) (err error)
    Enable (name string) (err error)
    Disable (name string) (err error)
    IsRunning (name string) (result bool, err error)
    IsEnabled(name string) (result bool, err error)
    Setup () (err error)
}

func ManagerFactory () (m Manager, err error) {
    uname, err := sys.Uname()
    lsb, err := sys.LSBFactory()
    if err != nil {
        return
    }

    if uname.Ident == sys.UNAME_LINUX {
        if lsb.Id == sys.DISTRO_ARCHLINUX {
            m = Manager(Systemd{})
        } else {
            err = errors.New("Unsupported distro")
        }
    } else if uname.Ident == sys.UNAME_FREEBSD {
        m = Manager(Service{})
    } else {
        err = errors.New("Unknown UNIX release")
    }
    return
}
