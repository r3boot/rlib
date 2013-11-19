package vpn

import (
    "errors"
    "os"
    "io/ioutil"
    "strconv"
    "strings"
    "github.com/r3boot/rlib/sys"
)

type AutoSSH struct {
    Id          int
    Name        string
    RunAs       string
    ConfigFile  string
    AgentFile   string
    EchoPort    int
}

func (assh *AutoSSH) Start () (err error) {
    a := *assh

    err = a.SetAgentVars()
    if err != nil {
        err = errors.New("Failed to set ssh-agent vars: " + err.Error())
        return
    }

    err = sys.Start("/usr/bin/autossh", "-M", strconv.Itoa(a.EchoPort), "-f", "-F", a.ConfigFile , "-N", "-y", a.Name)
    if err != nil {
        err = errors.New("Failed to start AutoSSH for " + a.Name + ": " + err.Error())
    }

    return
}

func (assh *AutoSSH) Stop () (err error) {
    a := *assh

    pid, err := a.GetPid()
    if err != nil {
        err = errors.New("Failed to retrieve pid: " + err.Error())
        return
    }

    process, err := os.FindProcess(pid)
    if err != nil {
        err = errors.New("Failed to convert pid to Process: " + err.Error())
        return
    }

    process.Kill()
    process.Wait()

    return
}

func (assh *AutoSSH) GetPid () (pid int, err error) {
    a := *assh
    cmdline := "/usr/bin/autossh -M " + strconv.Itoa(a.EchoPort) + " -F " + a.ConfigFile + " -N -y " + a.Name
    pid, err = sys.PidOf(cmdline)

    if err != nil {
        pid = 0
    }

    return
}

func (assh *AutoSSH) IsRunning () (result bool) {
    a := *assh
    pid, err := a.GetPid()
    if err != nil {
        result = false
    } else if pid > 1 {
        result = true
    }

    return
}

func (assh *AutoSSH) IsConnected () bool {

    return true
}

func (assh *AutoSSH) SetAgentVars () (err error) {
    a := *assh

    if ! sys.FileExists(a.AgentFile) {
        err = errors.New(a.AgentFile + " does not exist")
        return
    }

    content, err := ioutil.ReadFile(a.AgentFile)
    if err != nil {
        err = errors.New("Failed to read agent vars file: " + err.Error())
        return
    }

    for _, line := range strings.Split(string(content), "\n") {
        if strings.HasPrefix(line, "export ") {
            line = strings.Replace(line, "export ", "", -1)
        }

        t := strings.Split(line, "=")
        if len(t) != 2 {
            continue
        }

        k := t[0]
        v := strings.Replace(t[1], "\"", "", -1)
        os.Setenv(k, v)
    }

    return
}

func (assh *AutoSSH) ReadConfig () (err error) {
    a := *assh

    pwd := sys.GetPasswdInfo(a.RunAs)
    a.ConfigFile = pwd.Homedir + "/.ssh/config"
    a.AgentFile = pwd.Homedir + "/.ssh/ssh-agent.vars"
    a.EchoPort = 50100 + (a.Id * 2)

    if ! sys.FileExists(a.ConfigFile) {
        err = errors.New(a.ConfigFile + " does not exists")
        return
    }

    *assh = a

    return
}

func AutoSSHFactory (id int, name, runas string) (a AutoSSH, err error) {
    a = AutoSSH{Id: id, Name: name, RunAs: runas}
    err = a.ReadConfig()

    return
}
