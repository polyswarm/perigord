FROM ubuntu:16.04
LABEL maintainer="Maxwell Koo <mjkoo90@gmail.com>"

RUN set -x && \
    apt-get update && \
    apt-get -y install software-properties-common && \
    apt-add-repository -y ppa:longsleep/golang-backports && \
    apt-add-repository -y ppa:ethereum/ethereum && \
    apt-get update && \
    apt-get -y install \
        git \
        golang-go \
        solc && \
    apt-get -y purge software-properties-common && \
    apt-get -y autoremove && \
    rm -rf /var/lib/apt/lists/*

ENV GOPATH /go
ENV PATH $GOPATH/bin:$PATH

RUN set -x && \
    mkdir -p $GOPATH && \
    go get -u github.com/ethereum/go-ethereum && \
    cd $GOPATH/src/github.com/ethereum/go-ethereum/cmd/abigen && \
    go install && \
    go get -u github.com/jteeuwen/go-bindata/... && \
    (go get -u github.com/polyswarm/perigord/... || true) && \
    cd $GOPATH/src/github.com/polyswarm/perigord && \
    go generate && \
    cd perigord && \
    go install
