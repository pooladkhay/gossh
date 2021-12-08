package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
)

var cmdConnect = &cobra.Command{
	Use:                   "connect [server to connect to] [(Oprtional) -l local_port:remote_port]",
	Example:               "gossh connect server_name\ngossh connect server_name -l 3000:3001\ngossh connect server_name -l 3000:3000,4000:4000,5436:1234 ",
	Short:                 "Connects to a specific server",
	Args:                  cobra.MinimumNArgs(1),
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		ports, _ := cmd.Flags().GetString("forward-local")
		key, _ := cmd.Flags().GetString("key")
		connect(args[0], ports, key)
	},
}

func getConfig(srv, value string) string {
	return cfg.Section(srv).Key(value).String()
}

func connect(srv, ports, key string) {
	if cfg.Section(srv).Key("host").String() == "" {
		fmt.Printf("server \"%s\" not found\n", srv)
		os.Exit(0)
	}

	var password string

	if getConfig(srv, "encrypted") == "1" {
		if key == "" {
			fmt.Println("Password for this server is encrypted with a key.\nUse 'gossh connect --help' for more information.")
			os.Exit(0)
		}
		p, err := decryptPass([]byte(key), getConfig(srv, "password"))
		if err != nil {
			log.Fatalln("Error while dencrypting password:", err)
		}
		password = string(p)
	} else {
		password = getConfig(srv, "password")
	}

	ctx, cancel := context.WithCancel(context.Background())
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)

	session := sessionOpts{
		Remote:                  getConfig(srv, "remote"),
		SSHPort:                 getConfig(srv, "port"),
		User:                    getConfig(srv, "user"),
		Password:                password,
		Ctx:                     ctx,
		CtxCancel:               cancel,
		localPortForwardEnabled: false,
		portPairsMap:            make(map[string]string),
	}

	if strings.Contains(ports, ",") {
		allPorts := strings.Split(ports, ",")
		for _, lrPorts := range allPorts {
			if strings.Contains(lrPorts, ":") {
				prts := strings.Split(lrPorts, ":")
				if len(prts) == 2 {
					if prts[0] != "" && prts[1] != "" {
						session.localPortForwardEnabled = true
						session.portPairsMap[prts[0]] = prts[1]
					}
				}
			}
		}
	} else {
		if strings.Contains(ports, ":") {
			prts := strings.Split(ports, ":")
			if len(prts) == 2 {
				if prts[0] != "" && prts[1] != "" {
					session.localPortForwardEnabled = true
					session.portPairsMap[prts[0]] = prts[1]
				}
			}
		}
	}

	go func() {
		if err := session.start(); err != nil {
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
