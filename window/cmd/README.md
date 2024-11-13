# window
The window CLI allows you to manage application windows through basic commands to move, resize, open, close, minimize, focus, and list windows.

## Usage
Run window followed by a command and any required arguments.

### Commands
- move: Move a window to a specific position.
- resize: Resize a window to specified dimensions.
- open: Open a minimized or hidden window.
- close: Close a window.
- minimize: Minimize a window.
- find: Find a window by its title.
- focus: Bring a window to the foreground.
- list: List all currently open windows.

### Examples
```bash
Copy code
window move "My App" 100 150
window resize "My App" 800 600
window open "My App"
window close "My App"
window minimize "My App"
window find "My App"
window focus "My App"
window list
```

Each command will output relevant information or an error if the specified window cannot be found or the action cannot be completed.
