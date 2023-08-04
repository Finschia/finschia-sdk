#!/bin/bash

# functions
start_challenge() {
  # tools/bin/cannon load-elf --path tools/bin/mini.elf --out tools/miniapp/state.json --meta tools/miniapp/meta.json > /dev/null 2>&1
  tools/bin/cannon run --input tools/miniapp/state.json --meta tools/miniapp/meta.json --output out/output.json --proof-at =0 --proof-fmt $PROOF_FMT tools/bin/preimage tools/preimage/data/height-$BLOCK_HEIGHT-cha.json > /dev/null 2>&1
  STEP=$(jq -r .step out/output.json)
  echo_tx_link $(simd tx settlement start-challenge $CHALLENGER $DEFENDER $ROLLUP $BLOCK_HEIGHT $STEP --keyring-backend=test --home $BASE_DIR --chain-id $CHAIN_ID --from $CHALLENGER -y | jq -r .txhash)
  sleep 10
  CHALLENGE_ID=$(tools/bin/challenge id-from $ROLLUP $BLOCK_HEIGHT $CHALLENGER $DEFENDER)
  simd query settlement challenge $CHALLENGE_ID --home $BASE_DIR --output json | jq .challenge > out/challenge-0.json
}

nsect_challenge() {
  PREROUND=$((ROUND-1))
  STEPS=$(tools/bin/challenge steps --path out/challenge-$PREROUND.json)
  CHLG_HASHES=
  DEFD_HASHES=
  for STEP in $STEPS; do
    tools/bin/cannon run --input tools/miniapp/state.json --meta tools/miniapp/meta.json --output out/output.json --proof-at =$STEP --proof-fmt $PROOF_FMT tools/bin/preimage tools/preimage/data/height-$BLOCK_HEIGHT-cha.json > /dev/null 2>&1
    tools/bin/cannon run --input tools/miniapp/state.json --meta tools/miniapp/meta.json --output out/output.json --proof-at =$STEP --proof-fmt $PROOF_MAL_FMT tools/bin/preimage tools/preimage/data/height-$BLOCK_HEIGHT-def.json > /dev/null 2>&1
    if [[ ! $CHLG_HASHES ]]; then
        CHLG_HASHES=$(jq -r .pre out/proof-$STEP.json)
        DEFD_HASHES=$(jq -r .pre out/proof-$STEP-mal.json)
    else
        CHLG_HASHES=$CHLG_HASHES,$(jq -r .pre out/proof-$STEP.json)
        DEFD_HASHES=$DEFD_HASHES,$(jq -r .pre out/proof-$STEP-mal.json)
    fi
  done

  OPTIONS=$1
  echo_tx_link $(simd tx settlement nsect-challenge $CHALLENGER $CHALLENGE_ID $CHLG_HASHES --keyring-backend=test --home $BASE_DIR --chain-id $CHAIN_ID --from $CHALLENGER $OPTIONS -y | jq -r .txhash)
  sleep 10

  echo_tx_link $(simd tx settlement nsect-challenge $DEFENDER $CHALLENGE_ID $DEFD_HASHES --keyring-backend=test --home $BASE_DIR --chain-id $CHAIN_ID --from $DEFENDER $OPTIONS -y | jq -r .txhash)
  sleep 10

  simd query settlement challenge $CHALLENGE_ID --home $BASE_DIR --output json | jq .challenge > out/challenge-$ROUND.json
}

finish_challenge() {
  OPTIONS=$1
  STEP=$(jq -r .l out/challenge-$ROUND.json)
  STATE=$(jq -r .witness.state out/proof-$STEP.json)
  PROOFS=$(jq -r .witness.mem_proof out/proof-$STEP.json)
  PIKEY=$(jq -r .witness.preimage_key out/proof-$STEP.json)
  PIVAL=$(jq -r .witness.preimage_value out/proof-$STEP.json)
  PIOFFSET=$(jq -r .witness.preimage_offset out/proof-$STEP.json)
  TXID=$(simd tx settlement finish-challenge $CHALLENGER $CHALLENGE_ID $STATE $PROOFS $PIKEY,$PIVAL,$PIOFFSET --keyring-backend=test --home $BASE_DIR --chain-id $CHAIN_ID --from $CHALLENGER $OPTIONS -y | jq -r .txhash)
  sleep 10
  TXLINK=http://localhost:26657/tx?hash=0x$TXID
  echo $TXLINK
  WIN=$(curl -s $TXLINK | jq -r .result.tx_result.events[7].attributes[1].value | base64 --decode)
  if [ $WIN = true ]; then
    echo "Challenger wins:) demo successful:)"
  else
    echo "Challenger loses:( demo failure:("
  fi
}

echo_tx_link() {
  echo http://localhost:26657/tx?hash=0x$1
}

chekc_binary() {
  if ! [ -f "tools/bin/cannon" ]; then
    echo "Error: cannon doesn't exist. Try running 'cd tools && make all'" >&2
    exit 1
  fi
  if ! [ -f "tools/bin/challenge" ]; then
    echo "Error: challenge doesn't exist. Try running 'cd tools && make all'" >&2
    exit 1
  fi
  if ! [ -f "tools/bin/mini.elf" ]; then
    echo "Error: mini.elf doesn't exist. Try running 'cd tools && make all'" >&2
    exit 1
  fi
  if ! [ -f "tools/bin/preimage" ]; then
    echo "Error: preimage doesn't exist. Try running 'cd tools && make all'" >&2
    exit 1
  fi
}

# check binary
chekc_binary

# settings
BASE_DIR=~/.simapp
CHAIN_ID=sim
DEFENDER=$(simd keys show sequencer --keyring-backend=test --home $BASE_DIR --output json | jq -r .address)
CHALLENGER=$(simd keys show challenger --keyring-backend=test --home $BASE_DIR --output json | jq -r .address)
ROLLUP=test-rollup
BLOCK_HEIGHT=100
CHALLENGE_ID=
ROUND=
PROOF_FMT=out/proof-%d.json
PROOF_MAL_FMT=out/proof-%d-mal.json

# make output directory
mkdir -p out

echo "# Start challenge"
start_challenge

echo "# Nsect challenge"
IS_SEARCHING=$(tools/bin/challenge is-searching --path out/challenge-0.json)
ROUND=0
while [ "$IS_SEARCHING" == "true" ]
do
  ROUND=$((ROUND+1))
  echo "## Round $ROUND"
  nsect_challenge "--gas 300000"
  IS_SEARCHING=$(tools/bin/challenge is-searching --path out/challenge-$ROUND.json)
done

echo "# Finish challenge"
finish_challenge "--gas 300000"
