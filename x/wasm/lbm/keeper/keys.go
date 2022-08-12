package keeper

import sdk "github.com/line/lbm-sdk/types"

var (
	inactiveContractPrefix = []byte{0x90}
)

func getInactiveContractKey(contractAddress sdk.AccAddress) []byte {
	key := make([]byte, len(inactiveContractPrefix)+len(contractAddress))
	copy(key, inactiveContractPrefix)
	copy(key[len(inactiveContractPrefix):], contractAddress)
	return key
}
