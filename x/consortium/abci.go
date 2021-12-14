package consortium

import (
	"time"

	"github.com/line/lbm-sdk/telemetry"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/consortium/keeper"
	"github.com/line/lbm-sdk/x/consortium/types"
)

// EndBlocker called every block, update validator set.
func EndBlocker(ctx sdk.Context, keeper keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	if !keeper.GetEnabled(ctx) {
		return
	}

	// reject delegations whose operators are not allowed to create a validator.
	keeper.RejectDelegations(ctx)
}
