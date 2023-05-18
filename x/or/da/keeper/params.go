package keeper

import (
	"github.com/Finschia/finschia-sdk/x/or/da/types"

	sdk "github.com/Finschia/finschia-sdk/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) (types.Params, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte{types.ParamsKey})
	var params types.Params
	if bz == nil {
		return params, nil
	}
	err := k.cdc.Unmarshal(bz, &params)
	return params, err
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

func (k Keeper) CCBatchMaxBytes(ctx sdk.Context) (uint64, error) {
	params, err := k.GetParams(ctx)
	if err != nil {
		return 0, err
	}
	return params.CCBatchMaxBytes, nil
}

func (k Keeper) SCCBatchMaxBytes(ctx sdk.Context) (uint64, error) {
	params, err := k.GetParams(ctx)
	if err != nil {
		return 0, err
	}
	return params.SCCBatchMaxBytes, nil
}
