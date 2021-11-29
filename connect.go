package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
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

	session := sessionOpts{
		Remote:                  getConfig("remote"),
		SSHPort:                 getConfig("port"),
		User:                    getConfig("user"),
		Password:                password,
		Ctx:                     ctx,
		CtxCancel:               cancel,
		localPortForwardEnabled: false,
		portPairsMap:            make(map[string]string),
	}

	// ["gossh", "connect", "dev-pbx", "-f", "3306:3306,5038:5038,..."]
	//	 	 0			  1			 2			3			      4
	if len(os.Args) >= 5 {
		if os.Args[3] == "-f" {
			if strings.Contains(os.Args[4], ",") {
				allPorts := strings.Split(os.Args[4], ",")
				for _, lrPorts := range allPorts {
					if strings.Contains(lrPorts, ":") {
						ports := strings.Split(lrPorts, ":")
						if len(ports) == 2 {
							if ports[0] != "" && ports[1] != "" {
								session.localPortForwardEnabled = true
								session.portPairsMap[ports[0]] = ports[1]
							} else {
								helpErr()
							}
						} else {
							helpErr()
						}
					} else {
						helpErr()
					}
				}
			} else {
				if strings.Contains(os.Args[4], ":") {
					ports := strings.Split(os.Args[4], ":")
					if len(ports) == 2 {
						if ports[0] != "" && ports[1] != "" {
							session.localPortForwardEnabled = true
							session.portPairsMap[ports[0]] = ports[1]
						} else {
							helpErr()
						}
					} else {
						helpErr()
					}
				} else {
					helpErr()
				}
			}
		} else {
			helpErr()
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
