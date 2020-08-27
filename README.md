# Dockerate

![Go][1] [![MIT license][2]][3] ![GitHub release][4] [![Go Report Card][5]][6]

Decorate docker commands output.

## Installation

Here are two ways to install Dockerate on your system:

* Download pre-compiled binaries on the [releases page][7]

* Compile from source code

  - download and install [Go][8] (1.13+) if you haven't installed it yet
  - get Dockerate sources: `go get github.com/olshevskiy87/dockerate`
  - move to sources folder: `cd $GOPATH/src/github.com/olshevskiy87/dockerate`
  - install all dependencies (including golangci-lint): `make deps`
  - compile binary for your platform and architecture: `make build` or just `make`

## Commands

### dockerate-ps

Displays docker containers info like `docker ps`.

```
Usage: dockerate-ps [--all] [--no-trunc] [--quiet] [--size] [--color COLOR] [--name-like NAME-LIKE] [--name-ilike NAME-ILIKE] [--columns COLUMNS] [--apiver APIVER] [--verbose]

Options:
  --all, -a              show all containers
  --no-trunc             don't truncate output
  --quiet, -q            only display containers IDs
  --size, -s             display containers sizes
  --color COLOR          when to use colors: always, auto, never [default: auto]
  --name-like NAME-LIKE
                         container name pattern
  --name-ilike NAME-ILIKE
                         container name pattern (case insensitive)
  --columns COLUMNS      columns names to display (case insensitive, separated by commas)
  --apiver APIVER        docker server API version, env DOCKER_API_VERSION
  --verbose, -v          output more information
  --help, -h             display this help and exit
  --version              display version and exit
```

## Motivation

Package [docker-color-output][9] written in PHP by [Sergey Sorokin][10].

## License

MIT. See file LICENSE for details.

[1]: https://github.com/olshevskiy87/dockerate/workflows/Go/badge.svg
[2]: https://img.shields.io/badge/License-MIT-yellow.svg
[3]: https://lbesson.mit-license.org/
[4]: https://img.shields.io/github/v/tag/olshevskiy87/dockerate?label=release
[5]: https://goreportcard.com/badge/github.com/olshevskiy87/dockerate
[6]: https://goreportcard.com/report/github.com/olshevskiy87/dockerate
[7]: https://github.com/olshevskiy87/dockerate/releases
[8]: https://golang.org/dl/
[9]: https://github.com/devemio/docker-color-output
[10]: https://github.com/devemio
