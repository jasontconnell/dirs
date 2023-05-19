package commands

import (
	"strings"
)

type ClearSide struct{}

func (c ClearSide) Description() string {
	return "Given two directories, remove from the last directory (right) where that file only exists in that directory."
}

func (c ClearSide) Run(left, right string) Result {
	res := Result{}
	llist := allPaths(left)
	rlist := allPaths(right)

	lmap := make(map[string]bool)
	rmap := make(map[string]bool)

	for _, f := range llist {
		k := strings.TrimPrefix(f, left)
		lmap[k] = true
	}

	for _, f := range rlist {
		k := strings.TrimPrefix(f, right)
		rmap[k] = true
	}

	dellist := []string{}
	for k := range rmap {
		if _, ok := lmap[k]; !ok {
			dellist = append(dellist, k)
		}
	}

	err := removeItems(dellist, right)

	res.Error = err
	res.Success = err == nil
	res.Affected = len(dellist)

	return res
}
