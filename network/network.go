package network

import (
    "github.com/r3boot/rlib/logger"
)

type Network struct {}

var Log *logger.Log

func Setup (l *logger.Log) (err error) {
    myname := "networking.Setup"

    Log = l

    Log.Debug(myname, "Networking initialized")
    return
}
