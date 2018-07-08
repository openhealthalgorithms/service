# Open Health Algorithms Service

This service will generate binaries to be used as a package in any posix system. It will expose a `service` in `9595` port.

[CHANGELOG](./CHANGELOG.md)

## Table of Contents

* [Learning](#learning)
* [Developing](#developing)
* [Cross-Platform Code](./CONTRIBUTING.md#cross-platform-code)
* [Committing](./CONTRIBUTING.md#committing)
* [External Dependencies](./CONTRIBUTING.md#external-dependencies)
* [Tools](./CONTRIBUTING.md#tools)
* [Coding Style](./CONTRIBUTING.md#coding-style)

## Learning

Take a look at our [Reading list](./MATERIALS.md).

## Developing

**BEFORE** you start, please read the [Contributing Guideline](./CONTRIBUTING.md).

**Cross-Platform Guidelines** are located [here](./CONTRIBUTING.md#cross-platform-code).

### macOS

To be done.

### Linux

To be done.

#### Clone the repository

The Project follows the common convention in the Golang community about the workspace.  
This means that you should store all Go code under `GOPATH/src` and use full import paths.  

The repo should be cloned like this:

```bash
# open Git Bash command line

# change current directory to user's HOME directory
cd $HOME

# create directory for the project
mkdir -p $GOPATH/src/github.com/openhealthalgorithms

# change current directory
cd $GOPATH/src/github.com/openhealthalgorithms

# clone the repository using https
git clone https://github.com/openhealthalgorithms/service.git

# clone the repository using ssh and ssh keys
git clone git@github.com:openhealthalgorithms/service.git

# change directory to projects's directory
cd service

# and start having fun
```

### Documentation

Run `godoc` and open [documentation](http://localhost:6060/pkg/github.com/trilogy-group/aurea-eng-docker-automation-agent/) in your browser:

```bash
# run godoc
godoc -http=:6060 -tabwidth=4
```

### Golint

Run `golint`:

```bash
# cmd
golint ./cmd/...

# pkg
golint ./pkg/...
```

There is no need to run `golint ./...` since we're not interested in checking dependencies for the code standards.

### Gofmt

Run `gofmt`:

```bash
# cmd
gofmt -w ./cmd

# pkg
gofmt -w ./pkg
```

There is no need to run `gofmt -w ./` since we're not interested in formatting dependencies.

## Build

List all the available targets:

```bash
make
```

Build while developing:

```bash
# for mac
make build_dev_darwin

# for linux
make build_dev_linux

# for windows
make build_dev_windows
```

Build regular binaries:

```bash
# for mac
make build_darwin

# for linux
make build_linux

# for windows
make build_windows
```

Build the artifacts:

```bash
# build in a sequence
make artifacts

# build in parallel
make -j2 artifacts
```

## Contributing

Read about contributing [here](./CONTRIBUTING.md).
