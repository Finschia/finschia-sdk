package internal

import (
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"

	"github.com/line/lbm-sdk/x/foundation"
)

func (k Keeper) GetParams(ctx sdk.Context) foundation.Params {
	store := ctx.KVStore(k.storeKey)
	key := paramsKey
	bz := store.Get(key)

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

	store := ctx.KVStore(k.storeKey)
	key := paramsKey
	store.Set(key, bz)
}

// aliases
func (k Keeper) GetFoundationTax(ctx sdk.Context) sdk.Dec {
	params := k.GetParams(ctx)

	return params.FoundationTax
}

func (k Keeper) IsCensoredMessage(ctx sdk.Context, msgTypeURL string) bool {
	_, err := k.GetCensorship(ctx, msgTypeURL)
	return err == nil
}
