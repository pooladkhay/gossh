package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

func createKnownHosts() {
	f, err := os.OpenFile(filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"), os.O_CREATE, 0600)
	if err != nil {
		fmt.Printf("error creating known_hosts: %s\n", err)
		os.Exit(1)
	}
	defer f.Close()
}

func checkKnownHosts() ssh.HostKeyCallback {
	createKnownHosts()
	kh, err := knownhosts.New(filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"))
	if err != nil {
		fmt.Printf("error check known_hosts: %s\n", err)
		os.Exit(1)
	}
	return kh
}

func createHostKey(remote net.Addr, pubKey ssh.PublicKey) string {
	kh := knownhosts.Normalize(remote.String())
	return fmt.Sprintln(knownhosts.Line([]string{kh}, pubKey))
}

func addHostKey(host string, remote net.Addr, pubKey ssh.PublicKey) error {
	khFilePath := filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts")

	f, fErr := os.OpenFile(khFilePath, os.O_APPEND|os.O_WRONLY, 0600)
	if fErr != nil {
		return fErr
	}
	defer f.Close()

	_, fileErr := f.WriteString(createHostKey(remote, pubKey))
	return fileErr
}

func hostKeyCallback() ssh.HostKeyCallback {
	var keyErr *knownhosts.KeyError
	return ssh.HostKeyCallback(func(host string, remote net.Addr, pubKey ssh.PublicKey) error {
		fmt.Print("\033[31m")
		kh := checkKnownHosts()
		hErr := kh(host, remote, pubKey)
		if errors.As(hErr, &keyErr) && len(keyErr.Want) > 0 {
			fmt.Printf("WARNING: %v is not a key of %s, either a MiTM attack or %s has reconfigured the host pub key.", pubKey, host, host)
			fmt.Print("\033[0m")
			fmt.Println("Exiting...")
			return keyErr
		} else if errors.As(hErr, &keyErr) && len(keyErr.Want) == 0 {
			fmt.Printf(
				"WARNING: %s is not trusted, adding this key: %s to known_hosts file.",
				strings.Fields(createHostKey(remote, pubKey))[0],
				strings.Fields(createHostKey(remote, pubKey))[2],
			)
			fmt.Print("\033[0m")
			fmt.Println("Now Connecting...")
			return addHostKey(host, remote, pubKey)
		}
		fmt.Print("\033[0m")
		return nil
	})
}
