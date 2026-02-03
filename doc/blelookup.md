## blelookup



### Options

```
      --debug           enable debug output
      --format string   output format: table|json|csv (default "table")
  -h, --help            help for blelookup
      --verbose         enable verbose output
```

* [blelookup doc]()	 - Generate markdown documentation to stdout
* [blelookup scan]()	 - Scan for BLE devices.

## blelookup completion

Generate the autocompletion script for the specified shell

### Synopsis

Generate the autocompletion script for blelookup for the specified shell.
See each sub-command's help for details on how to use the generated script.


### Options

```
  -h, --help   help for completion
```

### Options inherited from parent commands

```
      --debug           enable debug output
      --format string   output format: table|json|csv (default "table")
      --verbose         enable verbose output
```

* [blelookup completion bash]()	 - Generate the autocompletion script for bash
* [blelookup completion fish]()	 - Generate the autocompletion script for fish
* [blelookup completion powershell]()	 - Generate the autocompletion script for powershell
* [blelookup completion zsh]()	 - Generate the autocompletion script for zsh

## blelookup completion bash

Generate the autocompletion script for bash

### Synopsis

Generate the autocompletion script for the bash shell.

This script depends on the 'bash-completion' package.
If it is not installed already, you can install it via your OS's package manager.

To load completions in your current shell session:

	source <(blelookup completion bash)

To load completions for every new session, execute once:

#### Linux:

	blelookup completion bash > /etc/bash_completion.d/blelookup

#### macOS:

	blelookup completion bash > $(brew --prefix)/etc/bash_completion.d/blelookup

You will need to start a new shell for this setup to take effect.


```
blelookup completion bash
```

### Options

```
  -h, --help              help for bash
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
      --debug           enable debug output
      --format string   output format: table|json|csv (default "table")
      --verbose         enable verbose output
```


## blelookup completion fish

Generate the autocompletion script for fish

### Synopsis

Generate the autocompletion script for the fish shell.

To load completions in your current shell session:

	blelookup completion fish | source

To load completions for every new session, execute once:

	blelookup completion fish > ~/.config/fish/completions/blelookup.fish

You will need to start a new shell for this setup to take effect.


```
blelookup completion fish [flags]
```

### Options

```
  -h, --help              help for fish
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
      --debug           enable debug output
      --format string   output format: table|json|csv (default "table")
      --verbose         enable verbose output
```


## blelookup completion help

Help about any command

### Synopsis

Help provides help for any command in the application.
Simply type completion help [path to command] for full details.

```
blelookup completion help [command] [flags]
```

### Options

```
  -h, --help   help for help
```

### Options inherited from parent commands

```
      --debug           enable debug output
      --format string   output format: table|json|csv (default "table")
      --verbose         enable verbose output
```


## blelookup completion powershell

Generate the autocompletion script for powershell

### Synopsis

Generate the autocompletion script for powershell.

To load completions in your current shell session:

	blelookup completion powershell | Out-String | Invoke-Expression

To load completions for every new session, add the output of the above command
to your powershell profile.


```
blelookup completion powershell [flags]
```

### Options

```
  -h, --help              help for powershell
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
      --debug           enable debug output
      --format string   output format: table|json|csv (default "table")
      --verbose         enable verbose output
```


## blelookup completion zsh

Generate the autocompletion script for zsh

### Synopsis

Generate the autocompletion script for the zsh shell.

If shell completion is not already enabled in your environment you will need
to enable it.  You can execute the following once:

	echo "autoload -U compinit; compinit" >> ~/.zshrc

To load completions in your current shell session:

	source <(blelookup completion zsh)

To load completions for every new session, execute once:

#### Linux:

	blelookup completion zsh > "${fpath[1]}/_blelookup"

#### macOS:

	blelookup completion zsh > $(brew --prefix)/share/zsh/site-functions/_blelookup

You will need to start a new shell for this setup to take effect.


```
blelookup completion zsh [flags]
```

### Options

```
  -h, --help              help for zsh
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
      --debug           enable debug output
      --format string   output format: table|json|csv (default "table")
      --verbose         enable verbose output
```


## blelookup doc

Generate markdown documentation to stdout

```
blelookup doc [flags]
```

### Options

```
  -h, --help   help for doc
```

### Options inherited from parent commands

```
      --debug           enable debug output
      --format string   output format: table|json|csv (default "table")
      --verbose         enable verbose output
```


## blelookup help

Help about any command

### Synopsis

Help provides help for any command in the application.
Simply type blelookup help [path to command] for full details.

```
blelookup help [command] [flags]
```

### Options

```
  -h, --help   help for help
```

### Options inherited from parent commands

```
      --debug           enable debug output
      --format string   output format: table|json|csv (default "table")
      --verbose         enable verbose output
```


## blelookup scan

Scan for BLE devices.

### Synopsis

Scan for BLE (Bluetooth Low Energy) devices.

```
blelookup scan [flags]
```

### Options

```
  -h, --help   help for scan
```

### Options inherited from parent commands

```
      --debug           enable debug output
      --format string   output format: table|json|csv (default "table")
      --verbose         enable verbose output
```


