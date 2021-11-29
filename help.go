package main

import (
	"fmt"
	"os"
)

func help() {
	fmt.Print(
		`
List all available servers:
$ gossh list

Add a new server:
$ gossh add -n [server name[no spaces]] -a [server address] -t (Optional)[port [default:22]] -u [user] -p [password] -e (Optional)[key to encrypt password with]

Connect to a server:
$ gossh connect [server name] [-e (Optional)[key to decrypt password]] [-f localPort:remotePort,...]

Delete a server (permanently):
$ gossh delete [server name]
`)
}

func helpErr() {
	fmt.Println("\nInvalid syntax!\nUse -h or --help for help.")
	os.Exit(0)
}
