package keeper

import (
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) InitGenesis(ctx sdk.Context, nonce math.Uint, ids []string) {
	k.setNonce(ctx, nonce)

	for _, id := range ids {
		k.addID(ctx, id)
	}
}

func (k Keeper) ExportGenesis(ctx sdk.Context) (nonce math.Uint, ids []string) {
	k.iterateIDs(ctx, func(id string) (stop bool) {
		ids = append(ids, id)
		return false
	})

	return k.getNonce(ctx), ids
}
