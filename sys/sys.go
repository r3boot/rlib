package sys

import (
    "github.com/r3boot/rlib/logger"
)

var Log *logger.Log

func Setup (l *logger.Log) (err error) {
    myname := "sys.Setup"

    Log = l

    Log.Debug(myname, "System functions initialized")

    return
}
