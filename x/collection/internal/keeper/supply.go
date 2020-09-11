package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/line/link-modules/x/collection/internal/types"
)

type SupplyKeeper interface {
	GetTotalInt(ctx sdk.Context, tokenID, target string) (supply sdk.Int, err error)
	GetSupply(ctx sdk.Context) (supply types.Supply)
	SetSupply(ctx sdk.Context, supply types.Supply)
	MintSupply(ctx sdk.Context, to sdk.AccAddress, amt types.Coins) error
	BurnSupply(ctx sdk.Context, from sdk.AccAddress, amt types.Coins) error
}

var _ SupplyKeeper = (*Keeper)(nil)

func (k Keeper) GetSupply(ctx sdk.Context) (supply types.Supply) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.SupplyKey(k.getContractID(ctx)))
	if b == nil {
		panic("stored supply should not have been nil")
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &supply)
	return
}

func (k Keeper) SetSupply(ctx sdk.Context, supply types.Supply) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(supply)
	store.Set(types.SupplyKey(supply.GetContractID()), b)
}

func (k Keeper) GetTotalInt(ctx sdk.Context, tokenID, target string) (supply sdk.Int, err error) {
	if _, err = k.GetToken(ctx, tokenID); err != nil {
		return sdk.NewInt(0), err
	}

	s := k.GetSupply(ctx)
	switch target {
	case types.QuerySupply:
		return s.GetTotalSupply().AmountOf(tokenID), nil
	case types.QueryBurn:
		return s.GetTotalBurn().AmountOf(tokenID), nil
	case types.QueryMint:
		return s.GetTotalMint().AmountOf(tokenID), nil
	default:
		return sdk.ZeroInt(), sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "invalid request target to query total %s", target)
	}
}

// MintCoins creates new coins from thin air and adds it to the module account.
// Panics if the name maps to a non-minter module account or if the amount is invalid.
func (k Keeper) MintSupply(ctx sdk.Context, to sdk.AccAddress, amt types.Coins) (err error) {
	defer func() {
		// to recover from overflows
		if r := recover(); r != nil {
			err = types.WrapIfOverflowPanic(r)
		}
	}()

	_, err = k.AddCoins(ctx, to, amt)
	if err != nil {
		return err
	}
	supply := k.GetSupply(ctx)
	supply = supply.Inflate(amt)
	// supply should never be negative. Big.Int.Add will be panic if it becomes overflow

	k.SetSupply(ctx, supply)
	return nil
}

// BurnCoins burns coins deletes coins from the balance of the module account.
// Panics if the name maps to a non-burner module account or if the amount is invalid.
func (k Keeper) BurnSupply(ctx sdk.Context, from sdk.AccAddress, amt types.Coins) (err error) {
	defer func() {
		// to recover from overflows
		// however, it will return insufficient fund error instead of panicking in the case
		if r := recover(); r != nil {
			err = types.WrapIfOverflowPanic(r)
		}
	}()

	_, err = k.SubtractCoins(ctx, from, amt)
	if err != nil {
		return err
	}
	supply := k.GetSupply(ctx)
	supply = supply.Deflate(amt)
	if supply.GetTotalSupply().IsAnyNegative() {
		return sdkerrors.Wrapf(types.ErrInsufficientSupply, "insufficient supply for token [%s]", k.getContractID(ctx))
	}
	k.SetSupply(ctx, supply)

	return nil
}
