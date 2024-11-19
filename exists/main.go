package main

import (
	"os/exec"

	"github.com/periaate/blume/clog"
	"github.com/periaate/blume/fsio"
)

func main() {
	args := fsio.Args()
	if len(args) == 0 {
		clog.Fatal("no arguments given")
	}

	cmd := args[0]
	_, err := exec.LookPath(cmd)
	if err != nil {
		clog.Fatal("command not found", "cmd", cmd)
	}

	clog.Info("command found", "cmd", cmd)
}
