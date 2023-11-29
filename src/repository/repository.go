package repository

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"git/src/objects"
	"git/src/utils"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/ini.v1"
)

type Repository struct {
	WorkTree string
	GitDir   string
	Config   *ini.File
}

func (r *Repository) WriteObject(object objects.GitObject) (string, error) {
	serializeData := objects.SerializeObject(object)
	sha1Hasher := sha1.New()
	sha1Hasher.Write(serializeData)
	shaHex := hex.EncodeToString(sha1Hasher.Sum(nil))
	prefix, remainder := shaHex[:2], shaHex[2:]
	objectPath := utils.Paths(r.GitDir, "objects", prefix, remainder)

	if err := os.MkdirAll(filepath.Dir(objectPath), os.ModePerm); err != nil {
		return "", err
	}

	file, err := os.Create(objectPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var compressedBuffer bytes.Buffer
	zlibWriter := zlib.NewWriter(&compressedBuffer)
	defer zlibWriter.Close()

	if _, err = zlibWriter.Write(serializeData); err != nil {
		return "", err
	}

	zlibWriter.Flush()

	if err := zlibWriter.Close(); err != nil {
		return "", err
	}

	if _, err = file.Write(compressedBuffer.Bytes()); err != nil {
		return "", err
	} else {
		return shaHex, nil
	}
}

func (r *Repository) ReadTreeObject(hash string) (objects.TreeObject, error) {
	gitObject, err := r.ReadObject(hash)
	if err != nil || gitObject.Type() != objects.BLOB {
		return objects.TreeObject{}, err
	}

	return gitObject.(objects.TreeObject), nil
}

func (r *Repository) ReadBlobObject(hash string) (objects.BlobObject, error) {
	gitObject, err := r.ReadObject(hash)
	if err != nil || gitObject.Type() != objects.BLOB {
		return objects.BlobObject{}, err
	}

	return gitObject.(objects.BlobObject), nil
}

func (r *Repository) ReadCommitObject(hash string) (objects.CommitObject, error) {
	gitObject, err := r.ReadObject(hash)
	if err != nil || gitObject.Type() != objects.COMMIT {
		return objects.CommitObject{}, err
	}

	return gitObject.(objects.CommitObject), nil
}

func (r *Repository) ReadObject(unformattedHash string) (objects.GitObject, error) {
	formattedHash := r.FormatObjectHash(unformattedHash)
	prefix, remainder := formattedHash[:2], formattedHash[2:]
	objectPath := utils.Paths(r.GitDir, "objects", prefix, remainder)
	objectFile, err := os.Open(objectPath)
	defer objectFile.Close()
	if err != nil {
		return nil, errors.New("Cannot open object file: " + formattedHash)
	}
	objectFileState, err := objectFile.Stat()
	if err != nil {
		return nil, errors.New("Cannot get stat from object file: " + formattedHash)
	}
	if objectFileState.IsDir() {
		return nil, errors.New("Object file: " + formattedHash + " cannot be a dir")
	}

	objectFileZlibReader, err := zlib.NewReader(objectFile)
	if err != nil {
		return nil, err
	}
	defer objectFileZlibReader.Close()

	return objects.DeserializeObject(objectFileZlibReader)
}

func (r *Repository) FormatObjectHash(hash string) string {
	return hash
}

func (r *Repository) WriteRef(reference objects.Reference) {
	utils.CreateFileIfNotExistsWithContent(utils.Paths(r.GitDir, "refs"), reference.NamePath, reference.Value+"\n")
}

func (r *Repository) ResolveRef(namePath string) (objects.Reference, error) {
	return r.resolveRefRecursive(namePath)
}

func (r *Repository) resolveRefRecursive(namePath string) (objects.Reference, error) {
	file, err := os.Open(r.GitDir + namePath)
	//This is normal in one specific case: we're looking for HEAD on a new repository with no commits
	if err != nil {
		return objects.Reference{}, nil
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return objects.Reference{}, err
	}

	stringRef := string(bytes)
	isRef := strings.HasPrefix(stringRef, "ref: ")
	if isRef {
		nextRefPath := strings.Split(stringRef, " ")[1]
		return r.resolveRefRecursive(nextRefPath)
	} else {
		return objects.Reference{NamePath: namePath, Value: stringRef}, nil
	}
}

func (r *Repository) GetAllRefs() (map[string]objects.Reference, error) {
	refsPath := utils.Paths(r.GitDir, "/refs")
	result := make(map[string]objects.Reference)

	err := r.readRefsRecursive(result, refsPath)

	return result, err
}

func (r *Repository) readRefsRecursive(result map[string]objects.Reference, dirPath string) error {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil
	}

	for _, file := range files {
		filePath := utils.Paths(dirPath, file.Name())

		if !file.IsDir() {
			if resolvedRef, err := r.ResolveRef(filePath); err != nil {
				result[filePath] = resolvedRef
			}
		} else {
			r.readRefsRecursive(result, filePath)
		}
	}

	return nil
}

func FindCurrentRepository(currentPath string) (*Repository, string, error) {
	paths := strings.Split(currentPath, string(filepath.Separator))

	for i := 0; i < len(paths); i++ {
		actualPath := utils.JoinStrings(paths[:len(paths)-i])
		if _, err := os.Open(utils.Path(actualPath, ".git")); err == nil {
			return CreateRepositoryObject(actualPath), actualPath, nil
		}
	}

	return nil, "", errors.New("fatal: not a git repository (or any of the parent directories): .git")
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
