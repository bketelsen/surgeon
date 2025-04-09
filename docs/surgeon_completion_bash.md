## surgeon completion bash

Generate the autocompletion script for bash

### Synopsis

Generate the autocompletion script for the bash shell.

This script depends on the 'bash-completion' package.
If it is not installed already, you can install it via your OS's package manager.

To load completions in your current shell session:

	source <(surgeon completion bash)

To load completions for every new session, execute once:

#### Linux:

	surgeon completion bash > /etc/bash_completion.d/surgeon

#### macOS:

	surgeon completion bash > $(brew --prefix)/etc/bash_completion.d/surgeon

You will need to start a new shell for this setup to take effect.


```
surgeon completion bash
```

### Options

```
  -h, --help              help for bash
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
  -c, --config-file string    (default "/var/home/bjk/projects/surgeon/.surgeon.yaml")
      --log-level log        logging level [debug|info|warn|error] (default info)
```

### SEE ALSO

* [surgeon completion](surgeon_completion.md)	 - Generate the autocompletion script for the specified shell

