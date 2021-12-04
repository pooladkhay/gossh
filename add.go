package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func sliceContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func argAfterFlag(s []string, flag string) string {
	if !sliceContains(os.Args, flag) {
		fmt.Printf("%s must be provided\n", flag)
		os.Exit(0)
	}
	for k, v := range s {
		if flag == v {
			if len(s) > k+1 {
				if strings.HasPrefix(s[k+1], "-") {
					helpErr()
				}
				return s[k+1]
			} else {
				helpErr()
			}
		}
	}
	helpErr()
	return ""
}

func resolveDNS(addr string) string {
	fmt.Printf("Resolving IP address for %s...\n", addr)
	ip, err := net.ResolveIPAddr("ip", addr)
	if err != nil {
		fmt.Printf("Failed to resolve IP address for %s\n", err)
		os.Exit(1)
	}
	return ip.String()
}

func add() {
	srv := new(sessionOpts)
	// parsing server name
	arg := argAfterFlag(os.Args, "-n")
	// check if name is unique or not
	if cfg.Section(arg).HasKey("host") {
		fmt.Println("-n [server name[no spaces]] must be unique")
		os.Exit(0)
	}
	srv.Name = arg

	// create a section in config file
	sec, _ := cfg.NewSection(srv.Name)

	// parsing address
	srv.Remote = resolveDNS(argAfterFlag(os.Args, "-a"))
	srv.Host = argAfterFlag(os.Args, "-a")

	// parsing port
	if !sliceContains(os.Args, "-t") {
		srv.SSHPort = "22"
	} else {
		srv.SSHPort = argAfterFlag(os.Args, "-t")
	}

	// parsing user
	srv.User = argAfterFlag(os.Args, "-u")

	// parsing password
	if sliceContains(os.Args, "-e") {
		rawPass := argAfterFlag(os.Args, "-p")
		key := argAfterFlag(os.Args, "-e")
		p, err := encryptPass([]byte(key), rawPass)
		if err != nil {
			log.Fatalln("Error while encrypting password:", err)
		}
		sec.NewKey("password", string(p))
		sec.NewKey("encrypted", "1")
	} else {
		srv.Password = argAfterFlag(os.Args, "-p")
		sec.NewKey("password", srv.Password)
		sec.NewKey("encrypted", "0")
	}

	// save to servers.ini
	sec.NewKey("host", srv.Host)
	sec.NewKey("remote", srv.Remote)
	sec.NewKey("port", srv.SSHPort)
	sec.NewKey("user", srv.User)

	err := cfg.SaveTo(srvFile)
	if err != nil {
		log.Fatalln("failed to add new server:", err)
	}

	fmt.Printf("-Name: %s\n-Host: %s\n-Port: %s\n-User: %s\n\nSaved successfully.\n", srv.Name, srv.Host, srv.SSHPort, srv.User)
	fmt.Printf("You can may now connect using:\n $ gossh connect %s\n", srv.Name)
}
