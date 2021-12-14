package types

import (
	"time"

	sdk "github.com/line/lbm-sdk/types"
	stakingtypes "github.com/line/lbm-sdk/x/staking/types"
)

type (
	StakingKeeper interface {
		GetDelegation(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (delegations stakingtypes.Delegation, found bool)
		Undelegate(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, sharesAmount sdk.Dec) (time.Time, error)
		IterateValidators(ctx sdk.Context, fn func(index int64, validator stakingtypes.ValidatorI) (stop bool))
		// jailValidator(ctx sdk.Context, validator stakingtypes.ValidatorI)
		GetValidator(ctx sdk.Context, addr sdk.ValAddress) (validator stakingtypes.Validator, found bool)
	}
)

