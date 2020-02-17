package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/bank/internal/types"
)

// When a safety box is created, blacklist its address
func (keeper Keeper) AfterSafetyBoxCreated(ctx sdk.Context, sbAddress sdk.AccAddress) {
	keeper.BlacklistAccountAction(ctx, sbAddress, types.ActionTransferTo)
}

//_________________________________________________________________________________________

// SafetyBoxHooks event hooks
type SafetyBoxHooks interface {
	AfterSafetyBoxCreated(ctx sdk.Context, sbAddress sdk.AccAddress) // Must be called when a safety box is created
}

// Hooks wrapper struct for safety box keeper
type Hooks struct {
	keeper Keeper
}

// Return the wrapper struct
func (keeper Keeper) Hooks() *Hooks {
	return &Hooks{keeper}
}

func (h Hooks) AfterSafetyBoxCreated(ctx sdk.Context, sbAddress sdk.AccAddress) {
	h.keeper.AfterSafetyBoxCreated(ctx, sbAddress)
}
