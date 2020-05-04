package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
)

type UpgradableModule interface {
	GetUpgradeHandler(version string) upgrade.UpgradeHandler
}

func (app *LinkApp) setUpgradeHandlers(versions ...string) {
	for _, version := range versions {
		app.upgradeKeeper.SetUpgradeHandler(version, upgradeHandler(version))
	}
}

func upgradeHandler(version string) upgrade.UpgradeHandler {
	return func(ctx sdk.Context, plan upgrade.Plan) {
		for _, m := range UpgradableModules {
			handler := m.GetUpgradeHandler(version)
			handler(ctx, plan)
		}
	}
}
