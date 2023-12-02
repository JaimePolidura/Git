package ignore

import (
	"errors"
	"git/src/utils"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type GitIgnore struct {
	ignoredRules []string
}

func (i *GitIgnore) IsIgnored(fileName string) (bool, error) {
	for _, ignoredRule := range i.ignoredRules {
		matchesIgnoreRule, err := filepath.Match(ignoredRule, fileName)
		if err != nil {
			return false, errors.New("Error while matching gitignore with " + fileName + " with pattern" + ignoredRule)
		}

		if matchesIgnoreRule {
			return true, nil
		}
	}

	return false, nil
}

func Deserialize(reader io.Reader) (GitIgnore, error) {
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return GitIgnore{}, err
	}

	actualOffset := 0
	ignored := make([]string, 0)

	for actualOffset < len(bytes) {
		lineBytes, newOffset, err := utils.ReadUntil(bytes, actualOffset, '\n')

		if err != nil {
			return GitIgnore{}, err
		}

		lineString := strings.Trim(string(lineBytes), " ")
		actualOffset = newOffset

		if strings.HasPrefix(lineString, "#") { //Commented
			continue
		}
		if lineString == "" {
			continue
		}

		ignored = append(ignored, lineString)
	}

	return GitIgnore{
		ignoredRules: ignored,
	}, nil

}
