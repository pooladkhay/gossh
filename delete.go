package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var cmdDelete = &cobra.Command{
	Use:                   "delete [server's name]",
	Example:               "gossh delete server_name",
	Short:                 "Deletes the specified server from server's list",
	Args:                  cobra.MinimumNArgs(1),
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		delete(args[0])
	},
}

func delete(srv string) {
	if cfg.Section(srv).Key("host").String() == "" {
		fmt.Printf("Server '%s' not found.\n", srv)
		os.Exit(0)
	}

	cfg.DeleteSection(srv)
	err := cfg.SaveTo(srvFile)
	if err != nil {
		log.Fatalln("failed to delete the server:", err)
	}
	fmt.Printf("Server '%s' deleted successfully.\n", srv)
}
