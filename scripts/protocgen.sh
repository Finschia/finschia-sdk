#!/usr/bin/env bash

# How to run manually:
# docker build --pull --rm -f "contrib/devtools/Dockerfile" -t cosmossdk-proto:latest "contrib/devtools"
# docker run --rm -v $(pwd):/workspace --workdir /workspace cosmossdk-proto sh ./scripts/protocgen.sh

#echo "Formatting protobuf files"
#find ./ -name "*.proto" -exec clang-format -i {} \;

set -e

echo "Generating gogo proto code"
cd proto
proto_dirs=$(find ./cosmos ./lbm -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq | grep -v '^./cosmos/store/')
for dir in $proto_dirs; do
  for file in $(find "${dir}" -maxdepth 1 -name '*.proto'); do
    if grep -q "option go_package" "$file"; then
      buf generate --template buf.gen.gogo.yaml "$file"
    fi
  done
done

cd ..

# generate tests proto code
(cd testutil/testdata; buf generate)
(cd baseapp/testutil; buf generate)

# move proto files to the right places
cp -r github.com/Finschia/finschia-sdk/* ./
rm -rf github.com

go mod tidy

./scripts/protocgen-pulsar.sh
