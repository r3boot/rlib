package logger

import (
    "fmt"
    "time"
)

const MSG_INFO byte    = 0x0
const MSG_WARNING byte = 0x1
const MSG_FATAL byte   = 0x2
const MSG_DEBUG byte   = 0x3

var MSG_STRING = map[byte]string{
    MSG_INFO:    "INFO",
    MSG_WARNING: "WARNING",
    MSG_FATAL:   "FATAL",
    MSG_DEBUG:   "DEBUG",
}

type Log struct {
    UseDebug bool
    UseTimestamp bool
    TimestampFormat string
}

func (l Log) Message (caller, message string, log_level byte) {
    msg := caller + "[" + MSG_STRING[log_level] + "]: " + message

    if l.UseTimestamp {
        if len(l.TimestampFormat) == 0 {
            l.TimestampFormat = time.RFC3339
        }
        timestamp := time.Now().Format(time.RFC3339)
        msg = timestamp + " " + msg
    }
    fmt.Println(msg)
}

func (l Log) Info (caller, message string) {
    l.Message(caller, message, MSG_INFO)
}

func (l Log) Warning (caller, message string) {
    l.Message(caller, message, MSG_WARNING)
}

func (l Log) Fatal (caller, message string) {
    l.Message(caller, message, MSG_FATAL)
}

func (l Log) Debug (caller, message string) {
    if l.UseDebug {
        l.Message(caller, message, MSG_DEBUG)
    }
}
