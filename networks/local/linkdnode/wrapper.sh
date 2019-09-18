#!/usr/bin/env sh

##
## Input parameters
##
BINARY=/linkd/${BINARY:-linkd}
ID=${ID:-0}
LOG=${LOG:-linkd.log}

##
## Assert linux binary
##
if ! [ -f "${BINARY}" ]; then
	echo "The binary $(basename "${BINARY}") cannot be found. Please add the binary to the shared folder. Please use the BINARY environment variable if the name of the binary is not 'linkd' E.g.: -e BINARY=linkd_my_test_version"
	exit 1
fi
BINARY_CHECK="$(file "$BINARY" | grep 'ELF 64-bit LSB executable, x86-64')"
if [ -z "${BINARY_CHECK}" ]; then
	echo "Binary needs to be OS linux, ARCH amd64"
	exit 1
fi

##
## Run binary with all parameters
##
export LINKDHOME="${LINKDHOME:-/linkd/node${ID}/linkd}"

if [ -d "`dirname ${LINKDHOME}/${LOG}`" ]; then
  "$BINARY" --home "$LINKDHOME" "$@" | tee "${LINKDHOME}/${LOG}"
else
  "$BINARY" --home "$LINKDHOME" "$@"
fi

chmod 777 -R /linkd

