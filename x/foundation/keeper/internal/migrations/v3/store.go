package v3

import (
	"cosmossdk.io/core/store"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Finschia/finschia-sdk/x/foundation"
)

// MigrateStore performs in-place store migrations from v2 to v3.
func MigrateStore(ctx sdk.Context, storeService store.KVStoreService, cdc codec.BinaryCodec, subspace Subspace) error {
	store := storeService.OpenKVStore(ctx)

	// migrate params
	if err := migrateParams(ctx, store, cdc, subspace); err != nil {
		return err
	}

	return nil
}

func migrateParams(ctx sdk.Context, store store.KVStore, cdc codec.BinaryCodec, subspace Subspace) error {
	var params foundation.Params
	subspace.GetParamSet(ctx, &params)

	bz, err := cdc.Marshal(&params)
	if err != nil {
		return err
	}

	if err := store.Set(ParamsKey, bz); err != nil {
		return err
	}

	return nil
}
