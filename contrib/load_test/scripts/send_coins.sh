#!/usr/bin/env bash
set -e

if [ -z "$1" ]; then
  MASTER_ADDR=tlink1mrgyg2l98t9l2rur8dum7yqsv5mfwy8usws8a8
fi

./build/linkcli tx send $(./build/linkcli keys show alice -a) ${MASTER_ADDR} 100000000stake -y -b block
