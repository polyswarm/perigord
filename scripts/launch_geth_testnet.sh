#!/bin/bash

DIR=/tmp/geth_private_testnet

# cd to script directory
mkdir -p $DIR
cd "${0%/*}"

cp ./genesis.json $DIR
cp -r ./keystore $DIR
mkdir -p $DIR/etc
echo -n blah > $DIR/etc/pw

geth --datadir $DIR --nodiscover --maxpeers 0 init $DIR/genesis.json
geth --datadir $DIR --nodiscover --maxpeers 0 --mine --minerthreads 1 --rpc --rpcapi "eth,web3,personal,net" console
