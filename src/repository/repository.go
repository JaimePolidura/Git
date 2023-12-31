package repository

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"git/src/ignore"
	"git/src/index"
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

func (r *Repository) WriteObject(object *objects.Object) (string, error) {
	serializeData := object.Serialize()
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
	gitObject, err := r.ReadObject(hash, objects.TREE)
	if err != nil {
		return objects.TreeObject{}, err
	}

	return gitObject.SerializableGitObject.(objects.TreeObject), nil
}

func (r *Repository) ReadBlobObject(hash string) (objects.BlobObject, error) {
	gitObject, err := r.ReadObject(hash, objects.BLOB)
	if err != nil {
		return objects.BlobObject{}, err
	}

	return gitObject.SerializableGitObject.(objects.BlobObject), nil
}

func (r *Repository) ReadCommitObject(hash string) (objects.CommitObject, error) {
	gitObject, err := r.ReadObject(hash, objects.COMMIT)
	if err != nil {
		return objects.CommitObject{}, err
	}

	return gitObject.SerializableGitObject.(objects.CommitObject), nil
}

func (r *Repository) ReadObject(unresolvedHash string, reqType objects.ObjectType) (objects.Object, error) {
	if resolvedHash, _, err := r.ResolveObjectName(unresolvedHash, reqType); err == nil {
		return r.readObjectByResolvedName(resolvedHash)
	} else {
		return objects.Object{}, err
	}
}

func (r *Repository) readObjectByResolvedName(resolvedHash string) (objects.Object, error) {
	prefix, remainder := resolvedHash[:2], resolvedHash[2:]
	objectPath := utils.Paths(r.GitDir, "objects", prefix, remainder)
	objectFile, err := os.Open(objectPath)
	defer objectFile.Close()
	if err != nil {
		return objects.Object{}, errors.New("Cannot open object file: " + resolvedHash)
	}
	objectFileState, err := objectFile.Stat()
	if err != nil {
		return objects.Object{}, errors.New("Cannot get stat from object file: " + resolvedHash)
	}
	if objectFileState.IsDir() {
		return objects.Object{}, errors.New("Object file: " + resolvedHash + " cannot be a dir")
	}

	objectFileZlibReader, err := zlib.NewReader(objectFile)
	if err != nil {
		return objects.Object{}, err
	}
	defer objectFileZlibReader.Close()

	return objects.DeserializeObject(objectFileZlibReader)
}

func (r *Repository) ReadIndex() (*index.IndexObject, error) {
	utils.CreateFileIfNotExists(r.GitDir, "index")

	if file, err := os.Open(utils.Path(r.GitDir, "index")); err == nil {
		return index.Deserialize(file)
	} else {
		return nil, err
	}
}

func (r *Repository) WriteIndex(index *index.IndexObject) error {
	if err := os.Remove(utils.Paths(r.GitDir, "index")); err != nil {
		return err
	}

	if file, err := os.Create(utils.Paths(r.GitDir, "index")); err == nil {
		_, err := file.Write(index.Serialize())
		file.Close()
		return err
	} else {
		return err
	}
}

func (r *Repository) AbsolutePathToRepositoryPath(path string) string {
	return utils.RemovePrefix(path, r.WorkTree+"/")
}

func (r *Repository) GetPathFileInRepository(path string) string {
	isAbsolute := strings.HasPrefix(path, "/")
	if isAbsolute {
		return path
	} else {
		return utils.Paths(utils.CurrentPath(), path)
	}
}

func (r *Repository) WriteRef(reference objects.Reference) {
	utils.CreateFileIfNotExistsWithContent(utils.Paths(r.GitDir, "refs"), reference.NamePath, reference.Value+"\n")
}

func (r *Repository) IsIgnored(pathInRepository string) (bool, error) {
	if strings.Contains(pathInRepository, ".git/") {
		return true, nil
	}

	gitIgnores, err := r.readGitIgnores(pathInRepository)
	if err != nil {
		return false, err
	}
	if len(gitIgnores) == 0 {
		return false, nil
	}

	parent := filepath.Dir(pathInRepository)

	for {
		if gitIgnore, gitIgnoreExists := gitIgnores[parent]; gitIgnoreExists {
			matchesSomeIgnore, err := gitIgnore.IsIgnored(pathInRepository)

			if err != nil {
				return false, err
			}
			if matchesSomeIgnore {
				return true, nil
			}
		}
		if parent == "" {
			break
		}

		parent = filepath.Dir(parent)
	}

	return false, nil
}

