package keeper

import (
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/fswap/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) (types.Params, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte{types.ParamsKey})
	var params types.Params
	if bz == nil {
		return types.Params{}, types.ErrParamsNotFound
	}
	err := k.cdc.Unmarshal(bz, &params)
	if err != nil {
		return types.Params{}, err
	}
	return params, nil
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
