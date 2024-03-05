package internal

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/Finschia/finschia-sdk/x/foundation"
	v2 "github.com/Finschia/finschia-sdk/x/foundation/keeper/internal/migrations/v2"
	v3 "github.com/Finschia/finschia-sdk/x/foundation/keeper/internal/migrations/v3"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	keeper   Keeper
	subspace paramstypes.Subspace
}

// NewMigrator returns a new Migrator.
func NewMigrator(keeper Keeper, subspace paramstypes.Subspace) Migrator {
	return Migrator{keeper: keeper}
}

func (m Migrator) Register(register func(moduleName string, fromVersion uint64, handler module.MigrationHandler) error) error {
	for fromVersion, handler := range map[uint64]module.MigrationHandler{
		1: func(ctx sdk.Context) error {
			return v2.MigrateStore(ctx, m.keeper.storeService, m.keeper.cdc, m.subspace)
		},
		2: func(ctx sdk.Context) error {
			return v3.MigrateStore(ctx, m.keeper.storeService, m.keeper.cdc, m.subspace)
		},
	} {
		if err := register(foundation.ModuleName, fromVersion, handler); err != nil {
			return err
		}
	}

	return nil
}
