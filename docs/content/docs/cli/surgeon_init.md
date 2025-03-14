---
date: 2025-03-10T19:03:11Z
title: "surgeon init"
slug: surgeon_init
url: /docs/cli/surgeon_init/
---
## surgeon init

Initialize a new surgical fork

### Synopsis

The init command will create a new '.surgeon.yaml'
file in the current directory.	This file will contain
the configuration for the surgeon command.

Example configuration file:

	upstream: https://github.com/community-scripts/ProxmoxVE
	modsdir: mymods


```
surgeon init [flags]
```

### Options

```
      --debug   debug output
  -h, --help    help for init
```

### Options inherited from parent commands

```
      --config string     config file (default is .surgeon.yaml)
      --modsdir string    directory containing code modification files
      --upstream string   upstream repository
      --verbose           display verbose output
```

### SEE ALSO

* [surgeon](/surgeon/docs/cli/surgeon/)	 - Surgical forks of upstream repositories

###### Auto generated by toolbox on 10-Mar-2025
