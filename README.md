# PLEASE NOTE: This is an experimental project still in early development

# Open Health Algorithms Service

[![Go Report Card](https://goreportcard.com/badge/github.com/openhealthalgorithms/service)](https://goreportcard.com/report/github.com/openhealthalgorithms/service)

This service will generate binaries to be used as a package in any posix system. It will expose a `service` in `9595` port.

[CHANGELOG](./CHANGELOG.md)

## Table of Contents

* [Learning](#learning)
* [Developing](#developing)
* [Build](#build)
* [Deploy](#deploy)
* [Committing](./CONTRIBUTING.md#committing)

## Learning

Take a look at our [Reading list](./MATERIALS.md).

## Developing

**BEFORE** you start, please read the [Contributing Guideline](./CONTRIBUTING.md).

### Clone the repository

~~The Project follows the common convention in the Golang community about the workspace.  
This means that you should store all Go code under `GOPATH/src` and use full import paths.~~

Starting with version 1.0.x, we are using `go modules` on this project. That means you will need to clone this
project outside of your `$GOPATH/src` directory.

The repo should be cloned like this:

```bash
# open Git Bash command line

# change current directory to user's HOME directory
cd $HOME

# create directory for the project (whatever you like)
mkdir -p $HOME/goprojects/openhealthalgorithms

# change current directory
cd $HOME/goprojects/openhealthalgorithms

# clone the repository using https
git clone https://github.com/openhealthalgorithms/service.git

# clone the repository using ssh and ssh keys
git clone git@github.com:openhealthalgorithms/service.git

# change directory to projects's directory
cd service
```

### Running in Development

You can run `go run main.go` from the root directory (i.e. `$HOME/goprojects/openhealthalgorithms/service`)

## Build

List all the available targets:

```bash
make
```

**Note: This package needs `gcc` installed on your system to help compile this. So, you need to use separate system to generate the binaries.**

Build regular binaries:

```bash
# on mac
make build_darwin

# on linux
make build_linux
```

> You can use a temporary VM to generate/compile this.

## Deploy

### Run Directly

After you run the build commands, it will generate two binaries in the `artifacts` directory:

- `ohas-darwin-v1.x.x`
- `ohas-linux-v1.x.x`

You can just run the application from command line:

```bash
# run in the foreground
./ohas-linux-v1.0.2

# run in the background (need to kill the process afterwards, if you want to stop)
./ohas-linux-v1.0.2 > /dev/null 2>&1 &

# run in the background when the terminal is closed (need to kill the process afterwards, if you want to stop)
nohup ./ohas-linux-v1.0.2 &
```

### Run as a Service

You can run this as a service on your system too. For this, we will show how to do that in a `Ubuntu 18.04` server.

- Save your configuration on `/etc/ohas/ohas.toml` file.
- Create a file `ohas.service` on `/lib/systemd/system/` directory.

```bash
sudo nano /lib/systemd/system/ohas.service
```

Write the following:

```bash
[Unit]
Description=Open Health Algorithms Service

[Service]
ExecStart=/usr/local/bin/ohas
User=root
Group=root
UMask=007
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

- Upload the binary (`ohas-linux-v1.0.2`)
- Run the followings:

```bash
cp ohas-linux-v1.0.2 ohas               # rename the binary
sudo mv ohas /usr/local/bin/            # move the binary to a directory in the path
sudo chown root: /usr/local/bin/ohas    # change the owner to root
sudo chmod +x /usr/local/bin/ohas       # make the file executable
sudo systemctl restart ohas             # start/restart the service
sudo systemctl status ohas              # check the status
```

> This service will write output on `stdout`. If you setup any logging system (like stackdriver logging),
> all the output can be viewed on the logging application.

## Contributing

Read about contributing [here](./CONTRIBUTING.md).
