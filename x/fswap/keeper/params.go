package keeper

import (
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/fswap/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte{types.ParamsKey})
	var params types.Params
	if bz == nil {
		panic(types.ErrParamsNotFound)
	}
	k.cdc.MustUnmarshal(bz, &params)
	return params
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := k.cdc.Marshal(&params)
	if err != nil {
		return err
	}
	store.Set([]byte{types.ParamsKey}, bz)
	return nil
}

func (k Keeper) NewCoinDenom(ctx sdk.Context) string {
	params := k.GetParams(ctx)
	return params.NewCoinDenom
}
