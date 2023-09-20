package commands

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"vc/workdir"
)

type VC struct {
	wd          *workdir.WorkDir
	status      *Status
	commits     map[int]*CommitContent
	fileHistory map[string]string
}

type CommitContent struct {
	Message string
	Files   map[string]string
}

type Status struct {
	ModifiedFiles []string
	StagedFiles   []string
}

func Init(wd *workdir.WorkDir) *VC {
	return &VC{
		wd: wd,
		status: &Status{
			ModifiedFiles: []string{},
			StagedFiles:   []string{},
		},
		commits:     make(map[int]*CommitContent),
		fileHistory: make(map[string]string),
	}
}

func (vc VC) GetWorkDir() *workdir.WorkDir {
	return vc.wd
}

func (vc VC) Status() *Status {
	if len(vc.fileHistory) == 0 {
		return vc.status
	}
	for _, file := range vc.wd.ListFilesRoot() {
		lastStatus, ok := vc.fileHistory[file]
		if ok {
			lastModified, _ := vc.wd.CatFile(file)
			if lastStatus != lastModified {
				vc.fileHistory[file] = lastModified
				vc.status.ModifiedFiles = append(vc.status.ModifiedFiles, file)
			}
		} else {
			vc.fileHistory[file], _ = vc.wd.CatFile(file)
			vc.status.ModifiedFiles = append(vc.status.ModifiedFiles, file)
		}
	}
	return vc.status
}

func (vc VC) AddAll() {
	files := vc.wd.ListFilesRoot()
	for _, file := range files {
		vc.fileHistory[file], _ = vc.wd.CatFile(file)
		vc.status.ModifiedFiles = removeFromSlice(vc.status.ModifiedFiles, file)
		vc.status.StagedFiles = removeFromSlice(vc.status.StagedFiles, file)
		vc.status.StagedFiles = append(vc.status.StagedFiles, file)
	}
}

func (vc VC) Commit(message string) {
	totalCommit := len(vc.commits)
	totalCommit++
	files := make(map[string]string)
	for key, value := range vc.fileHistory {
		files[key] = value
	}
	vc.commits[totalCommit] = &CommitContent{
		Message: message,
		Files:   files,
	}
	vc.status.StagedFiles = []string{}
	vc.status.ModifiedFiles = []string{}
}

func (vc VC) Add(files ...string) {
	for _, file := range files {
		vc.fileHistory[file], _ = vc.wd.CatFile(file)
		vc.status.ModifiedFiles = removeFromSlice(vc.status.ModifiedFiles, file)
		vc.status.StagedFiles = removeFromSlice(vc.status.StagedFiles, file)
		vc.status.StagedFiles = append(vc.status.StagedFiles, file)
	}
}

func (vc VC) Checkout(s string) (*workdir.WorkDir, error) {
	commitId := len(vc.commits) - convert(s)
	commit, ok := vc.commits[commitId]
	if ok {
		rootDirectory := vc.wd.RootDirectory + strconv.Itoa(commitId) + "/"
		if _, err := os.Stat(rootDirectory); errors.Is(err, os.ErrNotExist) {
			err := os.Mkdir(rootDirectory, 0666)
			if err != nil {
				panic(err)
			}
			for key, value := range commit.Files {
				if !strings.Contains(key, ".") {
					err := os.MkdirAll(rootDirectory+key, 0660)
					if err != nil {
						return nil, err
					}
				} else {
					err := os.MkdirAll(filepath.Dir(rootDirectory+key), 0660)
					if err != nil {
						return nil, err
					}

					f, err := os.Create(rootDirectory + key)
					if err != nil {
						return nil, err
					}
					err = os.WriteFile(rootDirectory+key, []byte(value), 0666)
					if err != nil {
						return nil, err
					}
					err = f.Close()
					if err != nil {
						return nil, err
					}
				}
			}
		}
		return &workdir.WorkDir{RootDirectory: rootDirectory}, nil
	} else {
		return nil, errors.New("invalid commit")
	}
}

func (vc VC) Log() []string {
	commitsMessage := make([]string, 0)
	for i := len(vc.commits); i >= 1; i-- {
		commitsMessage = append(commitsMessage, vc.commits[i].Message)
	}
	return commitsMessage
}

func removeFromSlice(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

func convert(s string) int {
	if strings.HasPrefix(s, "~") {
		i, _ := strconv.Atoi(s[1:])
		return i
	} else if strings.HasPrefix(s, "^") {
		return len(s)
	} else {
		return 0
	}
}
