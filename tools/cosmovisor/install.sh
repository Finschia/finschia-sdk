#!/bin/sh

set -e

# cosmovisor config assertions
[ -n "$COSMOVISOR_TAG" ]
[ -n "$DAEMON_HOME" ]
[ -d $DAEMON_HOME ]
[ -n "$DAEMON_NAME" ]

# install app
go install ./../../...

# install cosmovisor
go install cosmossdk.io/tools/cosmovisor/cmd/cosmovisor@$COSMOVISOR_TAG
cosmovisor init $(which $DAEMON_NAME)
