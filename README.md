# Barcelona Client

The command line client for [Barcelona](https://github.com/degica/barcelona).

## Installation

Go to [Releases page](https://github.com/degica/barcelona-cli/releases) and download the file for your platform.
Unzip the file and place the binary to your `PATH`

## Usage

`bcn help`

## Development

Requirements:

- [Go installed](https://golang.org/doc/install) and GOPATH setup.  
- [Install glide](https://github.com/Masterminds/glide#install)  

### Getting setup

- Example GOPATH: `/home/my_home/go`
- Clone this project into: `/home/my_home/go/src/github.com/degica
- `cd` into `barcelona-cli`
- Run `glide install`

### Creating a build

- Running `make dev` will issue a development executable, `barcelona-cli`, in the root of the project.
