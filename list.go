package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var cmdList = &cobra.Command{
	Use:                   "list",
	Short:                 "Lists all available servers",
	Args:                  cobra.MinimumNArgs(0),
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		list()
	},
}

func list() {
	names := cfg.SectionStrings()
	if len(names) <= 1 {
		fmt.Println("Server list is empty.\nUse 'gossh add -h' to get started!")
		return
	}
	fmt.Println("Available Servers:")
	for _, name := range names {
		if name != "DEFAULT" {
			fmt.Println("-", name)
		}
	}
}
