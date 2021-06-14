package main

import (
	"fmt"
	"log"
	"os"
)

func delete() {
	if len(os.Args) <= 2 {
		fmt.Println("server name must be provided as second argument")
		os.Exit(0)
	}
	if os.Args[2] == "" {
		fmt.Println("server name must be provided as second argument")
		os.Exit(0)
	}
	if cfg.Section(os.Args[2]).Key("host").String() == "" {
		fmt.Printf("server \"%s\" not found\n", os.Args[1])
		os.Exit(0)
	}

	sec := os.Args[2]
	cfg.DeleteSection(sec)
	err := cfg.SaveTo(srvFile)
	if err != nil {
		log.Fatalln("failed to add new server:", err)
	}
	fmt.Printf("Server %s deleted successfully\n", sec)
}
