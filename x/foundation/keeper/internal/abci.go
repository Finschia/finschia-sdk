package internal

import (
	"time"

	"github.com/Finschia/finschia-sdk/telemetry"
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/foundation"
)

// BeginBlocker withdraws rewards from fee-collector before the distribution
// module's withdraw.
func BeginBlocker(ctx sdk.Context, k Keeper) {
	defer telemetry.ModuleMeasureSince(foundation.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

	if err := k.CollectFoundationTax(ctx); err != nil {
		panic(err)
	}
}

func EndBlocker(ctx sdk.Context, k Keeper) {
	k.UpdateTallyOfVPEndProposals(ctx)
	k.PruneExpiredProposals(ctx)
}
