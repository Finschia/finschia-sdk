package internal

import (
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/foundation"
)

func (k Keeper) CollectFoundationTax(ctx sdk.Context) error {
	// fetch and clear the collected fees for the fund, since this is
	// called in BeginBlock, collected fees will be from the previous block
	feeCollector := k.authKeeper.GetModuleAccount(ctx, k.feeCollectorName)
	feesCollectedInt := k.bankKeeper.GetAllBalances(ctx, feeCollector.GetAddress())
	feesCollected := sdk.NewDecCoinsFromCoins(feesCollectedInt...)

	// calculate the tax
	taxRatio := k.GetFoundationTax(ctx)
	tax, _ := feesCollected.MulDecTruncate(taxRatio).TruncateDecimal()

	// update foundation treasury
	pool := k.GetPool(ctx)
	pool.Treasury = pool.Treasury.Add(sdk.NewDecCoinsFromCoins(tax...)...)
	k.SetPool(ctx, pool)

	// collect tax to the foundation treasury
	if err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, k.feeCollectorName, foundation.TreasuryName, tax); err != nil {
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
