package v2

import (
	"fmt"

	"github.com/Finschia/finschia-sdk/codec"
	storetypes "github.com/Finschia/finschia-sdk/store/types"
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/foundation"
)

// MigrateStore performs in-place store migrations from v1 to v2.
func MigrateStore(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.BinaryCodec, subspace Subspace) error {
	store := ctx.KVStore(storeKey)

	// migrate params
	if err := migrateParams(ctx, store, cdc, subspace); err != nil {
		return err
	}

	return nil
}

func migrateParams(ctx sdk.Context, store storetypes.KVStore, cdc codec.BinaryCodec, subspace Subspace) error {
	bz := store.Get(ParamsKey)
	if bz == nil {
		return fmt.Errorf("params not found")
	}
	store.Delete(ParamsKey)

	var params foundation.Params
	if err := cdc.Unmarshal(bz, &params); err != nil {
		return err
	}

	subspace.SetParamSet(ctx, &params)

	return nil
}
