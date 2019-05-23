package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/jasontconnell/dirs/commands"

	"github.com/pkg/errors"
)

var cmds map[string]commands.Command

func init() {
	cmds = make(map[string]commands.Command)
	installCommand("clearequal", commands.ClearEqual{})
}

func main() {
	c := flag.String("c", "", "command")
	flag.Parse()

	// two dirs
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(errors.Wrap(err, "getting current working directory"))
	}
	left := getDir(cwd, flag.Arg(0))
	right := getDir(cwd, flag.Arg(1))

	cmdinst, ok := cmds[*c]

	if !ok {
		log.Fatal(errors.Wrapf(errors.New("command not defined"), "command: %s", *c))
	}

	result := cmdinst.Run(left, right)

	if result.Error != nil {
		log.Fatal(errors.Wrapf(result.Error, "executing %s", *c))
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
