#!/bin/bash
L2BINARYNAME=rollupd
L1_KEYRING_DIR=~/.simapp
L2_KEYRING_DIR=~/.l2simapp

# Remove L1 data
pkill simd-darwin-arm64
rm -rf $L1_KEYRING_DIR

# Remove L2 data
if [ -x "$(which ${L2BINARYNAME})" ]; then
    echo "clear L2 Binary"
    pkill $L2BINARYNAME
    rm $GOPATH/bin/$L2BINARYNAME
fi

rm -rf $L2_KEYRING_DIR
