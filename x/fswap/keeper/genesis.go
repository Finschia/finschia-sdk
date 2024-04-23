package keeper

import (
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/fswap/types"
)

// InitGenesis initializes the module's state from a provided genesis
// state.
func (k Keeper) InitGenesis(ctx sdk.Context, genState types.GenesisState) {
	if err := k.SetSwapped(ctx, genState.Swapped); err != nil {
		panic(err)
	}
	if err := k.SetTotalSupply(ctx, genState.SwappableNewCoinAmount); err != nil {
		panic(err)
	}
}

// ExportGenesis returns the capability module's exported genesis.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	return &types.GenesisState{
		Swapped:                k.GetSwapped(ctx),
		SwappableNewCoinAmount: k.GetTotalSupply(ctx),
	}
}
