package crisis

import (
	"time"

	"github.com/line/lfb-sdk/telemetry"
	sdk "github.com/line/lfb-sdk/types"
	"github.com/line/lfb-sdk/x/crisis/keeper"
	"github.com/line/lfb-sdk/x/crisis/types"
)

// check all registered invariants
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	if k.InvCheckPeriod() == 0 || ctx.BlockHeight()%int64(k.InvCheckPeriod()) != 0 {
		// skip running the invariant check
		return
	}
	k.AssertInvariants(ctx)
}
