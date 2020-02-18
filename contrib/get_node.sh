#!/usr/bin/env bash
set -xe
VERSION=v11.15.0
if [[ "$OSTYPE" == "darwin"* ]]; then
  OS=darwin
else
  OS=linux
fi
NODE_FULL=node-${VERSION}-${OS}-x64

echo "Get node"
rm -rf ~/.local/bin
rm -rf ~/.local/node
mkdir -p ~/.local/bin
mkdir -p ~/.local/node
wget http://nodejs.org/dist/${VERSION}/${NODE_FULL}.tar.gz -O ~/.local/node/${NODE_FULL}.tar.gz
tar -xzf ~/.local/node/${NODE_FULL}.tar.gz -C ~/.local/node/
ln -s ~/.local/node/${NODE_FULL}/bin/node ~/.local/bin/node
ln -s ~/.local/node/${NODE_FULL}/bin/npm ~/.local/bin/npm
npm i -g dredd@12.1.0
ln -s ~/.local/node/${NODE_FULL}/bin/dredd ~/.local/bin/dredd
npm i -g swagger-cli@3.0.1
ln -s ~/.local/node/${NODE_FULL}/bin/swagger-cli ~/.local/bin/swagger-cli
