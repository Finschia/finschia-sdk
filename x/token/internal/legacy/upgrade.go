package legacy

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
)

func UpgradeHandler(version string) upgrade.UpgradeHandler {
	// XXX: return handler for the migration version
	return func(ctx sdk.Context, plan upgrade.Plan) {

	}
}
