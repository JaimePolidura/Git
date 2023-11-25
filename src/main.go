package main

import (
	"fmt"
	"git/src/commands"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Invalid argument use\n")
		os.Exit(1)
	}

	switch command := os.Args[1]; command {
	case "init":
		commands.Init()
	}
}
