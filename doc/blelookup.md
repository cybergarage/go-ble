## blelookup

Scan and inspect the Bluetooth Low Energy devices nearby

### Synopsis

blelookup scans for the Bluetooth Low Energy devices which are advertising
nearby, connects to one of them, and looks up the Bluetooth SIG assigned
numbers.

  blelookup scan                             list the advertising devices
  blelookup scan --service 0xFFF6            list only the devices of a service
  blelookup connect <address> <service>      list the characteristics of a service
  blelookup lookup <uuid|company id>         look up an assigned number

This tool is a central. Advertising as a peripheral is not supported yet.

### Options

```
      --debug              enable debug output
      --format string      output format: table|json|csv (default "table")
  -h, --help               help for blelookup
  -t, --timeout duration   operation timeout (default 5s)
      --verbose            enable verbose output
```

* [blelookup connect]()	 - Connect to a device and list the characteristics of a service
* [blelookup doc]()	 - Generate markdown documentation to stdout
* [blelookup lookup]()	 - Look up a Bluetooth SIG assigned number
* [blelookup scan]()	 - Scan for the advertising BLE devices

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
      --debug              enable debug output
      --format string      output format: table|json|csv (default "table")
  -t, --timeout duration   operation timeout (default 5s)
      --verbose            enable verbose output
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
      --debug              enable debug output
      --format string      output format: table|json|csv (default "table")
  -t, --timeout duration   operation timeout (default 5s)
      --verbose            enable verbose output
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
      --debug              enable debug output
      --format string      output format: table|json|csv (default "table")
  -t, --timeout duration   operation timeout (default 5s)
      --verbose            enable verbose output
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
      --debug              enable debug output
      --format string      output format: table|json|csv (default "table")
  -t, --timeout duration   operation timeout (default 5s)
      --verbose            enable verbose output
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
      --debug              enable debug output
      --format string      output format: table|json|csv (default "table")
  -t, --timeout duration   operation timeout (default 5s)
      --verbose            enable verbose output
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
      --debug              enable debug output
      --format string      output format: table|json|csv (default "table")
  -t, --timeout duration   operation timeout (default 5s)
      --verbose            enable verbose output
```


## blelookup connect

Connect to a device and list the characteristics of a service

### Synopsis

Scan for the device of the specified address, connect to it, and print the
characteristics of the specified service with the negotiated ATT MTU.

The device must be advertising, because it is found by a scan first.

```
blelookup connect <address> <service> [flags]
```

### Examples

```
  blelookup connect 11:22:33:AA:BB:CC 0xFFF6
  blelookup connect --format json 0102030A-0B0C-0D0E-0F10-111213141516 0xFFF6
```

### Options

```
  -h, --help   help for connect
```

### Options inherited from parent commands

```
      --debug              enable debug output
      --format string      output format: table|json|csv (default "table")
  -t, --timeout duration   operation timeout (default 5s)
      --verbose            enable verbose output
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
      --debug              enable debug output
      --format string      output format: table|json|csv (default "table")
  -t, --timeout duration   operation timeout (default 5s)
      --verbose            enable verbose output
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
      --debug              enable debug output
      --format string      output format: table|json|csv (default "table")
  -t, --timeout duration   operation timeout (default 5s)
      --verbose            enable verbose output
```


## blelookup lookup

Look up a Bluetooth SIG assigned number

### Synopsis

Look up a Bluetooth SIG assigned number in the bundled database, and print the
service, the characteristic or the company which it names.

The command needs no Bluetooth hardware, so it can be used to read a UUID which
was captured elsewhere.

```
blelookup lookup <uuid|company id> [flags]
```

### Examples

```
  blelookup lookup 0xFFF6
  blelookup lookup 0000fff6-0000-1000-8000-00805f9b34fb
  blelookup lookup --format json 0x004C
```

### Options

```
  -h, --help   help for lookup
```

### Options inherited from parent commands

```
      --debug              enable debug output
      --format string      output format: table|json|csv (default "table")
  -t, --timeout duration   operation timeout (default 5s)
      --verbose            enable verbose output
```


## blelookup scan

Scan for the advertising BLE devices

### Synopsis

Scan for the Bluetooth Low Energy devices which are advertising nearby, and
print them as they are discovered.

The scan stops after the timeout, or when it is interrupted.

```
blelookup scan [flags]
```

### Examples

```
  blelookup scan
  blelookup scan --service 0xFFF6 --timeout 30s
  blelookup scan --rssi -70 --format json
```

### Options

```
      --address strings   scan only the devices of the address
  -h, --help              help for scan
      --name strings      scan only the devices of the local name
      --rssi int          scan only the devices whose RSSI is equal to or greater than this value
  -s, --service strings   scan only the devices which advertise the service UUID
```

### Options inherited from parent commands

```
      --debug              enable debug output
      --format string      output format: table|json|csv (default "table")
  -t, --timeout duration   operation timeout (default 5s)
      --verbose            enable verbose output
```


