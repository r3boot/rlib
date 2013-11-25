package sys

const SYSCTL_ONE        byte  = 0x31
const SYSCTL_ZERO       byte = 0x30

const LSB_VERSION       string = "LSB_VERSION"
const LSB_D_ID          string = "DISTRIB_ID"
const LSB_D_RELEASE     string = "DISTRIB_RELEASE"
const LSB_D_DESCRIPTION string = "DISTRIB_DESCRIPTION"

const DISTRO_ARCHLINUX  string = "arch"

const PASSWD_USERNAME   int = 0
const PASSWD_UID        int = 2
const PASSWD_GID        int = 3
const PASSWD_REALNAME   int = 4
const PASSWD_HOMEDIR    int = 5
const PASSWD_SHELL      int = 6

var ETC_PATHS []string = []string{"/etc", "/usr/local/etc"}
