package commands

import (
	"path/filepath"
)

type ClearEmpty struct{}

func (c ClearEmpty) Description() string {
	return "Given two directories, remove directories that are empty"
}

func (c ClearEmpty) Run(left, right string) Result {
	res := Result{}
	var err error


	ldirs := allDirs(left)
	rdirs := allDirs(right)


	dellist := []string{}
	for _, d := range ldirs {
		if d.files == 0 && len(d.subs) == 0 {
			dellist = append(dellist, filepath.Join(left, d.path))
		}
	}

	for _, d := range rdirs {
		if d.files == 0 && len(d.subs) == 0 {
			dellist = append(dellist, filepath.Join(right, d.path))
		}
	}

	err = removeItems(dellist)

	res.Error = err
	res.Success = err == nil
	res.Affected = len(dellist)

	return res
}
