package types

import (
	sdk "github.com/line/lbm-sdk/types"

	stakingtypes "github.com/line/lbm-sdk/x/staking/types"
)

type (
	StakingKeeper interface {
		Validator(sdk.Context, sdk.ValAddress) stakingtypes.ValidatorI
		IterateValidators(ctx sdk.Context, fn func(index int64, validator stakingtypes.ValidatorI) (stop bool))
	}
	SlashingKeeper interface {
		Tombstone(sdk.Context, sdk.ConsAddress)
	}
)

