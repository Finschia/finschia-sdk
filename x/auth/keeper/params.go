package keeper

import (
	sdk "github.com/line/lbm-sdk/v2/types"
	"github.com/line/lbm-sdk/v2/x/auth/types"
)

var param types.Params
var paramCached bool = false

// SetParams sets the auth module's parameters.
func (ak AccountKeeper) SetParams(ctx sdk.Context, params types.Params) {
	param = params
	ak.paramSubspace.SetParamSet(ctx, &params)
	paramCached = true
}

// GetParams gets the auth module's parameters.
func (ak AccountKeeper) GetParams(ctx sdk.Context) (params types.Params) {
	if !paramCached {
		ak.paramSubspace.GetParamSet(ctx, &params)
		param = params
		paramCached = true
	}
	params = param
	return
}
