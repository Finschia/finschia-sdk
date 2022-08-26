package keeper

import (
	sdk "github.com/line/lbm-sdk/types"

	"github.com/line/lbm-sdk/x/ibc/applications/27-interchain-accounts/controller/types"
)

// IsControllerEnabled retrieves the host enabled boolean from the paramstore.
// True is returned if the controller submodule is enabled.
func (k Keeper) IsControllerEnabled(ctx sdk.Context) bool {
	var res bool
	k.paramSpace.Get(ctx, types.KeyControllerEnabled, &res)
	return res
}

// GetParams returns the total set of the host submodule parameters.
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams(k.IsControllerEnabled(ctx))
}

// SetParams sets the total set of the host submodule parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}
