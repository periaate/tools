package main

import (
	"os/exec"

	"github.com/periaate/blume/fsio"
	"github.com/periaate/blume/gen"
	"github.com/periaate/blume/gen/T"
	"github.com/periaate/blume/yap"
)

func main() {
	fsio.Args(T.Len[string](T.NotZero[int])).Match(
		func(s []string) {
			yap.Info("command ["+s[0]+"] found at:", gen.Must(exec.LookPath(s[0])))
		},
		func(e T.Error[any]) { yap.Fatal("no arguments given") },
	)
}
