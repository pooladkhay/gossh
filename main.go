package main

import (
	"log"
	"os"

	"gopkg.in/ini.v1"
)

const VERSION = "0.1.0"

var cfg *ini.File
var srvFile = "/usr/local/etc/gossh/servers.ini"

func init() {
	srvDir := "/usr/local/etc/gossh"

	// Create directory "/usr/local/etc/gossh" if not exists
	_, err := os.Stat(srvDir)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(srvDir, 0755)
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
	iniFile, err := ini.Load(srvFile)
	if err != nil {
		log.Fatalf("Failed to read servers.ini file: %v\n", err)
	}
	cfg = iniFile
}

func main() {
	if len(os.Args) > 1 && os.Args[1] != "" {
		switch os.Args[1] {
		case "list":
			list()
		case "add":
			add()
		case "connect":
			connect()
		case "delete":
			delete()
		case "-h", "--help":
			help()
		default:
			helpErr()
		}
	} else {
		helpErr()
	}
}
