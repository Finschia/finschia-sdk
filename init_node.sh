#!/bin/sh

display_usage() {
	echo "\nMissing $1 parameter. Please check if all parameters were specified."
	echo "\nUsage: ./run_node [CHAIN_ID]"
  echo "\nExample: ./init_node $BINARY test-chain-id 26657 26656 6060 9090 \n"
  exit 1
}

KEYRING=--keyring-backend="test"
SILENT=1

redirect() {
  if [ "$SILENT" -eq 1 ]; then
    "$@" > /dev/null 2>&1
  else
    "$@"
  fi
}

BINARY=simd
CHAINID=$1
CHAINDIR=~/.simapp
RPCPORT=26657
P2PPORT=26656
PROFPORT=6060
GRPCPORT=9090

if [ -z "$1" ]; then
  display_usage "[CHAIN_ID]"
fi

echo "Creating $BINARY instance: home=$CHAINDIR | chain-id=$CHAINID | p2p=:$P2PPORT | rpc=:$RPCPORT | profiling=:$PROFPORT | grpc=:$GRPCPORT"

# Add dir for chain, exit if error
if ! mkdir -p $CHAINDIR 2>/dev/null; then
    echo "Failed to create chain folder. Aborting..."
    exit 1
fi

# Build genesis file incl account for passed address
coins="100000000000stake,100000000000ukrw"
delegate="100000000000stake"

redirect $BINARY --home $CHAINDIR --chain-id $CHAINID init $CHAINID 
sleep 1
$BINARY --home $CHAINDIR keys add validator $KEYRING --output json > $CHAINDIR/validator_seed.json 2> /dev/null
sleep 1
$BINARY --home $CHAINDIR keys add user $KEYRING --recover --output json < user_mnemonic > $CHAINDIR/key_seed.json 2> /dev/null
sleep 1
redirect $BINARY --home $CHAINDIR add-genesis-account $($BINARY --home $CHAINDIR keys $KEYRING show user -a) $coins 
sleep 1
redirect $BINARY --home $CHAINDIR add-genesis-account $($BINARY --home $CHAINDIR keys $KEYRING show validator -a) $coins 
sleep 1
redirect $BINARY --home $CHAINDIR gentx validator $delegate $KEYRING --chain-id $CHAINID
sleep 1
redirect $BINARY --home $CHAINDIR collect-gentxs 
sleep 1

# Check platform
platform='unknown'
unamestr=`uname`
if [ "$unamestr" = 'Linux' ]; then
   platform='linux'
fi

# Set proper defaults and change ports (use a different sed for Mac or Linux)
echo "Change settings in config.toml file..."
if [ $platform = 'linux' ]; then
  sed -i 's#"tcp://127.0.0.1:26657"#"tcp://0.0.0.0:'"$RPCPORT"'"#g' $CHAINDIR/config/config.toml
  sed -i 's#"tcp://0.0.0.0:26656"#"tcp://0.0.0.0:'"$P2PPORT"'"#g' $CHAINDIR/config/config.toml
  sed -i 's#"localhost:6060"#"localhost:'"$P2PPORT"'"#g' $CHAINDIR/config/config.toml
  sed -i 's/timeout_commit = "5s"/timeout_commit = "1s"/g' $CHAINDIR/config/config.toml
  sed -i 's/timeout_propose = "3s"/timeout_propose = "1s"/g' $CHAINDIR/config/config.toml
  sed -i 's/index_all_keys = false/index_all_keys = true/g' $CHAINDIR/config/config.toml
  sed -i 's/pruning = "default"/pruning = "nothing"/g' $CHAINDIR/config/app.toml
  # sed -i '' 's#index-events = \[\]#index-events = \["message.action","send_packet.packet_src_channel","send_packet.packet_sequence"\]#g' $CHAINDIR/config/app.toml
else
  sed -i '' 's#"tcp://127.0.0.1:26657"#"tcp://0.0.0.0:'"$RPCPORT"'"#g' $CHAINDIR/config/config.toml
  sed -i '' 's#"tcp://0.0.0.0:26656"#"tcp://0.0.0.0:'"$P2PPORT"'"#g' $CHAINDIR/config/config.toml
  sed -i '' 's#"localhost:6060"#"localhost:'"$P2PPORT"'"#g' $CHAINDIR/config/config.toml
  sed -i '' 's/timeout_commit = "5s"/timeout_commit = "1s"/g' $CHAINDIR/config/config.toml
  sed -i '' 's/timeout_propose = "3s"/timeout_propose = "1s"/g' $CHAINDIR/config/config.toml
  sed -i '' 's/index_all_keys = false/index_all_keys = true/g' $CHAINDIR/config/config.toml
  sed -i '' 's/pruning = "default"/pruning = "nothing"/g' $CHAINDIR/config/app.toml
  # sed -i '' 's#index-events = \[\]#index-events = \["message.action","send_packet.packet_src_channel","send_packet.packet_sequence"\]#g' $CHAINDIR/config/app.toml
fi

