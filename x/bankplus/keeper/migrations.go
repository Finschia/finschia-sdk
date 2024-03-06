package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankv2 "github.com/cosmos/cosmos-sdk/x/bank/migrations/v2"
	bankv3 "github.com/cosmos/cosmos-sdk/x/bank/migrations/v3"
)

type Migrator struct {
	keeper BaseKeeper
}

func NewMigrator(keeper BaseKeeper) *Migrator {
	return &Migrator{keeper: keeper}
}

func (m Migrator) WrappedMigrateBankplusWithBankMigrate1to2n3(ctx sdk.Context) error {
	if err := DeprecateBankPlus(ctx, m.keeper.storeService); err != nil {
		return err
	}

	if err := bankv2.MigrateStore(ctx, m.keeper.storeService, m.keeper.cdc); err != nil {
		return err
	}

	return bankv3.MigrateStore(ctx, m.keeper.storeService, m.keeper.cdc)
}
