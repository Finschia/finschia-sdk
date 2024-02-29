package stakingplus

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// FoundationKeeper defines the expected foundation keeper
type FoundationKeeper interface {
	Accept(ctx sdk.Context, grantee sdk.AccAddress, msg sdk.Msg) error
}

// StakingMsgServer defines the expected staking keeper
type StakingMsgServer interface {
	// CreateValidator defines a method for creating a new validator.
	CreateValidator(context.Context, *staking.MsgCreateValidator) (*staking.MsgCreateValidatorResponse, error)
	// EditValidator defines a method for editing an existing validator.
	EditValidator(context.Context, *staking.MsgEditValidator) (*staking.MsgEditValidatorResponse, error)
	// Delegate defines a method for performing a delegation of coins
	// from a delegator to a validator.
	Delegate(context.Context, *staking.MsgDelegate) (*staking.MsgDelegateResponse, error)
	// BeginRedelegate defines a method for performing a redelegation
	// of coins from a delegator and source validator to a destination validator.
	BeginRedelegate(context.Context, *staking.MsgBeginRedelegate) (*staking.MsgBeginRedelegateResponse, error)
	// Undelegate defines a method for performing an undelegation from a
	// delegate and a validator.
	Undelegate(context.Context, *staking.MsgUndelegate) (*staking.MsgUndelegateResponse, error)
	// CancelUnbondingDelegation defines a method for performing canceling the unbonding delegation
	// and delegate back to previous validator.
	//
	// Since: cosmos-sdk 0.46
	CancelUnbondingDelegation(context.Context, *staking.MsgCancelUnbondingDelegation) (*staking.MsgCancelUnbondingDelegationResponse, error)
	// UpdateParams defines an operation for updating the x/staking module
	// parameters.
	// Since: cosmos-sdk 0.47
	UpdateParams(context.Context, *staking.MsgUpdateParams) (*staking.MsgUpdateParamsResponse, error)
}
