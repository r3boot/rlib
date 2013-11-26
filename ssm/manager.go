package ssm

import (
    "errors"
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

func ManagerFactory () (m Manager, err error) {
    uname, err := sys.Uname()
    lsb := sys.LSBFactory()

    if uname.Ident == UNAME_LINUX {
        if lsb.Id == sys.DISTRO_ARCHLINUX {
            systemd, err := SystemdFactory()
            if err != nil {
                return
            }
            m = Manager(systemd)
        } else {
            err = errors.New("Unsupported distro")
        }
    } else if uname.Ident == UNAME_FREEBSD {
        service, err := ServiceFactory()
        if err != nil {
            return
        }
        m = Manager(service)
    } else {
        err = errors.New("Unknown UNIX release")
    }
    return
}
