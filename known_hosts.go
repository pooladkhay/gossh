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
		kh := checkKnownHosts()
		hErr := kh(host, remote, pubKey)
		if errors.As(hErr, &keyErr) && len(keyErr.Want) > 0 {
			fmt.Print("\033[31m")
			fmt.Println("WARNING")
			fmt.Print("\033[0m")
			fmt.Printf("'%s' is not a key of '%s'.\nEither a MiTM attack or '%s' has reconfigured the host pub key.\n", strings.Fields(createHostKey(remote, pubKey))[2], host, host)
			fmt.Println("Please contact your system administrator.")
			fmt.Println("Exiting...")
			return keyErr
		} else if errors.As(hErr, &keyErr) && len(keyErr.Want) == 0 {
			fmt.Print("\033[31m")
			fmt.Println("WARNING")
			fmt.Print("\033[0m")
			fmt.Printf("The authenticity of host '%s' can't be established.\n", host)
			fmt.Printf("ED25519 key fingerprint is: %s\n", ssh.FingerprintSHA256(pubKey))
			fmt.Printf("Do you want to continue connecting and permanently add '%s' to 'known_hosts'?\n", host)
			fmt.Print("[yes/no]: ")
			var val string
			for {
				fmt.Scanln(&val)
				switch val {
				case "yes":
					fmt.Println("Connecting...")
					return addHostKey(host, remote, pubKey)
				case "no":
					return keyErr
				default:
					fmt.Print("Please type 'yes' or 'no': ")
				}
			}
		}
		fmt.Print("\033[0m")
		return nil
	})
}
