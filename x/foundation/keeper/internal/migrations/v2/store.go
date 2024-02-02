package v2

import (
	"fmt"

	"cosmossdk.io/core/store"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Finschia/finschia-sdk/x/foundation"
)

// MigrateStore performs in-place store migrations from v1 to v2.
func MigrateStore(ctx sdk.Context, storeService store.KVStoreService, cdc codec.BinaryCodec, subspace Subspace) error {
	store := storeService.OpenKVStore(ctx)

	// migrate params
	if err := migrateParams(ctx, store, cdc, subspace); err != nil {
		return err
	}

	return nil
}

func migrateParams(ctx sdk.Context, store store.KVStore, cdc codec.BinaryCodec, subspace Subspace) error {
	bz, err := store.Get(ParamsKey)
	if err != nil {
		return err
	}
	if bz == nil {
		return fmt.Errorf("params not found")
	}
	if err := store.Delete(ParamsKey); err != nil {
		return err
	}

	var params foundation.Params
	if err := cdc.Unmarshal(bz, &params); err != nil {
		return err
	}

	subspace.SetParamSet(ctx, &params)

	return nil
}
