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
    IsEnabled (name string) (result bool, err error)
}

func ManagerFactory () (m Manager, err error) {
    uname, err := sys.Uname()

    if uname.Ident == sys.UNAME_LINUX {
        lsb, _:= sys.LSBFactory()
        if err != nil {
            return
        }

        if lsb.Id == sys.DISTRO_ARCHLINUX {
            systemd := new(Systemd)
            if err = SystemdSetup(); err != nil {
                err = errors.New("Failed to initialize systemd: " + err.Error())
            }
            m = Manager(systemd)
        } else {
            err = errors.New("Unsupported Linux distro: (" + lsb.Id + ")")
        }
    } else if uname.Ident == sys.UNAME_FREEBSD {
        service := new(Service)
        if err = ServiceSetup(); err != nil {
            err = errors.New("Failed to initialize service: " + err.Error())
        }
        m = Manager(service)
    } else {
        err = errors.New("Unknown UNIX release")
    }

    return
}
