package types

import (
	sdk "github.com/line/lbm-sdk/types"

	"github.com/line/lbm-sdk/x/foundation"
)

// FoundationKeeper defines the expected foundation keeper
type FoundationKeeper interface {
	GetEnabled(ctx sdk.Context) bool
	GetValidatorAuth(ctx sdk.Context, valAddr sdk.ValAddress) (*foundation.ValidatorAuth, error)
}
