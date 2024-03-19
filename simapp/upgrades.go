package simapp

import (
	"context"
	"fmt"

	storetypes "cosmossdk.io/store/types"
	circuittypes "cosmossdk.io/x/circuit/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"

	"github.com/cosmos/cosmos-sdk/types/module"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/Finschia/finschia-sdk/simapp/internal"
)

// UpgradeName defines the on-chain upgrade name for the sample SimApp upgrade
// from v047 to v050.
//
// NOTE: This upgrade defines a reference implementation of what an upgrade
// could look like when an application is migrating from Cosmos SDK version
// v0.47.x to v0.50.x.
const UpgradeName = "v047-to-v050"

func (app SimApp) RegisterUpgradeHandlers() {
	app.UpgradeKeeper.SetUpgradeHandler(
		UpgradeName,
		func(ctx context.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
			app.deprecateBankPlusFromSimapp(ctx)
			return app.ModuleManager.RunMigrations(ctx, app.Configurator(), fromVM)
		},
	)

	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(err)
	}

	if upgradeInfo.Name == UpgradeName && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		storeUpgrades := storetypes.StoreUpgrades{
			Added: []string{
				circuittypes.ModuleName,
			},
		}

		// configure store loader that checks if version == upgradeHeight and applies store upgrades
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	}
}

// deprecateBankPlusFromSimapp remove all the states of x/bankplus module for deprecation
func (app SimApp) deprecateBankPlusFromSimapp(ctx context.Context) {
	for _, key := range app.kvStoreKeys() {
		if key.Name() == banktypes.StoreKey {
			err := internal.DeprecateBankPlus(ctx, key)
			if err != nil {
				panic(fmt.Errorf("failed to deprecate x/bankplus: %w", err))
			}
		}
	}
}
