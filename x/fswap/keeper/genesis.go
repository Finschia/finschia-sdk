package keeper

import (
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/fswap/types"
)

// InitGenesis initializes the module's state from a provided genesis
// state.
func (k Keeper) InitGenesis(ctx sdk.Context, genState types.GenesisState) {
	if err := k.SetFswapInit(ctx, genState.FswapInit); err != nil {
		panic(err)
	}
	if err := k.SetSwapped(ctx, genState.Swapped); err != nil {
		panic(err)
	}
}

// ExportGenesis returns the capability module's exported genesis.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	return &types.GenesisState{
		FswapInit: k.GetFswapInit(ctx),
		Swapped:   k.GetSwapped(ctx),
	}
}
