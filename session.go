package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

type serverOpts struct {
	Name     string
	Host     string
	Remote   string
	Port     string
	User     string
	Password string
	Ctx      context.Context
}

func (c *serverOpts) start() error {
	config := &ssh.ClientConfig{
		User: c.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(c.Password),
		},
		Timeout:           4 * time.Second,
		HostKeyCallback:   hostKeyCallback(),
		HostKeyAlgorithms: []string{ssh.KeyAlgoED25519},
	}

	fmt.Print("\033[2J")
	fmt.Print("\033[H")
	fmt.Printf("Conntecting to %s on port %s...\n", c.Remote, c.Port)

	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", c.Remote, c.Port), config)
	if err != nil {
		return fmt.Errorf("cannot connect to %s:%s. %s", c.Remote, c.Port, err)
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		return fmt.Errorf("cannot open new session: %s", err)
	}
	defer session.Close()

	go func() {
		<-c.Ctx.Done()
		conn.Close()
	}()

	fd := int(os.Stdin.Fd())
	state, err := term.MakeRaw(fd)
	if err != nil {
		return fmt.Errorf("terminal make raw failed: %s", err)
	}
	defer term.Restore(fd, state)

	w, h, err := term.GetSize(fd)
	if err != nil {
		return fmt.Errorf("terminal get size failed: %s", err)
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	if err := session.RequestPty("xterm-256color", h, w, modes); err != nil {
		return fmt.Errorf("session xterm failed: %s", err)
	}

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	if err := session.Shell(); err != nil {
		return fmt.Errorf("session shell failed: %s", err)
	}

	fmt.Print("\033[2J")
	fmt.Print("\033[H")
	if err := session.Wait(); err != nil {
		if e, ok := err.(*ssh.ExitError); ok {
			switch e.ExitStatus() {
			case 130:
				return nil
			}
		}
		return fmt.Errorf("ssh failed: %s", err)
	}
	return nil
}
