# surgeon init

Initialize a new surgical fork

## Synopsis

The init command will create a new '.surgeon.yaml'
file in the current directory.	This file will contain
the configuration for the surgeon command.

Example configuration file:

	upstream: https://github.com/community-scripts/ProxmoxVE
	modsdir: mymods


```
surgeon init [flags]
```

## Options

```
  -h, --help   help for init
```

## Options inherited from parent commands

```
  -c, --config-file string    (default "/var/home/bjk/projects/surgeon/.surgeon.yaml")
      --log-level log        logging level [debug|info|warn|error] (default info)
```

## See also

* [surgeon](surgeon.md)	 - Surgically modify your forks

