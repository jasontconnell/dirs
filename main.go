package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/jasontconnell/dirs/commands"
)

var cmds map[string]commands.Command

func init() {
	cmds = make(map[string]commands.Command)
	installCommand("clearequal", commands.ClearEqual{})
	installCommand("clearside", commands.ClearSide{})
	installCommand("clearempty", commands.ClearEmpty{})
}

func main() {
	c := flag.String("c", "", "command")
	flag.Parse()

	// two dirs
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(fmt.Errorf("getting current working directory. %w", err))
	}
	left := getDir(cwd, flag.Arg(0))
	right := getDir(cwd, flag.Arg(1))

	cmdinst, ok := cmds[*c]

	if !ok {
		log.Fatal(fmt.Errorf("command not defined %s", *c))
	}

	result := cmdinst.Run(left, right)

	if result.Error != nil {
		log.Fatal(fmt.Errorf("executing %s", *c))
	} else {
		log.Printf("%s: success. affected: %d\n", *c, result.Affected)
	}
}

func getDir(cwd, dir string) string {
	d := dir
	if !filepath.IsAbs(d) {
		d = filepath.Join(cwd, d)
	}
	return d
}

func installCommand(cmdnm string, cmd commands.Command) {
	cmds[cmdnm] = cmd
}
