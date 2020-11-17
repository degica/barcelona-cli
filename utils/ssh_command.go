package utils

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type SshConfig interface {
	GetCertPath() string
	GetPrivateKeyPath() string
	IsDebug() bool
}

type SshCommand interface {
	Run(command string) error
}

type SshCommandRunner interface {
	RunCommand(name string, arg ...string) error
}

type sshCommand struct {
	IP          string
	BastionIP   string
	Certificate string
	Config      SshConfig
	CmdRunner   SshCommandRunner
}

func NewSshCommand(IP string, bastionIP string, certificate string, sshConfig SshConfig, cmdRunner SshCommandRunner) SshCommand {
	return &sshCommand{
		IP:          IP,
		BastionIP:   bastionIP,
		Certificate: certificate,
		Config:      sshConfig,
		CmdRunner:   cmdRunner,
	}
}

func (ssh *sshCommand) Run(command string) error {
	err := ioutil.WriteFile(ssh.Config.GetCertPath(), []byte(ssh.Certificate), 0644)
	if err != nil {
		return err
	}

	sshArgs := [...]string{
		"-t", "-t",
		"-oStrictHostKeyChecking=no",
		"-oLogLevel=QUIET",
		"-oUserKnownHostsFile=/dev/null",
		"-oServerAliveInterval=60",
		"-oServerAliveCountMax=720", // 12 hours
		fmt.Sprintf("-oProxyCommand=ssh -W %%h:%%p -i %s hopper@%s", ssh.Config.GetPrivateKeyPath(), ssh.BastionIP),
		"-i", ssh.Config.GetPrivateKeyPath(),
		fmt.Sprintf("ec2-user@%s", ssh.IP),
		command,
	}
	if ssh.Config.IsDebug() {
		fmt.Printf("ssh %s\n", strings.Join(sshArgs[:], " "))
	}

	return ssh.CmdRunner.RunCommand("ssh", sshArgs[:]...)
}
