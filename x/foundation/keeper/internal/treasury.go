package internal

import (
	sdk "github.com/Finschia/finschia-sdk/types"
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
	"github.com/Finschia/finschia-sdk/x/foundation"
)

func (k Keeper) CollectFoundationTax(ctx sdk.Context) error {
	feeCollector := k.authKeeper.GetModuleAccount(ctx, k.feeCollectorName).GetAddress()
	feesCollectedInt := k.bankKeeper.GetAllBalances(ctx, feeCollector)
	feesCollected := sdk.NewDecCoinsFromCoins(feesCollectedInt...)

	// calculate the tax
	taxRatio := k.GetFoundationTax(ctx)
	tax, _ := feesCollected.MulDecTruncate(taxRatio).TruncateDecimal()

	// collect the tax
	if err := k.FundTreasury(ctx, feeCollector, tax); err != nil {
		return err
	}

	return nil
}

func (k Keeper) GetTreasury(ctx sdk.Context) sdk.DecCoins {
	return k.GetPool(ctx).Treasury
}

func (k Keeper) FundTreasury(ctx sdk.Context, from sdk.AccAddress, amt sdk.Coins) error {
	pool := k.GetPool(ctx)
	pool.Treasury = pool.Treasury.Add(sdk.NewDecCoinsFromCoins(amt...)...)
	k.SetPool(ctx, pool)

	return k.bankKeeper.SendCoinsFromAccountToModule(ctx, from, foundation.TreasuryName, amt)
}

func (k Keeper) WithdrawFromTreasury(ctx sdk.Context, to sdk.AccAddress, amt sdk.Coins) error {
	pool := k.GetPool(ctx)
	remains, hasNeg := pool.Treasury.SafeSub(sdk.NewDecCoinsFromCoins(amt...))
	if hasNeg {
		return sdkerrors.ErrInsufficientFunds.Wrapf("not enough coins in treasury, %s", pool.Treasury)
	}
	pool.Treasury = remains
	k.SetPool(ctx, pool)

	return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, foundation.TreasuryName, to, amt)
}

func (k Keeper) GetPool(ctx sdk.Context) foundation.Pool {
	store := ctx.KVStore(k.storeKey)
	key := poolKey
	bz := store.Get(key)

	var pool foundation.Pool
	k.cdc.MustUnmarshal(bz, &pool)

	return pool
}

func (k Keeper) SetPool(ctx sdk.Context, pool foundation.Pool) {
	bz := k.cdc.MustMarshal(&pool)

	store := ctx.KVStore(k.storeKey)
	key := poolKey
	store.Set(key, bz)
}
