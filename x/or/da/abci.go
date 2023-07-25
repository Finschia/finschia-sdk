package da

import (
	"time"

	"github.com/Finschia/finschia-sdk/telemetry"
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/or/da/keeper"
	"github.com/Finschia/finschia-sdk/x/or/da/types"
)

func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)
	if err := k.UpdateQueueTxsStatus(ctx); err != nil {
		panic(err)
	}
}
