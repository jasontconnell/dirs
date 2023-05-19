package commands

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

type hashResult struct {
	path, hash string
}

type dirResult struct {
	subs   []*dirResult
	parent *dirResult
	path   string
	files  int
}

func allDirs(dir string) []*dirResult {
	pm := make(map[string]*dirResult)
	list := []*dirResult{}
	filepath.Walk(dir, func(p string, f os.FileInfo, err error) error {
		parent, ok := pm[filepath.Dir(p)]
		if !f.IsDir() {
			if ok {
				parent.files++
			}
			return nil
		}

		reldir := strings.TrimPrefix(p, dir)

		res := &dirResult{path: reldir, files: 0, subs: []*dirResult{}, parent: parent}
		pm[p] = res
		if ok && parent != nil {
			parent.subs = append(parent.subs, res)
		}

		list = append(list, res)

		return nil
	})
	return list
}

func allPaths(dir string) []string {
	list := []string{}

	filepath.Walk(dir, func(p string, f os.FileInfo, err error) error {
		if f.IsDir() {
			return nil
		}

		list = append(list, p)
		return nil
	})

	return list
}

func allHash(root string, paths []string) []hashResult {
	files := make(chan string)
	res := make(chan hashResult)
	done := make(chan bool)
	hres := []hashResult{}

	go func() {
		for _, p := range paths {
			files <- p
		}
		close(files)
	}()

	go func() {
		for p := range res {
			hres = append(hres, p)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(runtime.NumCPU())
	for i := 0; i < runtime.NumCPU(); i++ {
		go func(fc chan string, hc chan hashResult, dc chan bool) {
			for fn := range fc {
				h, err := getHash(fn)
				if err == nil {
					hc <- hashResult{path: strings.TrimPrefix(fn, root), hash: h}
				}
			}

			wg.Done()
		}(files, res, done)
	}
	wg.Wait()

	return hres
}

func getHash(fn string) (string, error) {
	f, err := os.Open(fn)
	if err != nil {
		return "", err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}

	sum := md5.Sum(b)
	return fmt.Sprintf("%x", sum), nil
}

func removeItems(relpaths []string, basedirs ...string) error {
	var err error
	for _, p := range relpaths {
		fullPaths := []string{}
		if len(basedirs) > 0 {
			for _, d := range basedirs {
				fullPath := filepath.Join(d, p)
				fullPaths = append(fullPaths, fullPath)
			}
		} else {
			fullPaths = append(fullPaths, p)
		}

		for _, d := range fullPaths {
			err := os.Remove(d)
			if err != nil {
				err = fmt.Errorf("deleting %s. %w", d, err)
				break
			}
		}
	}
	return err
}
