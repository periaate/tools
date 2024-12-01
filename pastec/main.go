package main

import (
	"fmt"
	"log"
	"os"

	"github.com/atotto/clipboard"
)

func main() {
	if clipboard.Unsupported {
		fmt.Println("Clipboard access is not supported on this platform.")
		return
	}
	content, err := clipboard.ReadAll()
	if err != nil {
		log.Fatalf("Failed to read clipboard contents: %v", err)
	}
	os.Stdout.WriteString(content)
}
