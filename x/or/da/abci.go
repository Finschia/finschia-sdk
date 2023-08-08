package da

import (
	"time"

	"github.com/Finschia/finschia-rdk/telemetry"
	sdk "github.com/Finschia/finschia-rdk/types"
	"github.com/Finschia/finschia-rdk/x/or/da/keeper"
	"github.com/Finschia/finschia-rdk/x/or/da/types"
)

func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)
	if err := k.UpdateQueueTxsStatus(ctx); err != nil {
		panic(err)
	}
}
