#!/bin/bash

# Prepare chain
rm -rf ~/.simapp
zsh ../../../init_node.sh sim

sleep 2

simd start --home ~/.simapp/simapp0
