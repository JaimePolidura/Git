package commands

import (
	"bufio"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"git/src/objects"
	"git/src/repository"
	"git/src/utils"
	"os"
)

// HashObject Args: main.go hash-object -t blob -w <blob path>
func HashObject(args []string) {
	if len(args) != 7 {
		fmt.Fprintf(os.Stderr, "Invalid args. Use: hash-object -t blob -w <blob path>\n")
		os.Exit(1)
	}

	repository, err := repository.FindCurrentRepository(utils.CurrentPath())

	filePath := args[5]
	file, err := os.Open(filePath)
	utils.Check(err, "Error while opening the file")
	defer file.Close()

	objectType, err := objects.GetObjectTypeByString(args[3])
	if err != nil {
		utils.ExitError("Unknown type")
	}

	buffer := bufio.NewReader(file)
	bytesFromFile, err := buffer.ReadBytes('\x00')
	utils.Check(err, "Error while reading the bytes of the file")

	hasher := sha1.New()
	hasher.Write(bytesFromFile)
	shaHex := hex.EncodeToString(hasher.Sum(nil))

	object := &objects.Object{Type: objectType, Length: len(bytesFromFile), Data: bytesFromFile}

	sha, err := repository.WriteObject(object)
	if err != nil {
		utils.ExitError("Error while writting object")
	}

	os.Stdout.Write([]byte(sha))
}
