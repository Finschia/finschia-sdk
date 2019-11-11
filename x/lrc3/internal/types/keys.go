package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

const (
	ModuleName = "lrc3"

	StoreKey = ModuleName

	RouterKey = ModuleName
)

var (
	ApprovalsKeyPrefix         = []byte{0x02}
	OperatorApprovalsKeyPrefix = []byte{0x03}
)

func GetApprovalKey(denom string, tokenId string) []byte {
	h := tmhash.New()
	_, err := h.Write([]byte(tokenId))
	if err != nil {
		panic(err)
	}
	bs := h.Sum(nil)

	return append(append(ApprovalsKeyPrefix, []byte(denom)...), bs...)
}

func GetOperatorApprovalKey(denom string, ownerAddress sdk.AccAddress) []byte {
	h := tmhash.New()
	_, err := h.Write(ownerAddress)
	if err != nil {
		panic(err)
	}
	bs := h.Sum(nil)

	return append(append(OperatorApprovalsKeyPrefix, []byte(denom)...), bs...)
}
