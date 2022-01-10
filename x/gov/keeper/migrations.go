package keeper

import (
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/gov/migrations/sample"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	keeper Keeper
}

// NewMigrator returns a new Migrator.
func NewMigrator(keeper Keeper) Migrator {
	return Migrator{keeper: keeper}
}

// MigrateSample migrates from version 1 to a sample version for migration testing.
func (m Migrator) MigrateSample(ctx sdk.Context) error {
	return sample.MigrateStore(ctx, m.keeper.storeKey, m.keeper.cdc)
}
