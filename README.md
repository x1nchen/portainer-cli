![Tests](https://github.com/x1nchen/portainer-cli/workflows/Tests/badge.svg?branch=master)
# Portainer CLI


## Availability

## We need your feedback

## Usage

```
Work seamlessly with Portainer from the command line.

Usage:
  portainer-cli [command]

Available Commands:
  config      config to get/set configuration
  help        Help about any command
  login       login to get the auth token
  search      search container with fuzzy name
  sync        sync data from portainer instance to local
  version     Print the version of portainer

Flags:
  -h, --help          help for portainer-cli
      --host string   host base url such as http://localhost:9000

Use "portainer-cli [command] --help" for more information about a command.
```

1. first use login command to get the auth token

```bash
portainer-cli --host https://example.com login -u username -p password
```

2. sync data from remote server to local cache

```bash
portainer-cli --host https://example.com sync
```

3. search container with name fuzzy match

```bash
portainer-cli --host https://example.com search example
```

## Documentation

## Installation

### macOS

`portainer-cli` is available via Homebrew

#### go

```bash
GO111MODULE="on" go get -v github.com/x1nchen/portainer-cli
```

#### Homebrew

Install:

```bash
brew tap x1nchen/tap
brew install portainer-cli
```


Upgrade: 

```bash
brew upgrade portainer-cli
```


