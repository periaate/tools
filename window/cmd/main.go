package main

import (
	"strings"

	"github.com/periaate/tools/window"

	"github.com/periaate/blume/clog"
	"github.com/periaate/blume/fsio"
	"github.com/periaate/blume/gen"
)

func main() {
	args := fsio.QArgs()
	if len(args) == 0 {
		clog.Fatal("no arguments provided")
	}

	cmd := args[0]
	args = gen.Shift(args, 1)

	switch cmd {
	case "move":
		moveWindow(args[0].String(), args[1].Int(), args[2].Int())
	case "size", "resize":
		resizeWindow(args[0].String(), args[1].Int(), args[2].Int())
	case "open":
		openWindow(args[0].String())
	case "close":
		closeWindow(args[0].String())
	case "minimize", "min":
		minimizeWindow(args[0].String())
	case "find":
		s := args[0].String()
		if r := find(s); r != nil {
			clog.Info("found window", "title", r.Title())
		} else {
			clog.Fatal("couldn't find window", "title", s)
		}
	case "focus":
		focusWindow(args[0].String())
	case "list":
		fallthrough
	default:
		listWindows()
	}
}

func listWindows() {
	wnds, err := window.ListAllWindows()
	if err != nil || len(wnds) == 0 {
		clog.Fatal("couldn't list windows", "err", err)
	}

	for _, wnd := range wnds {
		clog.Info("found window", "title", wnd.Title())
	}
}

func find(title string) *window.Window {
	windows, err := window.ListAllWindows()
	if err != nil {
		clog.Fatal("couldn't list windows", "err", err)
	}

	for _, wnd := range windows {
		if strings.Contains(strings.ToLower(wnd.Title()), strings.ToLower(title)) {
			return wnd
		}
	}

	return nil
}

func moveWindow(title string, x, y int) {
	wnd := find(title)
	if wnd == nil {
		clog.Fatal("couldn't find window", "title", title)
	}

	if err := wnd.Move(x, y); err != nil {
		clog.Fatal("couldn't move window", "title", wnd.Title(), "err", err)
	}
}

func resizeWindow(title string, x, y int) {
	wnd := find(title)
	if wnd == nil {
		clog.Fatal("couldn't find window", "title", title)
	}

	clog.Info("resizing window", "title", wnd.Title(), "x", x, "y", y)

	if err := wnd.Resize(x+16, y+39); err != nil {
		clog.Fatal("couldn't move window", "title", wnd.Title(), "err", err)
	}
}

func openWindow(title string) {
	wnd := find(title)
	if wnd == nil {
		clog.Fatal("couldn't find window", "title", title)
	}

	if err := wnd.Open(); err != nil {
		clog.Fatal("couldn't show window", "title", wnd.Title(), "err", err)
	}
}

func closeWindow(title string) {
	wnd := find(title)
	if wnd == nil {
		clog.Fatal("couldn't find window", "title", title)
	}

	if err := wnd.Close(); err != nil {
		clog.Fatal("couldn't close window", "title", wnd.Title(), "err", err)
	}
}

func minimizeWindow(title string) {
	wnd := find(title)
	if wnd == nil {
		clog.Fatal("couldn't find window", "title", title)
	}

	if err := wnd.Minimize(); err != nil {
		clog.Fatal("couldn't minimize window", "title", wnd.Title(), "err", err)
	}
}

func focusWindow(title string) {
	wnd := find(title)
	if wnd == nil {
		clog.Fatal("couldn't find window", "title", title)
	}

	if err := wnd.Focus(); err != nil {
		clog.Fatal("couldn't focus window", "title", wnd.Title(), "err", err)
	}
}
