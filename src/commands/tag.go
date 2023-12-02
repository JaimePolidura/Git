package commands

import (
	"fmt"
	"git/src/objects"
	"git/src/repository"
	"git/src/utils"
	"strings"
)

// Tag List tags: main.go tag
// Tag Create tag object: main.go tag -a <name> [object default: HEAD]
func Tag(args []string) {
	if len(args) < 2 {
		utils.ExitError("Invalid arguments: tag [-a] [NANE]")
	}

	currentRepository, _, err := repository.FindCurrentRepository(utils.CurrentPath())
	if err != nil {
		utils.ExitError(err.Error())
	}

	if len(args) == 2 { //main.go tag -> List all tags
		listTags(currentRepository)
	} else {
		createTagObject := args[2] == "-a"
		tagName := args[3]
		object := extractObjectFromArgs(args)

		createTag(currentRepository, tagName, object, createTagObject)
	}
}

func createTag(repository *repository.Repository, name string, refValue string, createTagObject bool) {
	resolvedHashRefValue := repository.ResolveObjectName(refValue)
	tagNamePath := utils.Path("tags", name)

	if createTagObject {
		tagObject := objects.TagObject{
			Object:    objects.Object{Type: objects.TAG},
			ObjectTag: resolvedHashRefValue,
			Tag:       tagNamePath,
			Tagger:    "Jaime Polidura <jaime.polidura@gmail.com>",
		}

		if shaObjectTagWritten, err := repository.WriteObject(tagObject); err == nil {
			repository.WriteRef(objects.Reference{NamePath: tagNamePath, Value: shaObjectTagWritten})
		} else {
			utils.ExitError("Cannot create tag: " + err.Error())
		}
	} else {
		repository.WriteRef(objects.Reference{NamePath: tagNamePath, Value: resolvedHashRefValue})
	}
}

func listTags(repository *repository.Repository) {
	tags, err := repository.GetAllRefs()
	if err != nil {
		utils.ExitError("Cannot get references: " + err.Error())
	}

	for refPath, refValue := range tags {
		if strings.Contains(refPath, "tags") {
			fmt.Println(refValue.NamePath + " => " + refValue.Value)
		}
	}
}

func extractObjectFromArgs(args []string) string {
	if len(args) == 5 {
		return args[4]
	} else {
		return "HEAD"
	}
}
