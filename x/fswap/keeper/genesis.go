package keeper

import (
	"fmt"

	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/fswap/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func (k Keeper) InitGenesis(ctx sdk.Context, genState *types.GenesisState) error {
	if len(genState.GetFswapInit()) > 1 {
		return fmt.Errorf("cannot initialize genesis state, there are more than 1 fswapInit")
	}
	if len(genState.GetSwapped()) > 1 {
		return fmt.Errorf("cannot initialize genesis state, there are more than 1 swapped")
	}
	for _, fswapInit := range genState.GetFswapInit() {
		if err := k.SetFswapInit(ctx, fswapInit); err != nil {
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
		FswapInit: k.GetAllFswapInits(ctx),
		Swapped:   k.GetAllSwapped(ctx),
	}
}
