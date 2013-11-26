package ssm

import (
    "github.com/r3boot/rlib/sys"
)

type Service struct {
    CmdService  string
}

func (s Service) Start (name string) (err error) {
    _, _, err = sys.Run(s.CmdService, name, "start")
    return
}

func (s Service) Stop (name string) (err error) {
    _, _, err = sys.Run(s.CmdService, name, "stop")
    return
}

func (s Service) Restart (name string) (err error) {
    _, _, err = sys.Run(s.CmdService, name, "restart")
    return
}

func (s Service) IsRunning (name string) (result bool, err error) {
    err = errors.New("Feature not implemented")
    return
}

func (s Service) Enable (name string) (err error) {
    err = errors.New("Feature not implemented")
    return
}

func (s Service) Disable (name string) (err error) {
    err = errors.New("Feature not implemented")
    return
}

func (s Service) IsEnabled (name string) (result bool, err error) {
    err = errors.New("Feature not implemented")
    return
}

func ServiceFactory () (s Service, err error) {
    service, err := sys.BinaryPrefix("service")
    if err != nil {
        return
    }

    s = Service{service}

    return
}
