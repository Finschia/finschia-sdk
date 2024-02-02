package internal

import (
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Finschia/finschia-sdk/x/foundation"
)

// BeginBlocker withdraws rewards from fee-collector before the distribution
// module's withdraw.
func BeginBlocker(ctx sdk.Context, k Keeper) error {
	defer telemetry.ModuleMeasureSince(foundation.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

	if err := k.CollectFoundationTax(ctx); err != nil {
		return err
	}

	return nil
}

func EndBlocker(ctx sdk.Context, k Keeper) error {
	defer telemetry.ModuleMeasureSince(foundation.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	if err := k.UpdateTallyOfVPEndProposals(ctx); err != nil {
		return err
	}

	k.PruneExpiredProposals(ctx)

	return nil
}
