#!/bin/sh

set -e

# cosmovisor configs
[ -n "$COSMOVISOR_TAG" ]

# app configs
CHAIN_ID=sim
KEY_MNEMONIC="mind flame tobacco sense move hammer drift crime ring globe art gaze cinnamon helmet cruise special produce notable negative wait path scrap recall have"

# set DEBUG to non empty to retain workdir
cleanup() {
	if [ -n "$workdir" ]
	then
		if [ -n "$DEBUG" ]
		then
			echo "workdir at: $workdir"
			find $workdir
		else
			rm -rf $workdir
		fi
	fi
}
trap cleanup TERM INT EXIT

export DAEMON_NAME=simd

# make a temporary working directory
workdir=$(mktemp -d)
export DAEMON_HOME=$workdir

COSMOVISOR_TAG=$COSMOVISOR_TAG sh install.sh
CHAIN_ID=$CHAIN_ID KEY_MNEMONIC="$KEY_MNEMONIC" KEY_INDEX=0 sh configure.sh

result=$workdir/result.fifo
mkfifo $result

RESULT=$result sh start.sh &
CHAIN_ID=$CHAIN_ID sh upgrade.sh

case $(cat $result) in
	OK)
	;;
	*)
		false
	;;
esac
