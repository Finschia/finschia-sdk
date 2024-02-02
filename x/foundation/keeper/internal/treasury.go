package internal

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/Finschia/finschia-sdk/x/foundation"
)

func (k Keeper) CollectFoundationTax(ctx sdk.Context) error {
	// TODO(@0Tech): use auth keeper after applying global bech32 removal
	feeCollector := address.Module(k.feeCollectorName)
	feesCollectedInt := k.bankKeeper.GetAllBalances(ctx, feeCollector)
	if feesCollectedInt.Empty() {
		return nil
	}
	feesCollected := sdk.NewDecCoinsFromCoins(feesCollectedInt...)

	// calculate the tax
	taxRatio := k.GetFoundationTax(ctx)
	tax, _ := feesCollected.MulDecTruncate(taxRatio).TruncateDecimal()
	if tax.Empty() {
		return nil
	}

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
	store := k.storeService.OpenKVStore(ctx)
	key := poolKey
	bz, err := store.Get(key)
	if err != nil {
		panic(err)
	}

	var pool foundation.Pool
	k.cdc.MustUnmarshal(bz, &pool)

	return pool
}

func (k Keeper) SetPool(ctx sdk.Context, pool foundation.Pool) {
	bz := k.cdc.MustMarshal(&pool)

	store := k.storeService.OpenKVStore(ctx)
	key := poolKey
	if err := store.Set(key, bz); err != nil {
		panic(err)
	}
}
