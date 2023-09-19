package workdir

import (
	"errors"
	cp "github.com/otiai10/copy"
	"os"
	"strconv"
	"strings"
	"time"
)

type WorkDir struct {
	RootDirectory string
}

func InitEmptyWorkDir() *WorkDir {
	rootDirectory := "new_project/"
	if _, err := os.Stat(rootDirectory); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(rootDirectory, 0666)
		if err != nil {
			panic(err)
		}
	}
	return &WorkDir{RootDirectory: rootDirectory}
}

func (wd WorkDir) Clone() *WorkDir {
	t := time.Now()
	rootDirectory := "clones/" + strconv.FormatInt(t.Unix(), 10) + "/"
	err := os.MkdirAll(rootDirectory, 0666)
	if err != nil {
		panic(err)
	}
	cp.Copy(wd.RootDirectory, rootDirectory)
	return &WorkDir{RootDirectory: rootDirectory}
}

func (wd WorkDir) CreateFile(filename string) error {
	file, err := os.Create(wd.RootDirectory + filename)
	if err != nil {
		return err
	}
	err = file.Close()
	if err != nil {
		return err
	}
	return nil
}

func (wd WorkDir) CreateDir(path string) error {
	fullPath := wd.RootDirectory + path
	if _, err := os.Stat(fullPath); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(fullPath, 0666)
		if err != nil {
			return err
		}
	}
	return nil
}

func (wd WorkDir) WriteToFile(filename string, content string) error {
	filename = wd.RootDirectory + filename
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		return err
	} else {
		err := os.WriteFile(filename, []byte(content), 0666)
		if err != nil {
			return err
		}
	}
	return nil
}

func (wd WorkDir) AppendToFile(filename string, content string) error {
	beforeAppend, err := wd.CatFile(filename)
	if err != nil {
		return err
	}
	err = wd.WriteToFile(filename, beforeAppend+content)
	if err != nil {
		return err
	}
	return nil
}

func (wd WorkDir) CatFile(filename string) (string, error) {
	filename = wd.RootDirectory + filename
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		return "", err
	} else {
		dat, err := os.ReadFile(filename)
		if err != nil {
			return "", err
		}
		return string(dat), nil
	}
}

func (wd WorkDir) ListFilesRoot() []string {
	files, err := getListOfAllFile(wd.RootDirectory)
	if err != nil {
		panic(err)
	}
	var newFiles []string
	for _, item := range files {
		newFiles = append(newFiles, strings.TrimPrefix(item, wd.RootDirectory))
	}
	filesArr = []string{}
	return newFiles
}

func (wd WorkDir) ListFilesIn(path string) ([]string, error) {
	files, err := getListOfAllFile(wd.RootDirectory + path + "/")
	if err != nil {
		panic(err)
	}
	var newFiles []string
	for _, item := range files {
		newFiles = append(newFiles, strings.TrimLeft(item, wd.RootDirectory))
	}
	filesArr = []string{}
	return newFiles, nil
}

var filesArr []string

func getListOfAllFile(path string) ([]string, error) {
	dir, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func(dir *os.File) {
		err := dir.Close()
		if err != nil {

		}
	}(dir)

	files, err := dir.Readdir(-1)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if file.IsDir() {
			_, err := getListOfAllFile(path + file.Name() + "/")
			if err != nil {
				return nil, err
			}
		} else {
			filesArr = append(filesArr, path+file.Name())
		}
	}
	return filesArr, nil
}

func GetModTimeOfFile(filename string) int64 {
	stat, err := os.Stat(filename)
	if err != nil {
		panic(err)
	}
	return stat.ModTime().UnixNano()
}
