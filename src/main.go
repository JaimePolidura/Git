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
	case "log":
		commands.Log(os.Args)
	case "ls-tree":
		commands.LsTree(os.Args)
	case "checkout":
		commands.Checkout(os.Args)
	case "tag":
		commands.Tag(os.Args)
	case "ls-files":
		commands.LsFiles(os.Args)
	case "add":
		commands.Add(os.Args)
	case "commit":
		commands.Commit(os.Args)
	default:
		panic("Unknown command")
	}
}
