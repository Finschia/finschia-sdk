package keeper

import (
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/fbridge/types"
)

func (k Keeper) InitGenesis(ctx sdk.Context, gs *types.GenesisState) error {
	k.setNextSequence(ctx, gs.SendingState.NextSeq)

	return nil
}

func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {

	return &types.GenesisState{
		SendingState: types.SendingState{
			NextSeq: k.GetNextSequence(ctx),
		},
	}
}
