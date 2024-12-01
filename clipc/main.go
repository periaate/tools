package main

import (
	"fmt"
	"os"

	"github.com/atotto/clipboard"
	"github.com/periaate/blume/fsio"
)

func main() {
	if clipboard.Unsupported {
		fmt.Println("Clipboard access is not supported on this platform.")
		return
	}

	res := fsio.ReadRawPipe()
	clipboard.WriteAll(string(res))
	if fsio.HasOutPipe() {
		os.Stdout.Write(res)
	}
}
