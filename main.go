package main

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/ini.v1"
)

const VERSION = "0.2.0"

var cfg *ini.File
var srvFile = "/usr/local/etc/gossh/servers.ini"

func init() {
	srvDir := "/usr/local/etc/gossh"

	// Create directory "/usr/local/etc/gossh" if not exists
	_, err := os.Stat(srvDir)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(srvDir, os.ModePerm)
		if errDir != nil {
			log.Fatalln(err)
		}
	}

	// Create file "/usr/local/etc/gossh/servers.ini" if not exists
	_, err = os.Stat(srvFile)
	if os.IsNotExist(err) {
		iniFile, err := os.Create(srvFile)
		if err != nil {
			log.Fatalf("Failed to create server.ini file: %s\n", err)
		}
		// iniFile.Chmod(0664)
		iniFile.Close()
	}

	// Load config file. you can backup this file and put it on another machine.
	// That's actually why I made this CLI App :)
	iniOpts := ini.LoadOptions{
		SpaceBeforeInlineComment: true,
	}
	iniFile, err := ini.LoadSources(iniOpts, srvFile)
	if err != nil {
		log.Fatalf("Failed to read servers.ini file: %v\n", err)
	}
	cfg = iniFile
}

func main() {
	cmdConnect.PersistentFlags().StringP("forward-local", "l", "", "enable local port forwarding")
	cmdConnect.PersistentFlags().StringP("key", "k", "", "key to dencrypt password with (only if encrypted while adding)")

	cmdAdd.PersistentFlags().StringP("name", "n", "", "server's name")
	cmdAdd.PersistentFlags().StringP("address", "a", "", "server's address url")
	cmdAdd.PersistentFlags().StringP("port", "t", "22", "server's ssh port (optional)")
	cmdAdd.PersistentFlags().StringP("user", "u", "", "username")
	cmdAdd.PersistentFlags().StringP("password", "p", "", "password")
	cmdAdd.PersistentFlags().StringP("key", "k", "", "key to encrypt password with (optional)")

	var rootCmd = &cobra.Command{Use: "gossh", Version: VERSION}
	rootCmd.AddCommand(cmdConnect, cmdList, cmdDelete, cmdAdd)
	rootCmd.Execute()
}
