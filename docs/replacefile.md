# replacefile CodeMod

Replace a file with another

## Usage

```
Replace a file with another.
This codemod replaces the matched file(s) with a file from your fork.

Args (1 required):
	1. The path to the file (in your fork) to replace the matched file(s)

Example:
	upstream: https://github.com/community-scripts/ProxmoxVE
	modsdir: codemods
	codemods:
	- description: Replace create_lxc
		mod: replacefile
		match: ct/create_lxc.sh
		args:
		- codemods/create_lxc.sh
	
```
