# sed CodeMod

Replace strings in a file

## Usage

```
Replace a strings in a file.
This codemod replaces strings in the matched file(s) with a string
specified in the arguments.

Args (2 required):
	1. search string
	2. replacement string

Example:
	upstream: https://github.com/community-scripts/ProxmoxVE
	modsdir: codemods
	codemods:
	- description: Header Updates
		mod: sed
		match: misc/*.func
		args:
		- https://github.com/community-scripts/ProxmoxVE/raw/main/ct/headers/
		- https://github.com/bketelsen/IncusScripts/raw/main/ct/headers/
	
```
