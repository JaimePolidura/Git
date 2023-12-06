package commands

import (
	"bufio"
	"fmt"
	"git/src/objects"
	"git/src/repository"
	"git/src/utils"
	"os"
)

// HashObject Takes file and creates a blob object in .git. It returns the sha
// HashObject Args: main.go hash-object -t blob -w <blob path>
func HashObject(args []string) {
	if len(args) != 6 {
		fmt.Fprintf(os.Stderr, "Invalid args. Use: hash-object -t blob -w <blob path>\n")
		os.Exit(1)
	}

	currentRepository, _, err := repository.FindCurrentRepository(utils.CurrentPath())
	utils.Check(err, "fatal: not a git repository (or any of the parent directories): .git")

	filePath := args[5]
	file, err := os.Open(filePath)
	utils.Check(err, "Error while opening the file")
	defer file.Close()

	buffer := bufio.NewReader(file)
	bytesFromFile, _ := buffer.ReadBytes('\x00')

	object := objects.CreateBlobObject(bytesFromFile)

	sha, err := currentRepository.WriteObject(object)
	if err != nil {
		utils.ExitError("Error while writting object")
	}

	_, _ = os.Stdout.Write([]byte(sha))
}
