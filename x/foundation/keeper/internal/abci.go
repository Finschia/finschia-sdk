package internal

import (
	"time"

	"github.com/line/lbm-sdk/telemetry"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/foundation"
)

// BeginBlocker withdraws rewards from fee-collector before the distribution
// module's withdraw.
func BeginBlocker(ctx sdk.Context, k Keeper) {
	defer telemetry.ModuleMeasureSince(foundation.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

	feeCollector := k.authKeeper.GetModuleAccount(ctx, k.feeCollectorName).GetAddress()
	if err := k.CollectFoundationTax(ctx, feeCollector); err != nil {
		panic(err)
	}
}

func EndBlocker(ctx sdk.Context, k Keeper) {
	k.UpdateTallyOfVPEndProposals(ctx)
	k.PruneExpiredProposals(ctx)
}
