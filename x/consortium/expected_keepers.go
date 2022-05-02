package consortium

import (
	sdk "github.com/line/lbm-sdk/types"
	stakingtypes "github.com/line/lbm-sdk/x/staking/types"
)

type (
	// StakingKeeper defines the staking module interface contract needed by the
	// consortium module.
	StakingKeeper interface {
		// iterate through validators by operator address, execute func for each validator
		IterateValidators(ctx sdk.Context, fn func(index int64, validator stakingtypes.ValidatorI) (stop bool))
	}
)
