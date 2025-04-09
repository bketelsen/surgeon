# surgeon completion fish

Generate the autocompletion script for fish

## Synopsis

Generate the autocompletion script for the fish shell.

To load completions in your current shell session:

	surgeon completion fish | source

To load completions for every new session, execute once:

	surgeon completion fish > ~/.config/fish/completions/surgeon.fish

You will need to start a new shell for this setup to take effect.


```
surgeon completion fish [flags]
```

## Options

```
  -h, --help              help for fish
      --no-descriptions   disable completion descriptions
```

## Options inherited from parent commands

```
  -c, --config-file string    (default "/var/home/bjk/projects/surgeon/.surgeon.yaml")
      --log-level log        logging level [debug|info|warn|error] (default info)
```

## See also

* [surgeon completion](surgeon_completion.md)	 - Generate the autocompletion script for the specified shell

