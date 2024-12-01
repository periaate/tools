package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		os.Args = append(os.Args, ".")
	}

	dir := os.Args[1]
	goModPath := filepath.Join(dir, "go.mod")

	file, err := os.Open(goModPath)
	if err != nil {
		log.Fatalf("Failed to open go.mod file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "module ") {
			moduleName := strings.TrimSpace(strings.TrimPrefix(line, "module"))
			fmt.Print(moduleName)
			return
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading go.mod file: %v", err)
	}

	log.Fatalf("Failed to find module name in go.mod file")
}
