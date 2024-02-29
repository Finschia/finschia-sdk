package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankv2 "github.com/cosmos/cosmos-sdk/x/bank/migrations/v2"
)

type Migrator struct {
	keeper BaseKeeper
}

func NewMigrator(keeper BaseKeeper) *Migrator {
	return &Migrator{keeper: keeper}
}

func (m Migrator) WrappedMigrateBankplusWithBankMigrate1to2(ctx sdk.Context) error {
	err := DeprecateBankPlus(ctx, m.keeper.storeService)
	if err != nil {
		return err
	}

	return bankv2.MigrateStore(ctx, m.keeper.storeService, m.keeper.cdc)
}
