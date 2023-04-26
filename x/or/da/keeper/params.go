package keeper

import (
	"github.com/Finschia/finschia-sdk/x/or/da/types"

	sdk "github.com/Finschia/finschia-sdk/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams(
		k.Placeholder(ctx),
	)
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}

// Placeholder returns the Placeholder param
func (k Keeper) Placeholder(ctx sdk.Context) (res string) {
	k.paramstore.Get(ctx, types.KeyPlaceholder, &res)
	return
}
