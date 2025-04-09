# sjson CodeMod

Modify a JSON file in-place

## Usage

```
sjson modifies a JSON file in-place.
This codemod modifies the matched file(s) by injecting specified content.

Args (3 required for set, 2 required for del):
	1. Action (set, del)
	2. Key path
	3. Value (required for set)

Example:
	upstream: https://github.com/community-scripts/ProxmoxVE
	modsdir: codemods
	codemods:
	- description: change OS Key to debian
		mod: sjson
		match: json/debian-vm.json
		args:
		- set
		- install_methods.1.resources.os
		- debian
	
```
