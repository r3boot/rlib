package vpn

import (
    "errors"
    "net"
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

    apid, _, err := a.GetPid()
    if err != nil {
        err = errors.New("Failed to retrieve pid: " + err.Error())
        return
    }

    autossh, err := os.FindProcess(apid)
    if err != nil {
        err = errors.New("Failed to convert pid to Process: " + err.Error())
        return
    }
    autossh.Kill()
    autossh.Wait()

    return
}

func (assh *AutoSSH) GetPid () (assh_pid, ssh_pid int, err error) {
    a := *assh

    cmdline := "/usr/bin/autossh -M " + strconv.Itoa(a.EchoPort) + " -F " + a.ConfigFile + " -N -y " + a.Name
    assh_pid, err = sys.PidOf(cmdline)
    if err != nil {
        assh_pid = 0
        return
    }

    eq := strconv.Itoa(a.EchoPort)
    er := strconv.Itoa(a.EchoPort + 1)
    cmdline = "/usr/bin/ssh -L " + eq + ":127.0.0.1:" + eq + " -R " + eq + ":127.0.0.1:" + er + " -F " + a.ConfigFile + " -N -y " + a.Name
    assh_pid, err = sys.PidOf(cmdline)
    if err != nil {
        assh_pid = 0
        return
    }

    return
}

func (assh *AutoSSH) IsRunning () (result bool) {
    a := *assh
    pid, _, err := a.GetPid()
    if err != nil {
        result = false
    } else if pid > 1 {
        result = true
    }

    return
}

func (assh *AutoSSH) IsConnected () bool {
    a := *assh

    _, err := net.Dial("tcp", "127.0.0.1:" + strconv.Itoa(a.EchoPort))
    if err != nil {
        return false
    }

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

    pwd, err := sys.GetPasswdInfo(a.RunAs)
    if err != nil {
        return
    }

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
