# Dockerate

[![Build Status][1]][2] [![MIT license][3]][4] ![GitHub release][5] [![Go Report Card][6]][7]

Decorate docker commands output.

## Installation

Here are two ways to install Dockerate on your system:

* Download pre-compiled binaries on the [releases page][8]

* Compile from source code

  - download and install [Go][9] (1.14+) if you haven't installed it yet
  - get Dockerate sources: `go get github.com/olshevskiy87/dockerate`
  - move to sources folder: `cd $GOPATH/src/github.com/olshevskiy87/dockerate`
  - install all dependencies (including golangci-lint): `make deps`
  - compile binary for your platform and architecture: `make build` or just `make`

## Commands

### dockerate-ps

Displays docker containers info like `docker ps`.

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

Package [docker-color-output][10] written in PHP by [Sergey Sorokin][11].

## License

MIT. See file LICENSE for details.

[1]: https://travis-ci.org/olshevskiy87/dockerate.svg?branch=master
[2]: https://travis-ci.org/olshevskiy87/dockerate
[3]: https://img.shields.io/badge/License-MIT-yellow.svg
[4]: https://lbesson.mit-license.org/
[5]: https://img.shields.io/github/v/tag/olshevskiy87/dockerate?label=release
[6]: https://goreportcard.com/badge/github.com/olshevskiy87/dockerate
[7]: https://goreportcard.com/report/github.com/olshevskiy87/dockerate
[8]: https://github.com/olshevskiy87/dockerate/releases
[9]: https://golang.org/dl/
[10]: https://github.com/devemio/docker-color-output
[11]: https://github.com/devemio
