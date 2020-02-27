package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/collection/internal/types"
)

type SupplyKeeper interface {
	GetTotalInt(ctx sdk.Context, symbol, tokenID, target string) (supply sdk.Int, err sdk.Error)
	GetSupply(ctx sdk.Context, symbol string) (supply types.Supply)
	SetSupply(ctx sdk.Context, supply types.Supply)
	MintSupply(ctx sdk.Context, symbol string, to sdk.AccAddress, amt types.Coins) sdk.Error
	BurnSupply(ctx sdk.Context, symbol string, from sdk.AccAddress, amt types.Coins) sdk.Error
}

var _ SupplyKeeper = (*Keeper)(nil)

func (k Keeper) GetSupply(ctx sdk.Context, symbol string) (supply types.Supply) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.SupplyKey(symbol))
	if b == nil {
		panic("stored supply should not have been nil")
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &supply)
	return
}

func (k Keeper) SetSupply(ctx sdk.Context, supply types.Supply) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(supply)
	store.Set(types.SupplyKey(supply.GetSymbol()), b)
}

func (k Keeper) GetTotalInt(ctx sdk.Context, symbol, tokenID, target string) (supply sdk.Int, err sdk.Error) {
	if _, err = k.GetToken(ctx, symbol, tokenID); err != nil {
		return sdk.NewInt(0), err
	}

	s := k.GetSupply(ctx, symbol)
	switch target {
	case types.QuerySupply:
		return s.GetTotalSupply().AmountOf(tokenID), nil
	case types.QueryBurn:
		return s.GetTotalBurn().AmountOf(tokenID), nil
	case types.QueryMint:
		return s.GetTotalMint().AmountOf(tokenID), nil
	default:
		return sdk.ZeroInt(), sdk.ErrInternal(fmt.Sprintf("invalid request target to query total %s", target))
	}
}

// MintCoins creates new coins from thin air and adds it to the module account.
// Panics if the name maps to a non-minter module account or if the amount is invalid.
func (k Keeper) MintSupply(ctx sdk.Context, symbol string, to sdk.AccAddress, amt types.Coins) sdk.Error {
	_, err := k.AddCoins(ctx, symbol, to, amt)
	if err != nil {
		return err
	}
	supply := k.GetSupply(ctx, symbol)
	supply = supply.Inflate(amt)
	if supply.GetTotalSupply().IsAnyNegative() {
		return types.ErrInsufficientSupply(types.DefaultCodespace, fmt.Sprintf("insufficient supply for token [%s]", symbol))
	}

	k.SetSupply(ctx, supply)
	return nil
}

// BurnCoins burns coins deletes coins from the balance of the module account.
// Panics if the name maps to a non-burner module account or if the amount is invalid.
func (k Keeper) BurnSupply(ctx sdk.Context, symbol string, from sdk.AccAddress, amt types.Coins) sdk.Error {
	_, err := k.SubtractCoins(ctx, symbol, from, amt)
	if err != nil {
		return err
	}
	supply := k.GetSupply(ctx, symbol)
	supply = supply.Deflate(amt)
	if supply.GetTotalSupply().IsAnyNegative() {
		return types.ErrInsufficientSupply(types.DefaultCodespace, fmt.Sprintf("insufficient supply for token [%s]", symbol))
	}
	k.SetSupply(ctx, supply)

	return nil
}
