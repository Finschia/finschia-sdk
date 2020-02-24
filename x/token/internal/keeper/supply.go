package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/token/internal/types"
)

type SupplyKeeper interface {
	GetSupplyInt(ctx sdk.Context, symbol string) (sdk.Int, sdk.Error)
	MintSupply(ctx sdk.Context, symbol string, to sdk.AccAddress, amount sdk.Int) sdk.Error
	BurnSupply(ctx sdk.Context, symbol string, from sdk.AccAddress, amount sdk.Int) sdk.Error
}

var _ SupplyKeeper = (*Keeper)(nil)

func (k Keeper) GetSupplyInt(ctx sdk.Context, symbol string) (sdk.Int, sdk.Error) {
	supply, err := k.getSupply(ctx, symbol)
	if err != nil {
		return sdk.ZeroInt(), err
	}
	return supply.GetTotal(), nil
}
func (k Keeper) getSupply(ctx sdk.Context, symbol string) (supply types.Supply, err sdk.Error) {
	if _, err := k.GetToken(ctx, symbol); err != nil {
		return nil, err
	}
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.SupplyKey(symbol))
	if b == nil {
		panic("stored supply should not have been nil")
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &supply)
	return
}

func (k Keeper) setSupply(ctx sdk.Context, supply types.Supply) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(supply)
	store.Set(types.SupplyKey(supply.GetSymbol()), b)
}

func (k Keeper) MintSupply(ctx sdk.Context, symbol string, to sdk.AccAddress, amount sdk.Int) sdk.Error {
	_, err := k.addBalance(ctx, symbol, to, amount)
	if err != nil {
		return err
	}

	supply, err := k.getSupply(ctx, symbol)
	if err != nil {
		return err
	}
	oldSupplyAmount := supply.GetTotal()
	newSupplyAmount := oldSupplyAmount.Add(amount)
	if newSupplyAmount.IsNegative() {
		return types.ErrInsufficientSupply(types.DefaultCodespace, fmt.Sprintf("insufficient supply for token [%s]; %s < %s", symbol, oldSupplyAmount, amount))
	}
	supply = supply.Inflate(amount)
	k.setSupply(ctx, supply)

	return nil
}

func (k Keeper) BurnSupply(ctx sdk.Context, symbol string, from sdk.AccAddress, amount sdk.Int) sdk.Error {
	_, err := k.subtractBalance(ctx, symbol, from, amount)
	if err != nil {
		return err
	}

	supply, err := k.getSupply(ctx, symbol)
	if err != nil {
		return err
	}
	oldSupplyAmount := supply.GetTotal()
	newSupplyAmount := oldSupplyAmount.Sub(amount)
	if newSupplyAmount.IsNegative() {
		return types.ErrInsufficientSupply(types.DefaultCodespace, fmt.Sprintf("insufficient supply for token [%s]; %s < %s", symbol, oldSupplyAmount, amount))
	}
	supply = supply.Deflate(amount)
	k.setSupply(ctx, supply)

	return nil
}
