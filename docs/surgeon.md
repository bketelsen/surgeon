# surgeon

Surgically modify your forks

## Synopsis

Surgeon is a tool to make surgical changes to forks of upstream repositories.

The surgeon command reads a configuration file in the current directory
named '.surgeon.yaml'.  This file contains the configuration for the
surgeon command.  The configuration file contains the upstream repository
URL, the directory containing the code modification files, and a list of
code modifications to apply to the forked repository.

The surgeon command will clone the upstream repository into a temporary directory,
then apply the code modifications to the cloned repository.  The contents of the
modified repository are copied to the current directory, overwriting any existing
files.

Important: modifications are applied in the order they are listed in the configuration,
and have a cumulative effect.  Be sure to verify your modifications before committing.

```
surgeon [flags]
```

## Options

```
  -c, --config-file string    (default "/var/home/bjk/projects/surgeon/.surgeon.yaml")
  -h, --help                 help for surgeon
      --log-level log        logging level [debug|info|warn|error] (default info)
```

## See also

* [surgeon codemod](surgeon_codemod.md)	 - Work with codemods
* [surgeon completion](surgeon_completion.md)	 - Generate the autocompletion script for the specified shell
* [surgeon init](surgeon_init.md)	 - Initialize a new surgical fork

