package keeper

import (
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/fswap/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func (k Keeper) InitGenesis(ctx sdk.Context, genState *types.GenesisState) error {
	if len(genState.GetSwapInit()) > 1 {
		return types.ErrSwapCanNotBeInitializedTwice.Wrap("cannot initialize genesis state, there are more than 1 swapInit")
	}
	if len(genState.GetSwapped()) > 1 {
		return types.ErrSwapCanNotBeInitializedTwice.Wrap("cannot initialize genesis state, there are more than 1 swapped")
	}
	for _, swapInit := range genState.GetSwapInit() {
		if err := k.setSwapInit(ctx, swapInit); err != nil {
			panic(err)
		}
	}
	for _, swapped := range genState.GetSwapped() {
		if err := swapped.Validate(); err != nil {
			panic(err)
		}
	}
	return nil
}

// ExportGenesis returns the module's exported genesis.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	return &types.GenesisState{
		SwapInit: k.getAllSwapInits(ctx),
		Swapped:  k.getAllSwapped(ctx),
	}
}
