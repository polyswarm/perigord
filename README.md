# Perigord: Golang Tools for Ethereum Development

## Install

These instructions assume an Ubuntu 16.04 x86_64 environment.

### Golang

Some dependencies require Go 1.7+, but Go 1.6 is in Ubuntu 16.04's default repos.
The below will install Go 1.8.


```
$ sudo add-apt-repository -y ppa:longsleep/golang-backports
$ sudo apt-get update
$ sudo apt-get install -y golang-go
$ mkdir $HOME/golang
$ echo "export GOPATH=$HOME/golang:`pwd`" >> ~/.bashrc
$ echo "export PATH=$PATH:$HOME/golang/bin" >> ~/.bashrc
```

Close / re-open your terminal or re-`source` your `.bashrc`.


### Install `abigen`

```
$ go get github.com/ethereum/go-ethereum
$ pushd $HOME/golang/src/github.com/ethereum/go-ethereum
$ go install ./cmd/abigen
$ popd
```

### Install Dependencies

```
$ go get -u github.com/jteeuwen/go-bindata/...
$ go get -u github.com/spf13/cobra
$ go get -u github.com/spf13/viper
```

## Setup

```
$ go get -u github.com/swarmdotmarket/perigord
$ pushd $HOME/golang/src/github.com/ethereum/go-ethereum
$ go generate
# pushd perigord
$ go install
$ popd
$ popd
```

## Usage

Run for usage information:

```
$ perigord
```
