package main

import (
	"fmt"

	"github.com/periaate/blume/fsio"
	"github.com/periaate/blume/gen"
	"github.com/periaate/blume/str"
)

func main() {
	sargs := fsio.SepArgs()
	res := gen.Filter(str.Contains(sargs[0]...))(sargs[1])
	for _, k := range res {
		fmt.Println(k)
	}
}
