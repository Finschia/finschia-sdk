#!/usr/bin/env sh
set -ex

export LINKCLI="docker run -i --net=host -v ${HOME}/.linkd:/root/.linkd -v ${HOME}/.linkcli:/root/.linkcli line/link linkcli"
export LINKD="docker run -i -p 26656:26656 -p 26657:26657 -v ${HOME}/.linkd:/root/.linkd -v ${HOME}/.linkcli:/root/.linkcli line/link linkd"

./initialize.sh
