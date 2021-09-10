package keeper

import (
	"fmt"

	sdk "github.com/line/lbm-sdk/types"
	wasmvm "github.com/line/wasmvm"
)

type cosmwasmAPIImpl struct {
	gasMultiplier uint64
}

func (a cosmwasmAPIImpl) humanAddress(canon []byte) (string, uint64, error) {
	gas := 5 * a.gasMultiplier
	if len(canon) != sdk.BytesAddrLen {
		//nolint:stylecheck
		return "", gas, fmt.Errorf("expected %d byte address", sdk.BytesAddrLen)
	}

	return sdk.BytesToAccAddress(canon).String(), gas, nil
}

func (a cosmwasmAPIImpl) canonicalAddress(human string) ([]byte, uint64, error) {
	bz, err := sdk.AccAddressToBytes(human)
	return bz, 4 * a.gasMultiplier, err
}

func (k Keeper) cosmwasmAPI(ctx sdk.Context) wasmvm.GoAPI {
	x := cosmwasmAPIImpl{
		gasMultiplier: k.getGasMultiplier(ctx),
	}
	return wasmvm.GoAPI{
		HumanAddress:     x.humanAddress,
		CanonicalAddress: x.canonicalAddress,
	}
}
