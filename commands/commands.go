package commands

import "vc/workdir"

type VC struct {
	RootDirectory string
}

func Init(wd *workdir.WorkDir) *VC {
	return &VC{
		RootDirectory: wd.RootDirectory,
	}
}
