package ssm

import (
    "errors"
    "strings"
    "github.com/r3boot/rlib/sys"
)

var CmdService string

type Service struct {
}

func (s Service) Start (name string) (err error) {
    _, _, err = sys.Run(CmdService, name, "start")
    return
}

func (s Service) Stop (name string) (err error) {
    _, _, err = sys.Run(CmdService, name, "stop")
    return
}

func (s Service) Restart (name string) (err error) {
    _, _, err = sys.Run(CmdService, name, "restart")
    return
}

func (s Service) IsRunning (name string) (result bool, err error) {
    stdout, _, _ := sys.Run(CmdService, name, "status")
    if strings.Contains(stdout[0], SERVICE_SERVICE_RUNNING) {
        result = true
    } else if strings.Contains(stdout[0], SERVICE_SERVICE_NOT_RUNNING) {
        result = false
    } else if strings.Contains(stdout[0], SERVICE_SERVICE_NOT_ENABLED) {
        err = errors.New("Service " + name + " is not enabled")
    } else {
        err = errors.New("Unknown error trying to determine service status of " + name)
    }
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

func ServiceSetup () (err error) {
    service, err := sys.BinaryPrefix("service")
    if err != nil {
        err = errors.New("Failed to initialize Service: " + err.Error())
    }

    CmdService = service

    return
}
