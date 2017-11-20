# Perigord: Golang Tools for Ethereum Development

*Note:* Perigord is in development and it's API is subject to change.

![Perigord Image (not covered by LICENSE)](https://polyswarm.io/img/perigord-logo-small.jpg)

This image is not covered by LICENSE. 

## Install

There is a Dockerfile in `docker/Dockerfile` to build a perigord image, to build
run

```
$ pushd docker
$ docker build -t perigord .
$ popd
```

These instructions assume an Ubuntu 16.04 x86\_64 environment.

### Prerequisite: Golang 1.8

Some dependencies require Go 1.7+, but Go 1.6 is in Ubuntu 16.04's default repos.
The below will install Go 1.8.


```
$ sudo add-apt-repository -y ppa:longsleep/golang-backports
$ sudo apt-get update
$ sudo apt-get install -y golang-go
$ mkdir $HOME/golang
$ echo "export GOPATH=$HOME/golang" >> ~/.bashrc
$ echo "export PATH=$PATH:$HOME/golang/bin" >> ~/.bashrc
```

Close / re-open your terminal or re-`source` your `.bashrc`.

### Prerequisite: `solc`

```
$ sudo add-apt-repository -y ppa:ethereum/ethereum
$ sudo apt-get update
$ sudo apt-get install -y solc
```

### Prerequisite: `abigen`

```
$ go get github.com/ethereum/go-ethereum
$ pushd $GOPATH/src/github.com/ethereum/go-ethereum
$ go install ./cmd/abigen
$ popd
```

### Build Dependency: `go-bindata`

```
$ go get -u github.com/jteeuwen/go-bindata/...
```

## Setup

```
$ go get -u github.com/polyswarm/perigord/...
$ pushd $GOPATH/src/github.com/polyswarm/perigord
$ go generate
$ cd perigord
$ go install
$ popd
```

## Usage

Run for usage information:

```
$ perigord
```

## Tutorial

[Refer to our introductory blog post for now.](https://medium.com/@swarmmarket/introducing-perigord-golang-tools-for-ethereum-dapp-development-60556c2d9fd)

