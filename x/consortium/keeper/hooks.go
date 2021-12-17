package keeper

import (
	sdk "github.com/line/lbm-sdk/types"
	stakingtypes "github.com/line/lbm-sdk/x/staking/types"
)

// Wrapper struct
type Hooks struct {
	k Keeper
}

var _ stakingtypes.StakingHooks = Hooks{}

// Create new consortium hooks
func (k Keeper) Hooks() Hooks { return Hooks{k} }

func (h Hooks) AfterDelegationModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	if !h.k.GetEnabled(ctx) {
		return
	}

	if auth, err := h.k.GetValidatorAuth(ctx, valAddr); err != nil || !auth.CreationAllowed {
		h.k.addPendingRejectedDelegation(ctx, delAddr, valAddr)
	}
}

func (h Hooks) BeforeDelegationRemoved(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	if !h.k.GetEnabled(ctx) {
		return
	}

	h.k.deletePendingRejectedDelegation(ctx, delAddr, valAddr)
}

func (h Hooks) AfterValidatorCreated(_ sdk.Context, _ sdk.ValAddress)                    {}
func (h Hooks) AfterValidatorRemoved(_ sdk.Context, _ sdk.ConsAddress, _ sdk.ValAddress) {}

func (h Hooks) BeforeDelegationCreated(_ sdk.Context, _ sdk.AccAddress, _ sdk.ValAddress)        {}
func (h Hooks) BeforeDelegationSharesModified(_ sdk.Context, _ sdk.AccAddress, _ sdk.ValAddress) {}
func (h Hooks) BeforeValidatorSlashed(_ sdk.Context, _ sdk.ValAddress, _ sdk.Dec)                {}

func (h Hooks) BeforeValidatorModified(_ sdk.Context, _ sdk.ValAddress)                         {}
func (h Hooks) AfterValidatorBonded(_ sdk.Context, _ sdk.ConsAddress, _ sdk.ValAddress)         {}
func (h Hooks) AfterValidatorBeginUnbonding(_ sdk.Context, _ sdk.ConsAddress, _ sdk.ValAddress) {}
