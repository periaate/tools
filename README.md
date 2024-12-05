# tools

includes various tools, individually not significant enough to be repositories.

network utilities:
- `auth`: reverse proxy implementation of `blume/hnet/auth`
- `fwauth`: forward auth implementation of `blume/hnet/auth`

dev utilities:
- `tagver`: git tag version utility. `tagver` for current tag version, `tagver patch` for patch+1, ...
- `licenser`: writes a license of choice to current working directory.
- `getmod`: gets the module name of a Go project.

system utilities:
- `pastec`: pastes the contents of the clipboard to the terminal.
- `clipc`: copies the piped content to the clipboard. (can be used in a pass-through fashion)
- `runn`: runs any number of commands, writing all of their stdout, stderr to the terminals.
- `rangen`: generates ranges given two numbers.
- `timer`: simple timer, also plays an audio file.
- `exists`: looks for the input argument in the path.
- `window`: Windows cli utility for manipulating windows.

Barely functional tier:
- `pingc`: http call tool.
- `str`: string tool.
- `tdn`: time date tool.

## License
All code available under GPL-3.0, unless otherwise specified.
