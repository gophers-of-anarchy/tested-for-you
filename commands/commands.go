package commands

import (
	"time"
	"vc/workdir"
)

type FileHistory struct {
	FilesInfo map[string]int64
}

type VC struct {
	wd          *workdir.WorkDir
	status      *Status
	commits     map[int]string
	fileHistory *FileHistory
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
		commits: make(map[int]string),
		fileHistory: &FileHistory{
			FilesInfo: make(map[string]int64),
		},
	}
}

func (vc VC) GetWorkDir() *workdir.WorkDir {
	return vc.wd
}

func (vc VC) Status() *Status {
	if len(vc.fileHistory.FilesInfo) == 0 {
		return vc.status
	}
	for _, file := range vc.wd.ListFilesRoot() {
		lastStatusTime, ok := vc.fileHistory.FilesInfo[file]
		if ok {
			lastModifiedTime := workdir.GetModTimeOfFile(vc.wd.RootDirectory + file)
			if lastStatusTime < lastModifiedTime {
				vc.fileHistory.FilesInfo[file] = time.Now().UnixNano()
				vc.status.ModifiedFiles = append(vc.status.ModifiedFiles, file)
			}
		} else {
			vc.fileHistory.FilesInfo[file] = time.Now().UnixNano()
			vc.status.ModifiedFiles = append(vc.status.ModifiedFiles, file)
		}
	}
	return vc.status
}

func (vc VC) AddAll() {
	files := vc.wd.ListFilesRoot()
	for _, file := range files {
		vc.fileHistory.FilesInfo[file] = time.Now().UnixNano()
		vc.status.ModifiedFiles = removeFromSlice(vc.status.ModifiedFiles, file)
		vc.status.StagedFiles = removeFromSlice(vc.status.StagedFiles, file)
		vc.status.StagedFiles = append(vc.status.StagedFiles, file)
	}
}

func (vc VC) Commit(message string) {
	totalCommit := len(vc.commits)
	totalCommit++
	vc.commits[totalCommit] = message
	vc.status.StagedFiles = []string{}
	vc.status.ModifiedFiles = []string{}
}

func (vc VC) Add(files ...string) {
	for _, file := range files {
		vc.fileHistory.FilesInfo[file] = time.Now().UnixNano()
		vc.status.ModifiedFiles = removeFromSlice(vc.status.ModifiedFiles, file)
		vc.status.StagedFiles = removeFromSlice(vc.status.StagedFiles, file)
		vc.status.StagedFiles = append(vc.status.StagedFiles, file)
	}
}

func (vc VC) Checkout(s string) (*workdir.WorkDir, error) {
	return nil, nil
}

func (vc VC) Log() []string {
	commitsMessage := make([]string, 0)
	for i := len(vc.commits); i >= 1; i-- {
		commitsMessage = append(commitsMessage, vc.commits[i])
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
