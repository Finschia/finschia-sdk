package internal

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/Finschia/finschia-sdk/x/foundation"
	v2 "github.com/Finschia/finschia-sdk/x/foundation/keeper/internal/migrations/v2"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	keeper Keeper
}

// NewMigrator returns a new Migrator.
func NewMigrator(keeper Keeper) Migrator {
	return Migrator{keeper: keeper}
}

func (m Migrator) Register(register func(moduleName string, fromVersion uint64, handler module.MigrationHandler) error) error {
	for fromVersion, handler := range map[uint64]module.MigrationHandler{
		1: func(ctx sdk.Context) error {
			return v2.MigrateStore(ctx, m.keeper.storeService, m.keeper.cdc, m.keeper.subspace)
		},
	} {
		if err := register(foundation.ModuleName, fromVersion, handler); err != nil {
			return err
		}
	}

	return nil
}
