package network

import (
    "net"
    "io/ioutil"
    "strconv"
    "github.com/r3boot/rlib/sys"
)

type Link struct {
    Interface       net.Interface
    CmdIfconfig     string
    CmdIp           string
}

/*
 * Check link status of interface. Returns true if interface is up, else
 * returns false.
 */
func (l Link) HasLink () bool {
    return (l.Interface.Flags & net.FlagUp) == net.FlagUp 
}
