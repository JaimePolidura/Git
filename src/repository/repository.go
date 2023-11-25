package repository

import (
	"git/src/utils"
	"gopkg.in/ini.v1"
	"os"
	"path/filepath"
	"strings"
)

type Repository struct {
	WorkTree string
	GitDir   string
	Config   *ini.File
}

func FindCurrentRepository(currentPath string) *Repository {
	paths := strings.Split(currentPath, string(filepath.Separator))

	for i := 0; i < len(paths); i++ {
		actualPath := utils.JoinStrings(paths[:len(paths)-i])
		if _, err := os.Open(utils.Path(actualPath, ".git")); err == nil {
			return CreateRepositoryObject(actualPath)
		}
	}

	utils.ExitError("Cannot find git repository")

	return nil
}

func CreateRepositoryObject(path string) *Repository {
	workTree := path
	gitDir := utils.Path(path, ".git")

	gitPathFile, err := os.Open(gitDir)
	defer gitPathFile.Close()
	utils.Check(err, "Cannot open .git directory")

	gitPathFileStat, err := gitPathFile.Stat()
	utils.Check(err, "Cannot call Stat of .git directory")

	if !gitPathFileStat.IsDir() {
		utils.ExitError(".git is not a directory")
	}

	configFile, err := ini.Load(utils.Path(gitDir, "config"))
	utils.Check(err, "Cannot open config ini file in .git")

	version, err := configFile.Section("core").Key("repositoryformatversion").Int()

	if err != nil || version != 0 {
		utils.ExitError("Cannot get version in config file in .git")
	}

	return &Repository{
		WorkTree: workTree,
		GitDir:   gitDir,
		Config:   configFile,
	}
}

func InitializeRepository(workTreePath string) *Repository {
	gitDir := utils.Path(workTreePath, ".git")

	workDirFile, err := os.Open(workTreePath)
	utils.Check(err, "Cannot open "+workTreePath)
	defer workDirFile.Close()

	stat, err := workDirFile.Stat()
	utils.Check(err, "Cannot get stat from "+workTreePath)

	if !stat.IsDir() {
		utils.ExitError(workTreePath + "is not a directory")
	}

	utils.CreateDirIfNotExists(workTreePath, ".git")
	utils.CreateDirIfNotExists(gitDir, "branches")
	utils.CreateDirIfNotExists(gitDir, "refs")
	utils.CreateDirIfNotExists(utils.Path(gitDir, "refs"), "heads")
	utils.CreateDirIfNotExists(utils.Path(gitDir, "refs"), "tags")

	utils.CreateFileIfNotExistsWithContent(gitDir, "description", "Unnamed repository; edit this file 'description' to name the repository.\n")
	utils.CreateFileIfNotExistsWithContent(gitDir, "HEAD", "ref: refs/heads/master\n")

	utils.CreateFileIfNotExists(gitDir, "config")
	config, err := ini.Load(utils.Path(gitDir, "config"))
	utils.Check(err, "Cannot open config in .git")

	addDefaultConfigToIniFile(config)

	return &Repository{
		WorkTree: workTreePath,
		GitDir:   gitDir,
		Config:   config,
	}
}

func addDefaultConfigToIniFile(iniFile *ini.File) {
	section, err := iniFile.NewSection("core")
	utils.Check(err, "Cannot create core section in ini file")

	section.NewKey("repositoryformatversion", "0")
	section.NewKey("filemode", "fals")
	section.NewKey("bare", "false")

	iniFile.SaveTo("./.git/config")
}
