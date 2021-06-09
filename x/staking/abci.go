package staking

import (
	"time"

	abci "github.com/line/ostracon/abci/types"

	"github.com/line/lfb-sdk/telemetry"
	sdk "github.com/line/lfb-sdk/types"
	"github.com/line/lfb-sdk/x/staking/keeper"
	"github.com/line/lfb-sdk/x/staking/types"
)

// BeginBlocker will persist the current header and validator set as a historical entry
// and prune the oldest entry based on the HistoricalEntries parameter
func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

	k.TrackHistoricalInfo(ctx)
}

// Called every block, update validator set
func EndBlocker(ctx sdk.Context, k keeper.Keeper) []abci.ValidatorUpdate {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	return k.BlockValidatorUpdates(ctx)
}
