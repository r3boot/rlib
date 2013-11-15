package vpn

import (
    "github.com/r3boot/rlib/logger"
)

var Log *logger.Log

func Setup (l *logger.Log) (err error) {
    Log = l
    return
}
