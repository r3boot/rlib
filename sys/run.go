package sys

import (
    "bytes"
    "io"
    "os/exec"
    "strings"
)

func Run (command string, args ...string) (stdout, stderr []string, err error) {
    var stdout_buf, stderr_buf bytes.Buffer
    cmd := exec.Command(command, args...)
    cmd.Stdout = &stdout_buf
    cmd.Stderr = &stderr_buf

    if err = cmd.Run(); err != nil {
        return
    }

    stdout = strings.Split(stdout_buf.String(), "\n")
    stderr = strings.Split(stderr_buf.String(), "\n")

    return
}

func Start (command string, args ...string) (err error) {
    cmd := exec.Command(command, args...)
    err = cmd.Start()
    return
}

func RunWithInput (stdin io.Reader, command string, args ...string) (stdout, stderr  []string, err error) {
    var stdout_buf, stderr_buf bytes.Buffer
    cmd := exec.Command(command, args...)
    cmd.Stdin = stdin
    cmd.Stdout = &stdout_buf
    cmd.Stderr = &stderr_buf

    if err = cmd.Run(); err != nil {
        return
    }

    stdout = strings.Split(stdout_buf.String(), "\n")
    stderr = strings.Split(stderr_buf.String(), "\n")

    return
}
