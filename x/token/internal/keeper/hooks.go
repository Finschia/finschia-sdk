package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/bank"
)

// When a safety box is created, blacklist its address
func (k Keeper) AfterSafetyBoxCreated(ctx sdk.Context, sbAddress sdk.AccAddress) {
	k.BlacklistAccountAction(ctx, sbAddress, bank.ActionTransferTo)
}

//_________________________________________________________________________________________

// SafetyBoxHooks event hooks
type SafetyBoxHooks interface {
	AfterSafetyBoxCreated(ctx sdk.Context, sbAddress sdk.AccAddress) // Must be called when a safety box is created
}

// Hooks wrapper struct for safety box keeper
type Hooks struct {
	k Keeper
}

// Return the wrapper struct
func (k Keeper) Hooks() *Hooks {
	return &Hooks{k}
}

func (h Hooks) AfterSafetyBoxCreated(ctx sdk.Context, sbAddress sdk.AccAddress) {
	h.k.AfterSafetyBoxCreated(ctx, sbAddress)
}
