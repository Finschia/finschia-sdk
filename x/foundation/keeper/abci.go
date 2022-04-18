package keeper

import (
	"time"

	abci "github.com/line/ostracon/abci/types"

	"github.com/line/lbm-sdk/telemetry"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/foundation"
)

// BeginBlocker withdraws rewards from fee-collector before the distribution
// module's withdraw.
func BeginBlocker(ctx sdk.Context, _ abci.RequestBeginBlock, k Keeper) {
	defer telemetry.ModuleMeasureSince(foundation.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

	if err := k.collectFoundationTax(ctx); err != nil {
		panic(err)
	}
}
