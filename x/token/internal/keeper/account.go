package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetAccountBalance(ctx sdk.Context, symbol string, addr sdk.AccAddress) sdk.Int {
	return k.bankKeeper.GetCoins(ctx, addr).AmountOf(symbol)
}
