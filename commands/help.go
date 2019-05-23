package commands

import (
	"path/filepath"
	"os"
	"crypto/md5"
	"io/ioutil"
	"fmt"
	"runtime"
	"sync"
	"strings"
)

type hashResult struct {
	path, hash string
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