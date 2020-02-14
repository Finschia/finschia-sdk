package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// combine multiple safety box hooks, all hook functions are run in array sequence
type MultiSafetyBoxHooks []SafetyBoxHooks

func NewMultiSafetyBoxHooks(hooks ...SafetyBoxHooks) MultiSafetyBoxHooks {
	return hooks
}

func (h MultiSafetyBoxHooks) AfterSafetyBoxCreated(ctx sdk.Context, sbAddress sdk.AccAddress) {
	for _, hook := range h {
		hook.AfterSafetyBoxCreated(ctx, sbAddress)
	}
}
