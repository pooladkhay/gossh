package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

type sessionOpts struct {
	Name                    string
	Host                    string
	Remote                  string
	SSHPort                 string
	User                    string
	Password                string
	Ctx                     context.Context
	CtxCancel               context.CancelFunc
	localPortForwardEnabled bool
	portPairsMap            map[string]string
}

func (ss *sessionOpts) start() error {
	config := &ssh.ClientConfig{
		User: ss.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(ss.Password),
		},
		Timeout:           4 * time.Second,
		HostKeyCallback:   hostKeyCallback(),
		HostKeyAlgorithms: []string{ssh.KeyAlgoED25519},
	}

	fmt.Print("\033[2J")
	fmt.Print("\033[H")
	fmt.Printf("Conntecting to %s on port %s...\n", ss.Remote, ss.SSHPort)

	sshConn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", ss.Remote, ss.SSHPort), config)
	if err != nil {
		return fmt.Errorf("cannot connect to %s:%s. %s", ss.Remote, ss.SSHPort, err)
	}
	defer sshConn.Close()

	go ss.startPortForwarding(sshConn)

	session, err := sshConn.NewSession()
	if err != nil {
		return fmt.Errorf("cannot open new session: %s", err)
	}
	defer session.Close()

	go func() {
		<-ss.Ctx.Done()
		sshConn.Close()
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

func (ss *sessionOpts) startPortForwarding(sshConn *ssh.Client) {
	if ss.localPortForwardEnabled {
		if len(ss.portPairsMap) > 0 {
			for localPort := range ss.portPairsMap {

				localListener, err := net.Listen("tcp", fmt.Sprintf(":%s", localPort))
				if err != nil {
					fmt.Print("\033[95m")
					fmt.Print("\r\nMessage from local machine:")
					fmt.Print("\033[0m")
					fmt.Printf("\r\nfailed to listen on port %s on local machine: %s\r\n", localPort, err.Error())
					fmt.Print("(Press enter to ignore this message)")
					continue
				}

				go func() {
					for {
						lps := strings.Split(localListener.Addr().String(), ":")
						rp := ss.portPairsMap[lps[len(lps)-1]]

						localConn, err := localListener.Accept()
						if err != nil {
							fmt.Print("\033[95m")
							fmt.Print("\r\nMessage from local machine:")
							fmt.Print("\033[0m")
							fmt.Printf("\r\nfailed to accept connection on port %s on local machine: %s\r\n", lps[len(lps)-1], err.Error())
							fmt.Print("(Press enter to ignore this message)")
							continue
						}

						remoteConn, err := sshConn.Dial("tcp", fmt.Sprintf("127.0.0.1:%s", rp))
						if err != nil {
							fmt.Print("\033[95m")
							fmt.Print("\r\nMessage from remote machine:")
							fmt.Print("\033[0m")
							fmt.Printf("\r\nfailed to dial port %s on remote machine: %s\r\n", rp, err.Error())
							fmt.Print("(Press enter to ignore this message)")
							localConn.Close()
							continue
						}

						go func(lc, rc net.Conn) {
							defer rc.Close()
							defer lc.Close()
							_, err = io.Copy(rc, lc)
							if err != nil {
								fmt.Print("\033[95m")
								fmt.Print("\r\nMessage from local machine:")
								fmt.Print("\033[0m")
								fmt.Printf("\r\nio.Copy failed: %s\r\n", err.Error())
								fmt.Print("(Press enter to ignore this message)")
								return
							}
						}(localConn, remoteConn)
						go func(lc, rc net.Conn) {
							defer rc.Close()
							defer lc.Close()
							_, err = io.Copy(lc, rc)
							if err != nil {
								fmt.Print("\033[95m")
								fmt.Print("\r\nMessage from local machine:")
								fmt.Print("\033[0m")
								fmt.Printf("\r\nio.Copy failed: %s\r\n", err.Error())
								fmt.Print("(Press enter to ignore this message)")
								return
							}
						}(localConn, remoteConn)
					}
				}()
			}
		}
	}
}
