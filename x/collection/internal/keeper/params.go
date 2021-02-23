package keeper

import (
	sdk "github.com/line/lbm-sdk/types"

	"github.com/line/lbm-sdk/x/collection/internal/types"
)

func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramsSpace.SetParamSet(ctx, &params)
}

func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramsSpace.GetParamSet(ctx, &params)
	return
}
