package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/Finschia/finschia-sdk/x/collection"
)

func (k Keeper) GetParams(ctx sdk.Context) collection.Params {
	store := k.storeService.OpenKVStore(ctx)
	key := paramsKey
	bz, err := store.Get(key)
	if err != nil {
		panic(err)
	}
	if bz == nil {
		panic(sdkerrors.ErrNotFound.Wrap("params does not exist"))
	}

	var params collection.Params
	k.cdc.MustUnmarshal(bz, &params)

	return params
}

func (k Keeper) SetParams(ctx sdk.Context, params collection.Params) {
	store := k.storeService.OpenKVStore(ctx)
	key := paramsKey

	bz, err := params.Marshal()
	if err != nil {
		panic(err)
	}

	if err := store.Set(key, bz); err != nil {
		panic(err)
	}
}