func (r *Repository) readGitIgnores(pathInRepository string) (map[string]ignore.GitIgnore, error) {
	index, err := r.ReadIndex()
	if err != nil {
		return nil, err
	}

	gitIgnores := make(map[string]ignore.GitIgnore)

	for _, entry := range index.Entries {
		if entry.FullPathName == ".gitignore" || strings.HasSuffix(entry.FullPathName, "/.gitignore") {
			file, err := os.Open(entry.FullPathName)
			if err != nil {
				return nil, err
			}
			gitIgnore, err := ignore.Deserialize(file)
			if err != nil {
				return nil, err
			}

			gitIgnores[filepath.Dir(entry.FullPathName)] = gitIgnore
		}
	}

	return gitIgnores, nil
}

func (r *Repository) ResolveRef(namePath string) (objects.Reference, error) {
	return r.resolveRefRecursive(namePath)
}

func (r *Repository) resolveRefRecursive(namePath string) (objects.Reference, error) {
	file, err := os.Open(utils.Path(r.GitDir, namePath))
	defer file.Close()
	if err != nil {
		return objects.Reference{}, NoCommits{}
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return objects.Reference{}, err
	}
	if len(bytes) == 0 {
		return objects.Reference{}, NoCommits{}
	}

	stringRef := string(bytes)
	isRef := strings.HasPrefix(stringRef, "ref: ")
	if isRef {
		nextRefPath := strings.Split(stringRef, " ")[1]
		return r.resolveRefRecursive(utils.SanitizePath(nextRefPath))
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

func (r *Repository) ResolveObjectName(name string, reqObjectType objects.ObjectType) (string, bool, error) {
	candidatesHash, isHead, err := r.getCandidatesResolveObjectName(name)
	if err != nil {
		return "", false, err
	}
	if len(candidatesHash) > 1 {
		utils.ExitError("Ambiguous reference ")
	}

	candidateHash := candidatesHash[0]

	for {
		candidateObject, err := r.readObjectByResolvedName(candidateHash)

		if err != nil {
			return "", false, err
		}

		if reqObjectType == objects.ANY || reqObjectType == candidateObject.Type {
			return candidateHash, isHead, nil
		}

		if candidateObject.Type == objects.TAG {
			candidateHash = candidateObject.SerializableGitObject.(objects.TagObject).ObjectTag
		} else if candidateObject.Type == objects.COMMIT && reqObjectType == objects.TREE {
			candidateHash = candidateObject.SerializableGitObject.(objects.CommitObject).Tree
		} else {
			return "", false, errors.New("cannot get type")
		}
	}
}

func (r *Repository) getCandidatesResolveObjectName(objectName string) ([]string, bool, error) {
	if strings.ToUpper(objectName) == "HEAD" {
		headHash, err := r.ResolveRef(objectName)
		return []string{headHash.Value}, true, err
	}

	candidatesHash := make([]string, 0)

	if utils.IsValidGitHash(objectName) {
		prefix := objectName[:2]
		remaining := objectName[2:]
		pathPrefix := utils.Paths(r.GitDir, "objects", prefix)
		if utils.CheckFileOrDirExists(pathPrefix) {
			dirs, _ := os.ReadDir(pathPrefix)
			for _, file := range dirs {
				if strings.HasPrefix(file.Name(), remaining) {
					candidatesHash = append(candidatesHash, prefix+file.Name())
				}
			}
		}
	}

	candidateIsHead := false

	if ref, err := r.ResolveRef("refs/tags/" + objectName); err == nil {
		candidatesHash = append(candidatesHash, ref.Value)
	}

	if ref, err := r.ResolveRef("refs/heads/" + objectName); err == nil {
		candidatesHash = append(candidatesHash, ref.Value)
		candidateIsHead = true
	}

	if len(candidatesHash) == 0 {
		return candidatesHash, false, errors.New("Cannot find object with name " + objectName)
	} else {
		return candidatesHash, candidateIsHead, nil
	}

}

func (r *Repository) WriteToHead(value string) error {
	file, err := os.OpenFile(utils.Paths(r.GitDir, "HEAD"), os.O_WRONLY, 0777)
	defer file.Close()
	utils.CheckError(err)
	_, err = file.Write([]byte(value))
	return err
}

func (r *Repository) GetActiveBranch() (_name string, _detached bool, _err error) {
	file, err := os.Open(utils.Path(r.GitDir, "HEAD"))
	if err != nil {
		return "", false, err
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return "", false, err
	}

	headContentString := string(bytes)

	if strings.HasPrefix(headContentString, "ref: refs/heads/") {
		return utils.SanitizePath(string(bytes[16:])), false, nil
	} else {
		return utils.SanitizePath(headContentString), true, nil
	}
}

func FindCurrentRepository(currentPath string) (*Repository, string, error) {
	paths := strings.Split(currentPath, string(filepath.Separator))

	for i := 0; i < len(paths); i++ {
		actualPath := "/" + utils.Paths(paths[:len(paths)-i]...)

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
	utils.CreateFileIfNotExists(utils.Paths(gitDir, "refs", "heads"), "master")

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
