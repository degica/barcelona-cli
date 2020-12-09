package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/degica/barcelona-cli/config"
)

type SSH struct {
	IP          string
	BastionIP   string
	Certificate string
}

func (ssh *SSH) Run(command string) error {
	err := ioutil.WriteFile(config.Get().GetCertPath(), []byte(ssh.Certificate), 0644)
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
		fmt.Sprintf("-oProxyCommand=ssh -W %%h:%%p -i %s hopper@%s", config.Get().GetPrivateKeyPath(), ssh.BastionIP),
		"-i", config.Get().GetPrivateKeyPath(),
		fmt.Sprintf("ec2-user@%s", ssh.IP),
		command,
	}
	if config.Debug {
		fmt.Printf("ssh %s\n", strings.Join(sshArgs[:], " "))
	}

	cmd := exec.Command("ssh", sshArgs[:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if cmd.Run() != nil {
		return err
	}
	return nil
}
