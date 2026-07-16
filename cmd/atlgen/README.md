# ATLGen - Command line utility for sending and receiving data over sound on Linux
`atlgen` allows you to easily send and recieve data over any pulseaudio device.

```
$ atlgen

---

Utilities for generating and parsing ATL (Audio Transport Language).

Usage:
  atlgen [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  list        List all audio devices that can be used for input and output.
  recv        Receive data from an audio input device.
  send        Send data over an audio output device.

Flags:
  -h, --help   help for atlgen

Use "atlgen [command] --help" for more information about a command.
```

## Installation
Either:
    - **Preffered**: Download a `.deb` from the releases page and `$ sudo apt install <name>.deb`
    - Or: Download this repo, `cd` into `cmd/atlgen`, then `$ make install`
    - Or: `$ go install github.com/JoshPattman/atl/cmd/atlgen@latest`

## Usage
The command is well documented - run `$ atlgen help` for usage.