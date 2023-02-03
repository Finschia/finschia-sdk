#!/bin/sh

set -e

[ -n "$DAEMON_HOME" ]
[ -d $DAEMON_HOME ]
[ -n "$DAEMON_NAME" ]
[ -x $(which $DAEMON_NAME) ]

[ -n "$CHAIN_ID" ]

UPGRADE_HEIGHT=5

assert_begin() {
	info=$($DAEMON_NAME q block --log_level info 2>/dev/null)
	if [ -z "$info" ]
	then
		return 1
	fi

	height=$(echo "$info" | jq -r '.block.header.height')
	if [ "$height" -ge 0 ] 2>/dev/null
	then
		return 0
	else
		return 1
	fi
}

wait_for_begin() {
	while ! assert_begin
	do
		sleep 1
	done
}

keyring_backend=test

new_binary=$(realpath ./dummyd)
prepare_upgrade() {
	name=$1
	bindir=$DAEMON_HOME/cosmovisor/upgrades/$name/bin
	mkdir -p $bindir
	cp $new_binary $bindir/$DAEMON_NAME
}

submit_upgrade() {
	name=$1
	checksum=sha256:$(sha256sum $new_binary | awk '{print $1}')
	info='{"binaries":{"any":"file://'$new_binary'?checksum='$checksum'"}}'
	$DAEMON_NAME --home $DAEMON_HOME tx --keyring-backend $keyring_backend gov submit-proposal software-upgrade $name --upgrade-height $UPGRADE_HEIGHT --upgrade-info $info --title upgrade --description "test upgrade" --deposit 1stake --broadcast-mode block --from validator0 --chain-id $CHAIN_ID --yes
}

vote() {
	proposal=$1
	$DAEMON_NAME --home $DAEMON_HOME tx --keyring-backend $keyring_backend gov vote $proposal VOTE_OPTION_YES --broadcast-mode sync --from validator0 --chain-id $CHAIN_ID --yes
}

wait_for_begin

if [ "$DAEMON_ALLOW_DOWNLOAD_BINARIES" != true ]
then
	prepare_upgrade testing
fi
submit_upgrade testing
vote 1
