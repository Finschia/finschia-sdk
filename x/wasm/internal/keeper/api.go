package keeper

import (
	"fmt"

	cosmwasm "github.com/CosmWasm/wasmvm"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/wasm/internal/types"
)

type cosmwasmAPIImpl struct {
	humanizeCost     uint64
	canonicalizeCost uint64
}

func (a cosmwasmAPIImpl) humanAddress(canon []byte) (string, uint64, error) {
	if len(canon) != sdk.AddrLen {
		return "", a.humanizeCost, fmt.Errorf("expected %d byte address", sdk.AddrLen)
	}

	return sdk.AccAddress(canon).String(), types.DefaultHumanizeCost, nil
}

func (a cosmwasmAPIImpl) canonicalAddress(human string) ([]byte, uint64, error) {
	bz, err := sdk.AccAddressFromBech32(human)
	return bz, a.canonicalizeCost, err
}

func (k Keeper) cosmwasmAPI(ctx sdk.Context) cosmwasm.GoAPI {
	x := cosmwasmAPIImpl{
		humanizeCost:     k.getHumanizeCost(ctx),
		canonicalizeCost: k.getCanonicalCost(ctx),
	}
	return cosmwasm.GoAPI{
		HumanAddress:     x.humanAddress,
		CanonicalAddress: x.canonicalAddress,
	}
}
