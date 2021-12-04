package main

import (
	"fmt"
)

func list() {
	names := cfg.SectionStrings()
	if len(names) <= 1 {
		fmt.Println("Server list is empty.\nUse -h or --help to get started!")
		return
	}
	fmt.Println("\nAvailable Servers:")
	for _, name := range names {
		if name != "DEFAULT" {
			fmt.Println("-", name)
		}
	}
}
