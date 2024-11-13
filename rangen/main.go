package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/periaate/blume/clog"
	"github.com/periaate/blume/fsio"
	"github.com/periaate/blume/gen"
)

func main() {
	args := fsio.Args()
	if len(args) == 0 {
		clog.Fatal("no arguments given")
	}

	n := strings.Split(args[0], "..")
	a := gen.Must(strconv.Atoi(n[0]))
	b := gen.Must(strconv.Atoi(n[1]))

	if a > b {
		for i := a; i > 1; i-- {
			fmt.Println(i)
		}
	} else {
		for i := a; i <= b; i++ {
			fmt.Println(i)
		}
	}
}
