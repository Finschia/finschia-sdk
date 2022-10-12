package keeper

import (
	gogotypes "github.com/gogo/protobuf/types"

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
	tax := feesCollected.MulDecTruncate(taxRatio)

	// update foundation treasury
	pool := k.GetPool(ctx)
	pool.Treasury = pool.Treasury.Add(tax...)
	k.SetPool(ctx, pool)

	// collect tax to the foundation treasury
	amount, _ := tax.TruncateDecimal()
	if err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, k.feeCollectorName, foundation.TreasuryName, amount); err != nil {
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

func (k Keeper) GetOneTimeMintCount(ctx sdk.Context) uint32 {
	store := ctx.KVStore(k.storeKey)
	key := oneTimeMintKey
	bz := store.Get(key)

	var count gogotypes.UInt32Value
	k.cdc.MustUnmarshal(bz, &count)

	return count.Value
}

func (k Keeper) SetOneTimeMintCount(ctx sdk.Context, count uint32) {
	store := ctx.KVStore(k.storeKey)
	key := oneTimeMintKey
	bz := k.cdc.MustMarshal(&gogotypes.UInt32Value{Value: count})
	store.Set(key, bz)
}
