package types

import (
	"time"

	sdk "github.com/line/lfb-sdk/types"
	stakingtypes "github.com/line/lfb-sdk/x/staking/types"
)

// StakingKeeper expected staking keeper
type StakingKeeper interface {
	GetHistoricalInfo(ctx sdk.Context, height int64) (stakingtypes.HistoricalInfo, bool)
	UnbondingTime(ctx sdk.Context) time.Duration
}
