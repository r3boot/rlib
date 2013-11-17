package ssm

import (
    "github.com/r3boot/rlib/logger"
    "github.com/r3boot/rlib/sys"
)

var Log *logger.Log

func Setup (l *logger.Log) (err error) {
    Log = l

    err = sys.Setup(l)
    return
}
