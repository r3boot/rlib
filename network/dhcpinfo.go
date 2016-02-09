package network

import (
	"errors"
	"fmt"
	"github.com/d2g/dhcp4"
	"github.com/d2g/dhcp4client"
	"net"
)

func GetOffer() (ip net.IP, subnet net.IPNet, err error) {
	ip = nil
	subnet.IP = nil
	subnet.Mask = nil

	// Generate a random mac address for this request
	hwaddr, e := GenerateMacAddr()
	if e != nil {
		err = e
		return
	}

	// Setup DHCP socket
	//socket, e := dhcp4client.NewInetSock(local_addr, remote_addr)
	socket, e := dhcp4client.NewInetSock()
	if e != nil {
		err = e
		return
	}

	// Setup DHCP client
	client, e := dhcp4client.New(
		dhcp4client.HardwareAddr(hwaddr),
		dhcp4client.Connection(socket),
	)

	// Send a DHCPREQUEST
	result, packet, e := client.Request()

	if !result {
		msg := fmt.Sprintf("DHCPREQUEST failed: %v", e)
		err = errors.New(msg)
		return
	}

	if e != nil {
		networkError, ok := e.(*net.OpError)
		if ok && networkError.Timeout() {
			msg := fmt.Sprintf("Cannot find DHCP server: %v", networkError)
			err = errors.New(msg)
		} else {
			msg := fmt.Sprintf("DHCPREQUEST failed: %v", networkError)
			err = errors.New(msg)
		}
		return
	}

	// Determine the netmask of the network
	options := packet.ParseOptions()
	mask := make(net.IPMask, net.IPv4len)
	mask = options[dhcp4.OptionCode(dhcp4.OptionSubnetMask)]

	// Setup the return structures
	ip = packet.YIAddr()
	subnet.IP = ip.To4().Mask(mask)
	subnet.Mask = mask

	return
}
