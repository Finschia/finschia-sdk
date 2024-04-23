package keeper

import (
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/fswap/types"
)

// InitGenesis initializes the module's state from a provided genesis
// state.
func (k Keeper) InitGenesis(ctx sdk.Context, bk types.BankKeeper, genState types.GenesisState) {
	if err := k.SetSwapped(ctx, genState.Swapped); err != nil {
		panic(err)
	}
	totalOldCoinsSupply := bk.GetSupply(ctx, k.config.OldCoinDenom).Amount
	totalNewCoinsSupply := k.config.SwapRate.MulInt(totalOldCoinsSupply)
	totalNewCoins := sdk.NewDecCoinFromDec(k.config.NewCoinDenom, totalNewCoinsSupply)
	if err := k.SetTotalSupply(ctx, totalNewCoins); err != nil {
		panic(err)
	}
}

// ExportGenesis returns the capability module's exported genesis.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	return &types.GenesisState{
		Swapped: k.GetSwapped(ctx),
	}
}
