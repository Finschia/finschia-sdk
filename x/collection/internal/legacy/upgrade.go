package legacy

import (
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/upgrade"
)

func UpgradeHandler(version string) upgrade.UpgradeHandler {
	// XXX: return handler for the migration version
	return func(ctx sdk.Context, plan upgrade.Plan) {

	}
}
