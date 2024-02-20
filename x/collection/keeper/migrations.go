package keeper

import (
	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Finschia/finschia-sdk/x/collection/keeper/migrations/v3"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	keeper Keeper
}

// NewMigrator returns a new Migrator.
func NewMigrator(keeper Keeper) Migrator {
	return Migrator{keeper: keeper}
}

func (m Migrator) Migrate1to2(ctx sdk.Context) error {
	return nil
}

func (m Migrator) Migrate2to3(ctx sdk.Context) error {
	store := runtime.KVStoreAdapter(m.keeper.storeService.OpenKVStore(ctx))
	oldClassStore := runtime.KVStoreAdapter(runtime.NewKVStoreService(storetypes.NewKVStoreKey(v3.ClassStoreKey)).OpenKVStore(ctx))
	return v3.MigrateStore(store, oldClassStore, m.keeper.cdc)
}
