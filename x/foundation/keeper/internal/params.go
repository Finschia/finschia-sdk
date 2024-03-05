package internal

import (
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/Finschia/finschia-sdk/x/foundation"
)

func (k Keeper) GetParams(ctx sdk.Context) foundation.Params {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(paramsKey)
	if err != nil {
		panic(err)
	}

	var params foundation.Params
	k.cdc.MustUnmarshal(bz, &params)

	return params
}

func (k Keeper) UpdateParams(ctx sdk.Context, params foundation.Params) error {
	// not allowed to set the tax, if it has been already disabled
	if k.GetFoundationTax(ctx).IsZero() && !params.FoundationTax.IsZero() {
		return sdkerrors.ErrInvalidRequest.Wrap("foundation tax has been already disabled")
	}

	k.SetParams(ctx, params)

	return nil
}

func (k Keeper) SetParams(ctx sdk.Context, params foundation.Params) {
	bz := k.cdc.MustMarshal(&params)

	store := k.storeService.OpenKVStore(ctx)
	if err := store.Set(paramsKey, bz); err != nil {
		panic(err)
	}
}

// aliases
func (k Keeper) GetFoundationTax(ctx sdk.Context) math.LegacyDec {
	params := k.GetParams(ctx)

	return params.FoundationTax
}

func (k Keeper) IsCensoredMessage(ctx sdk.Context, msgTypeURL string) bool {
	_, err := k.GetCensorship(ctx, msgTypeURL)
	return err == nil
}
