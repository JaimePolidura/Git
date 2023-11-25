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
	case "cat-file":
		commands.CatFile(os.Args)
	case "hash-object":
		commands.HashObject(os.Args)
	}
}
