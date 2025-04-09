# inject CodeMod

Inject contents into a file

## Usage

```
Inject contents into a file.
This codemod modifies the matched file(s) by injecting specified content.

Args (2 required):
	1. Injection point in the file. Valid: "start", "end", <line number>
	2. The content to inject

Example:
	upstream: https://github.com/community-scripts/ProxmoxVE
	modsdir: codemods
	codemods:
	- description: Inject Modification notice
		mod: inject
		match: install/*.sh
		args:
		- end
		- # Modified by surgeon
	
```
