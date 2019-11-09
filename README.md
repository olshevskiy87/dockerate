# Dockerate

[![Build Status](https://travis-ci.org/olshevskiy87/dockerate.svg?branch=master)](https://travis-ci.org/olshevskiy87/dockerate) [![MIT license](https://img.shields.io/badge/License-MIT-yellow.svg)](https://lbesson.mit-license.org/) ![GitHub release](https://img.shields.io/github/v/tag/olshevskiy87/dockerate?label=release)

Decorate docker commands output.

## Installation

Here are two ways to install Dockerate on your system:

* Download pre-compiled binaries on the [releases page](https://github.com/olshevskiy87/dockerate/releases)

* Compile from source code

  - download and install [Go](https://golang.org/dl/) if you haven't installed it yet
  - get Dockerate sources: ```go get github.com/olshevskiy87/dockerate```
  - move to sources folder: ```cd $GOPATH/src/github.com/olshevskiy87/dockerate```
  - install all dependencies (including golangci-lint): ```make deps```
  - compile binary for your platform and architecture: ```make build``` or just ```make```

## Commands

### dockerate-ps

Displays docker containers info like ```docker ps```.

```
Usage: dockerate-ps [--all] [--no-trunc] [--quiet] [--size] [--color COLOR] [--apiver APIVER] [--verbose]

Options:
  --all, -a              show all containers
  --no-trunc             don't truncate output
  --quiet, -q            only display containers IDs
  --size, -s             display containers sizes
  --color COLOR          when to use colors: always, auto, never [default: auto]
  --apiver APIVER        docker server API version, env DOCKER_API_VERSION
  --verbose, -v          output more information
  --help, -h             display this help and exit
  --version              display version and exit
```

## Motivation

Package [docker-color-output](https://github.com/devemio/docker-color-output) written in PHP by [Sergey Sorokin](https://github.com/devemio).


## License

MIT. See file LICENSE for details.
