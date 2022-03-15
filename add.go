package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/spf13/cobra"
)

var cmdAdd = &cobra.Command{
	Use:                   "add [FLAGS]",
	Example:               "Without password encryption:\ngossh add -n my-server -a myserver.com -t 7121 -u root -p strong@pass\nWith password encryption:\ngossh add -n my-server -a myserver.com -t 7121 -u root -p strong@pass -k strong@key",
	Short:                 "Adds a new server to the list",
	Args:                  cobra.MinimumNArgs(0),
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		address, _ := cmd.Flags().GetString("address")
		port, _ := cmd.Flags().GetString("port")
		user, _ := cmd.Flags().GetString("user")
		password, _ := cmd.Flags().GetString("password")
		key, _ := cmd.Flags().GetString("key")
		if name != "" && address != "" && user != "" && password != "" {
			add(name, address, port, user, password, key)
		} else {
			fmt.Println("All flags except port and key are mandatory and must be provided.\nUse 'gossh add --help' for more information.")
			os.Exit(0)
		}
	},
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

func add(name, address, port, user, password, key string) {
	// srv := new(sessionOpts)
	// check if name is unique or not
	if cfg.Section(name).HasKey("host") {
		fmt.Println("server name must be unique")
		os.Exit(0)
	}

	// create a section in config file
	sec, _ := cfg.NewSection(name)

	// parsing password
	if key != "" {
		p, err := encryptPass([]byte(key), password)
		if err != nil {
			log.Fatalln("Error while encrypting password:", err)
		}
		sec.NewKey("password", string(p))
		sec.NewKey("encrypted", "1")
	} else {
		sec.NewKey("password", password)
		sec.NewKey("encrypted", "0")
	}

	// save to servers.ini
	sec.NewKey("host", address)
	sec.NewKey("remote", resolveDNS(address))
	sec.NewKey("port", port)
	sec.NewKey("user", user)

	err := cfg.SaveTo(cfgFileAddr)
	if err != nil {
		log.Fatalln("failed to add new server:", err)
	}
	fmt.Printf("Name: %s\nHost: %s\nPort: %s\nUser: %s\nSaved successfully.\n", name, address, port, user)
	if key == "" {
		fmt.Printf("You can may now connect using:\n $ gossh connect %s\n", name)
	} else {
		fmt.Printf("You can may now connect using:\n $ gossh connect %s -k [key]\n", name)
	}
}
