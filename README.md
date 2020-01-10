[![CircleCI](https://circleci.com/gh/amitizle/muffin/tree/master.svg?style=svg)](https://circleci.com/gh/amitizle/muffin/tree/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/amitizle/muffin)](https://goreportcard.com/report/github.com/amitizle/muffin)
[![GolangCI](https://golangci.com/badges/github.com/amitizle/muffin.svg)](https://golangci.com/r/github.com/amitizle/muffin)

# Muffin

A small application to make checks (HTTP for now, more in the future).
It's stateless and does not require any external dependency.

All of the confugration is done by using a YAML file and environment variables.


## Build

There is a supplied `Makefile` in the repository, you can use the `build` target (`make build`).
You can optionally choose to change the output, for example:

```bash
$ make build BINARY=~/bin/muffin
```

## Testing

```bash
$ make test
```

## Types Of Checks

TODO

## Configuration

TODO

## Docker Image

TODO
