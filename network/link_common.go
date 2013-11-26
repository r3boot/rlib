package network

import (
    "net"
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
