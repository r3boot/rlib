package sys

import (
    "bytes"
    "os/exec"
    "strings"
)

func Run (command string, args ...string) (stdout, stderr  []string, err error) {
    myname := "sys.Run"
    var stdout_buf, stderr_buf bytes.Buffer
    cmd := exec.Command(command, args...)
    cmd.Stdout = &stdout_buf
    cmd.Stderr = &stderr_buf

    Log.Debug(myname, command + strings.Join(args, " "))
    if err = cmd.Run(); err != nil {
        return
    }

    stdout = strings.Split(stdout_buf.String(), "\n")
    stderr = strings.Split(stderr_buf.String(), "\n")

    return
}
