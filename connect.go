package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func getConfig(value string) string {
	return cfg.Section(os.Args[2]).Key(value).String()
}

func connect() {
	if len(os.Args) <= 2 {
		fmt.Println("\nserver name must be provided as second argument")
		os.Exit(0)
	}
	if os.Args[2] == "" {
		fmt.Println("\nserver name must be provided as second argument")
		os.Exit(0)
	}
	if cfg.Section(os.Args[2]).Key("host").String() == "" {
		fmt.Printf("\nserver \"%s\" not found\n", os.Args[2])
		os.Exit(0)
	}

	var password string

	if getConfig("encrypted") == "1" {
		key := argAfterFlag(os.Args, "-e")
		p, err := decryptPass([]byte(key), getConfig("password"))
		if err != nil {
			log.Fatalln("Error while dencrypting password:", err)
		}
		password = string(p)
	} else {
		password = getConfig("password")
	}

	ctx, cancel := context.WithCancel(context.Background())
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)

	server := serverOpts{
		Remote:   getConfig("remote"),
		Port:     getConfig("port"),
		User:     getConfig("user"),
		Password: password,
		Ctx:      ctx,
	}

	go func() {
		if err := server.start(); err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		cancel()
	}()

	select {
	case <-sig:
		cancel()
	case <-ctx.Done():
	}
}
