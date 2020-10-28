# Barcelona Client

The command line client for [Barcelona](https://github.com/degica/barcelona).

## Installation

Go to [Releases page](https://github.com/degica/barcelona-cli/releases) and download the file for your platform.
Unzip the file and place the binary to your `PATH`

## Usage

`bcn help`

## Development

Requirements:

- [Install Go](https://golang.org/doc/install)

### Getting setup

Simply check out the repository and download the modules required by barcelona-cli. Run `make test` and ensure the tests pass

```bash
git clone https://github.com/degica/barcelona-cli bcn
cd bcn
go mod download
make test
```

### Creating a build

Running `make` will issue a development executable, `barcelona-cli`, in the root of the project.

### Formatting

Run `make format` to format your code!

### Vetting

Run `make vet` to ensure your code meets all the go conventions
