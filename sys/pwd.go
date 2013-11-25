package sys

import (
    "errors"
    "io/ioutil"
    "strconv"
    "strings"
)

type PasswdInfo struct {
    Username    string
    Uid         int
    Gid         int
    Realname    string
    Homedir     string
    Shell       string
}

func GetPasswdInfo (user string) (info PasswdInfo, err error) {
    passwd_file := "/etc/passwd"

    if ! FileExists(passwd_file) {
        err = errors.New(passwd_file + " does not exists")
        return
    }

    content, err := ioutil.ReadFile(passwd_file)
    if err != nil {
        err = errors.New("Failed to read " + passwd_file + ": " + err.Error())
        return
    }

    for _, line := range strings.Split(string(content), "\n") {
        t := strings.Split(line, ":")
        if t[PASSWD_USERNAME] == user {
            info.Username = t[PASSWD_USERNAME]
            info.Uid, _ = strconv.Atoi(t[PASSWD_UID])
            info.Gid, _ = strconv.Atoi(t[PASSWD_GID])
            info.Realname = t[PASSWD_REALNAME]
            info.Homedir = t[PASSWD_HOMEDIR]
            info.Shell = t[PASSWD_SHELL]

            return
        }
    }

    return
}
