package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"

	"github.com/periaate/blume/clog"
	"github.com/periaate/blume/fsio"
	"github.com/periaate/blume/gen"
	"github.com/periaate/blume/options"
)

type Command struct {
	cmd   *exec.Cmd
	label string
}

func main() {
	fmt.Println("Running command(s)...")
	args, ans := fsio.Args(options.LongerThan[[]string](0))
	if ans != nil {
		clog.Fatal(ans.Name, "reason", ans.Reason)
	}

	cargs := gen.Split(func(s string) bool { return s == "??" })(args)
	cmds := []Command{}
	for _, arguments := range cargs {
		label := arguments[0]
		clog.Info("arguments", "args", arguments)
		cmd := exec.Command(arguments[1], arguments[2:]...)

		cmds = append(cmds, Command{
			label: label,
			cmd:   cmd,
		})
	}

	err := RunCommands(cmds...)
	if err != nil {
		clog.Fatal("error running commands", "error", err)
	}
}

// RunCommands executes a variadic number of exec.Cmd and streams all their outputs to stdout in real-time.
func RunCommands(cmds ...Command) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(cmds))

	for _, cmd := range cmds {
		wg.Add(1)
		go func(cm Command) {
			c := cm.cmd
			defer wg.Done()

			stdoutPipe, err := c.StdoutPipe()
			if err != nil {
				errChan <- fmt.Errorf("error creating stdout pipe: %w", err)
				return
			}

			stderrPipe, err := c.StderrPipe()
			if err != nil {
				errChan <- fmt.Errorf("error creating stderr pipe: %w", err)
				return
			}

			// Start the command
			if err := c.Start(); err != nil {
				errChan <- fmt.Errorf("error starting command: %w", err)
				return
			}

			// Stream the output
			go streamOutput(stdoutPipe, cm.label)
			go streamOutput(stderrPipe, cm.label+" stderr")

			// Wait for the command to finish
			if err := c.Wait(); err != nil {
				errChan <- fmt.Errorf("command finished with error: %w", err)
			}
		}(cmd)
	}

	// Wait for all goroutines to finish
	wg.Wait()
	close(errChan)

	// Collect errors, if any
	var errs []error
	for err := range errChan {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return fmt.Errorf("encountered errors: %v", errs)
	}

	return nil
}

// streamOutput streams the given reader line by line to stdout.
func streamOutput(reader io.Reader, label string) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		fmt.Printf("[%s] %s\n", label, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", label, err)
	}
}
