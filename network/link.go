package network

import (
    "net"
    "io/ioutil"
    "strconv"
    "github.com/r3boot/rlib/sys"
)

type Link struct {
    Interface net.Interface
}

/*
 * Check link status of interface. Returns true if interface is up, else
 * returns false.
 */
func (l Link) HasLink () bool {
    return (l.Interface.Flags & net.FlagUp) == net.FlagUp 
}

/*
 * Open /sys/class/net/<interface>/carrier and determine link status. Return
 * true if the content equals "1" (0x31), false otherwise. If the carrier file
 * cannot be read, return an error.
 */
func (l Link) HasCarrier () bool {
    carrier_file := "/sys/class/net/" + l.Interface.Name + "/carrier"

    content, err := ioutil.ReadFile(carrier_file)
    if err != nil {
        return false
    }

    return content[0] == LINK_UP
}

func (l Link) SetLinkStatus (link_status byte) {
    myname := "nic.setLinkStatus"

    var status string
    if link_status == LINK_UP {
        status = "up"
    } else if link_status == LINK_DOWN {
        status = "down"
    } else {
        Log.Fatal(myname, "Unknown link_status: " + strconv.Itoa(int(link_status)) + ", ignoring")
        return
    }

    _, _, err := sys.Run("/sbin/ip", "link", "set", l.Interface.Name, status)
    if err != nil {
        Log.Warning(myname, "Failed to set link status for " + l.Interface.Name)
    }
}

/*
 * Look in /sys/class/net/<interface>/type to see if this interface is
 * a loopback interface. Return if it is. Afterwards, look in
 * /sys/class/net/<interface>/device/class and check the pci class of the
 * device. If * this is "20000", it's an ethernet nic, if it's "28000", it's
 * a wireless nic. All other pci classes get flagged unknown.
 */
func (l Link) GetType () (intf_type byte) {
    if ! l.HasLink() {
        l.SetLinkStatus(LINK_UP)
    }

    flags_file := "/sys/class/net/" + l.Interface.Name + "/type"
    content, err := ioutil.ReadFile(flags_file)
    if err != nil {
        return
    }

    value := string(content[0:3])
    if value == LINK_LOOPBACK {
        intf_type = INTF_TYPE_LOOPBACK
        return
    }

    class_file := "/sys/class/net/" + l.Interface.Name + "/device/class"
    content, err = ioutil.ReadFile(class_file)
    if err != nil {
        return
    }

    value = string(content[0:8])
    if value == LINK_WIRELESS {
        intf_type = INTF_TYPE_WIRELESS
    } else if value == LINK_ETHERNET {
        intf_type = INTF_TYPE_ETHERNET
    }

    return
}

func (l Link) GetMTU () (mtu int, err error) {
    mtu_file := "/sys/class/net/" + l.Interface.Name + "/mtu"

    content, err := ioutil.ReadFile(mtu_file)
    if err != nil { return }

    mtu, err  = strconv.Atoi(string(content[0:3]))

    return
}

func (l Link) SetMTU (mtu int) (err error) {
    cur_mtu, err := l.GetMTU()
    if err != nil {
        return
    }

    if cur_mtu != mtu {
        mtu_file := "/sys/class/net/" + l.Interface.Name + "/mtu"
        value := []byte(strconv.Itoa(mtu))
        err = ioutil.WriteFile(mtu_file, value, 0755)
    }

    return
}
