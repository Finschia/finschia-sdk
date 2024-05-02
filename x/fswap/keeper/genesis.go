package keeper

import (
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/fswap/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func (k Keeper) InitGenesis(ctx sdk.Context, genState *types.GenesisState) error {
	if err := genState.Validate(); err != nil {
		return err
	}

	if len(genState.GetSwaps()) > k.config.MaxSwaps && !k.isUnlimited() {
		return types.ErrCanNotHaveMoreSwap.Wrapf("cannot initialize genesis state, there are more than %d swaps", k.config.MaxSwaps)
	}

	if len(genState.GetSwappeds()) > k.config.MaxSwaps && !k.isUnlimited() {
		return types.ErrCanNotHaveMoreSwap.Wrapf("cannot initialize genesis state, there are more than %d swapped", k.config.MaxSwaps)
	}

	if err := k.setSwapStats(ctx, genState.GetSwapStats()); err != nil {
		return err
	}

	for _, swap := range genState.GetSwaps() {
		if err := k.setSwap(ctx, swap); err != nil {
			panic(err)
		}
	}

	for _, swapped := range genState.GetSwappeds() {
		err := k.setSwapped(ctx, swapped)
		if err != nil {
			return err
		}
	}

	return nil
}

func (k Keeper) isUnlimited() bool {
	return k.config.MaxSwaps == 0
}

// ExportGenesis returns the module's exported genesis.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	stats, err := k.getSwapStats(ctx)
	if err != nil {
		panic(err)
	}
	return &types.GenesisState{
		Swaps:     k.getAllSwaps(ctx),
		SwapStats: stats,
		Swappeds:  k.getAllSwapped(ctx),
	}
}
