#!/usr/bin/env sh

##
## Input parameters
##
LINKDHOME=${LINKDHOME:-/testhome}

linkd --home "$LINKDHOME" "$@"

chmod 777 -R /linkd

