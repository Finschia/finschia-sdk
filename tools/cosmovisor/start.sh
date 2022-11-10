#!/bin/sh

set -e

[ -n "$DAEMON_HOME" ]
[ -d $DAEMON_HOME ]
[ -n "$RESULT" ]
[ -p $RESULT ]

if cosmovisor run start --home $DAEMON_HOME
then
	echo OK >>$RESULT
else
	echo FAILED >>$RESULT
fi
