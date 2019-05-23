package commands

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

type ClearEqual struct{}

func (c ClearEqual) Run(left, right string) Result {
	res := Result{}
	llist := allPaths(left)
	rlist := allPaths(right)

	lhashes := allHash(left, llist)
	rhashes := allHash(right, rlist)

	lmap := make(map[string]string)
	rmap := make(map[string]string)

	for _, h := range lhashes {
		lmap[h.path] = h.hash
	}

	for _, h := range rhashes {
		rmap[h.path] = h.hash
	}

	dellist := []string{}

	for k, v := range lmap {
		if rh, ok := rmap[k]; ok && rh == v {
			dellist = append(dellist, k)
		}
	}

	err := removeFiles(dellist, left, right)

	res.Error = err
	res.Success = err == nil
	res.Affected = len(dellist)

	return res
}

func removeFiles(relpaths []string, left, right string) error {
	for _, p := range relpaths {
		for _, d := range []string{left, right} {
			fullPath := filepath.Join(d, p)
			err := os.Remove(fullPath)
			if err != nil {
				return errors.Wrapf(err, "deleting file %s", fullPath)
			}
		}
	}
	return nil
}
