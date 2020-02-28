package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/collection/internal/types"
)

type SupplyKeeper interface {
	GetTotalInt(ctx sdk.Context, contractID, tokenID, target string) (supply sdk.Int, err sdk.Error)
	GetSupply(ctx sdk.Context, contractID string) (supply types.Supply)
	SetSupply(ctx sdk.Context, supply types.Supply)
	MintSupply(ctx sdk.Context, contractID string, to sdk.AccAddress, amt types.Coins) sdk.Error
	BurnSupply(ctx sdk.Context, contractID string, from sdk.AccAddress, amt types.Coins) sdk.Error
}

var _ SupplyKeeper = (*Keeper)(nil)

func (k Keeper) GetSupply(ctx sdk.Context, contractID string) (supply types.Supply) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.SupplyKey(contractID))
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

func (k Keeper) GetTotalInt(ctx sdk.Context, contractID, tokenID, target string) (supply sdk.Int, err sdk.Error) {
	if _, err = k.GetToken(ctx, contractID, tokenID); err != nil {
		return sdk.NewInt(0), err
	}

	s := k.GetSupply(ctx, contractID)
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
func (k Keeper) MintSupply(ctx sdk.Context, contractID string, to sdk.AccAddress, amt types.Coins) sdk.Error {
	_, err := k.AddCoins(ctx, contractID, to, amt)
	if err != nil {
		return err
	}
	supply := k.GetSupply(ctx, contractID)
	supply = supply.Inflate(amt)
	if supply.GetTotalSupply().IsAnyNegative() {
		return types.ErrInsufficientSupply(types.DefaultCodespace, fmt.Sprintf("insufficient supply for token [%s]", contractID))
	}

	k.SetSupply(ctx, supply)
	return nil
}

// BurnCoins burns coins deletes coins from the balance of the module account.
// Panics if the name maps to a non-burner module account or if the amount is invalid.
func (k Keeper) BurnSupply(ctx sdk.Context, contractID string, from sdk.AccAddress, amt types.Coins) sdk.Error {
	_, err := k.SubtractCoins(ctx, contractID, from, amt)
	if err != nil {
		return err
	}
	supply := k.GetSupply(ctx, contractID)
	supply = supply.Deflate(amt)
	if supply.GetTotalSupply().IsAnyNegative() {
		return types.ErrInsufficientSupply(types.DefaultCodespace, fmt.Sprintf("insufficient supply for token [%s]", contractID))
	}
	k.SetSupply(ctx, supply)

	return nil
}
