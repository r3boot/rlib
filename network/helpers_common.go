package network

import (
	"crypto/rand"
	"fmt"
	"net"
)

func GenerateMacAddr() (hwaddr net.HardwareAddr, err error) {
	buf := make([]byte, 6)
	if _, e := rand.Read(buf); e != nil {
		err = e
		return
	}

	// set local bit
	buf[0] |= 2

	hwaddr_s := fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x",
		buf[0], buf[1], buf[2], buf[3], buf[4], buf[5],
	)

	if hwaddr, err = net.ParseMAC(hwaddr_s); err != nil {
		hwaddr = nil
	}

	return
}
