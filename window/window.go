// window/window.go
package window

import (
	"syscall"
	"unsafe"

	"github.com/periaate/blume/clog"
	"golang.org/x/sys/windows"
)

// Window defines a Windows window and various methods to manipulate them.
type Window struct {
	hwnd  windows.HWND
	title string
}

func (w *Window) Title() string { return w.title }

// ListAllWindows enumerates all top-level windows and returns a slice of Window pointers.
func ListAllWindows() (res []*Window, err error) {
	var windowsList []*Window

	// Callback function for EnumWindows
	cb := syscall.NewCallback(func(hwnd syscall.Handle, lParam uintptr) uintptr {
		// Check if window is visible
		isVisible, err := isWindowVisible(windows.HWND(hwnd))
		if err != nil || !isVisible {
			return 1 // Continue enumeration
		}

		// Get window title
		title, err := getWindowText(windows.HWND(hwnd))
		if err != nil || title == "" {
			return 1 // Continue enumeration
		}

		windowsList = append(windowsList, &Window{
			hwnd:  windows.HWND(hwnd),
			title: title,
		})

		return 1 // Continue enumeration
	})

	procEnumWindows.Call(cb, 0)
	if err := windows.GetLastError(); err != nil && err != windows.ERROR_SUCCESS {
		return nil, err
	}

	return windowsList, nil
}

// Resize changes the size of the window to the specified width (x) and height (y).
func (w *Window) Resize(width int, height int) error {
	var rect windows.Rect
	err := getWindowRect(w.hwnd, &rect)
	if err != nil {
		return err
	}

	return moveWindow(w.hwnd, int(rect.Left), int(rect.Top), width, height, true)
}

// Move changes the position of the window to the specified x and y coordinates.
func (w *Window) Move(x int, y int) error {
	var rect windows.Rect
	err := getWindowRect(w.hwnd, &rect)
	if err != nil {
		return err
	}

	width := int(rect.Right - rect.Left)
	height := int(rect.Bottom - rect.Top)

	return moveWindow(w.hwnd, x, y, width, height, true)
}

// Open makes the window visible if it is hidden or minimized.
func (w *Window) Open() error {
	return showWindow(w.hwnd, SW_SHOW)
}

// Minimize minimizes the window.
func (w *Window) Minimize() error {
	return showWindow(w.hwnd, SW_MINIMIZE)
}

// Focus brings the window to the foreground.
func (w *Window) Focus() error {
	return setForegroundWindow(w.hwnd)
}

// Close sends a WM_CLOSE message to the window to close it gracefully.
func (w *Window) Close() error {
	return sendMessage(w.hwnd, WM_CLOSE, 0, 0)
}

// Internal helper functions and constants

var (
	user32                  = windows.NewLazySystemDLL("user32.dll")
	procEnumWindows         = user32.NewProc("EnumWindows")
	procIsWindowVisible     = user32.NewProc("IsWindowVisible")
	procGetWindowTextLength = user32.NewProc("GetWindowTextLengthW")
	procGetWindowText       = user32.NewProc("GetWindowTextW")
	procMoveWindow          = user32.NewProc("MoveWindow")
	procShowWindow          = user32.NewProc("ShowWindow")
	procSetForegroundWindow = user32.NewProc("SetForegroundWindow")
	procSendMessage         = user32.NewProc("SendMessageW")
	procGetWindowRect       = user32.NewProc("GetWindowRect")
)

const (
	SW_HIDE          = 0
	SW_SHOW          = 5
	SW_MINIMIZE      = 6
	SW_SHOWMINIMIZED = 2
	SW_SHOWMAXIMIZED = 3
	WM_CLOSE         = 0x0010
	WM_DESTROY       = 0x0002
	WM_QUIT          = 0x0012
	GWL_STYLE        = -16
	WS_VISIBLE       = 0x10000000
	WS_MINIMIZE      = 0x20000000
	WS_MAXIMIZEBOX   = 0x00010000
	WS_MINIMIZEBOX   = 0x00020000
	WS_SYSMENU       = 0x00080000
	WS_CAPTION       = 0x00C00000
)

