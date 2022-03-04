package keeper

import (
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/token"
)

func (k Keeper) InitGenesis(ctx sdk.Context, data *token.ClassGenesisState) {
	k.setNonce(ctx, data.Nonce)

	for _, id := range data.Ids {
		k.addID(ctx, id)
	}
}

func (k Keeper) ExportGenesis(ctx sdk.Context) *token.ClassGenesisState {
	var ids []string
	k.iterateIDs(ctx, func(id string) (stop bool) {
		ids = append(ids, id)
		return false
	})

	return &token.ClassGenesisState{
		Nonce: k.getNonce(ctx),
		Ids:   ids,
	}
}
