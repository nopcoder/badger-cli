# badger-cli


### Install

```sh
go install github.com/nopcoder/badger-cli@latest
```


### Usage

```
Command line client for managing a badger database

Usage:
  badger-cli [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  delete      Delete a key and its contents
  get         Get content of a specific key
  help        Help about any command
  list        List keys in the database
  set         Set a key and its value

Flags:
  -d, --dir string   Path to the badger database directory (default ".")
  -h, --help         help for badger-cli
  -v, --version      version for badger-cli

Use "badger-cli [command] --help" for more information about a command.
```