// isWindowVisible checks if the window is visible.
func isWindowVisible(hwnd windows.HWND) (bool, error) {
	ret, _, err := procIsWindowVisible.Call(uintptr(hwnd))
	if ret == 0 {
		if err != nil && err.Error() != "The operation completed successfully." {
			return false, err
		}
		return false, nil
	}
	return true, nil
}

// getWindowText retrieves the window's title text.
func getWindowText(hwnd windows.HWND) (string, error) {
	length, _, err := procGetWindowTextLength.Call(uintptr(hwnd))
	if length == 0 {
		if err != nil && err.Error() != "The operation completed successfully." {
			return "", err
		}
		return "", nil
	}

	buffer := make([]uint16, length+1)
	ret, _, err := procGetWindowText.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&buffer[0])), uintptr(length+1))
	if ret == 0 {
		if err != nil && err.Error() != "The operation completed successfully." {
			return "", err
		}
		return "", nil
	}

	return syscall.UTF16ToString(buffer), nil
}

// getWindowRect retrieves the window's bounding rectangle.
func getWindowRect(hwnd windows.HWND, rect *windows.Rect) error {
	ret, _, err := procGetWindowRect.Call(uintptr(hwnd), uintptr(unsafe.Pointer(rect)))
	if ret == 0 {
		return err
	}
	return nil
}

// moveWindow moves and/or resizes the window.
func moveWindow(hwnd windows.HWND, x, y, width, height int, repaint bool) error {
	repaintFlag := 0
	if repaint {
		repaintFlag = 1
	}
	ret, _, err := procMoveWindow.Call(
		uintptr(hwnd),
		uintptr(x),
		uintptr(y),
		uintptr(width),
		uintptr(height),
		uintptr(repaintFlag),
	)
	if ret == 0 {
		return err
	}
	return nil
}

// showWindow changes the window's show state.
func showWindow(hwnd windows.HWND, cmdShow int) error {
	ret, _, err := procShowWindow.Call(uintptr(hwnd), uintptr(cmdShow))
	if err != nil {
		clog.Fatal("couldn't show window", "err", err)
	}
	if ret == 0 {
		// According to documentation, return value is non-zero if window was previously visible
		// It's not necessarily an error if the window was already in the desired state
		// So we don't treat ret == 0 as an error here
	}
	return nil
}

// setForegroundWindow brings the window to the foreground.
func setForegroundWindow(hwnd windows.HWND) error {
	ret, _, err := procSetForegroundWindow.Call(uintptr(hwnd))
	if ret == 0 {
		return err
	}
	return nil
}

// sendMessage sends a message to the window.
func sendMessage(hwnd windows.HWND, msg uint32, wParam, lParam uintptr) error {
	ret, _, err := procSendMessage.Call(
		uintptr(hwnd),
		uintptr(msg),
		wParam,
		lParam,
	)
	if ret == 0 {
		return err
	}
	return nil
}

// Helper function to convert syscall error to Go error
func toError(e error) error {
	if e == syscall.Errno(0) {
		return nil
	}
	return e
}

// Additional function: Enumerate Windows using EnumWindows
func enumerateWindows(enumFunc func(windows.HWND) bool) error {
	enumCallback := syscall.NewCallback(func(hwnd syscall.Handle, lParam uintptr) uintptr {
		if enumFunc(windows.HWND(hwnd)) {
			return 1 // Continue enumeration
		}
		return 0 // Stop enumeration
	})

	ret, _, err := procEnumWindows.Call(enumCallback, 0)
	if ret == 0 {
		if err != nil && err.Error() != "The operation completed successfully." {
			return err
		}
	}
	return nil
}
