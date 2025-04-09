# bashfunc CodeMod

Replace a bash function with another

## Usage

```
Replace a bash function with another.
This codemod searches for a bash function in the matched file(s)
and replaces it with another function.

Args (2 required):
	1. The name of the function to replace
	2. The path to the file (in your fork) containing the replacement function

Example:
	upstream: https://github.com/community-scripts/ProxmoxVE
	modsdir: codemods
	codemods:
	- description: PVE Check Function
		mod: bashfunc
		match: misc/build.func
		args:
		- pve_check
		- codemods/pve_check.sh
	
```
