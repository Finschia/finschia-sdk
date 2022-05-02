package types

import (
	sdk "github.com/line/lbm-sdk/types"

	"github.com/line/lbm-sdk/x/consortium"
)

// ConsortiumKeeper defines the expected consortium keeper
type ConsortiumKeeper interface {
	GetEnabled(ctx sdk.Context) bool
	GetValidatorAuth(ctx sdk.Context, valAddr sdk.ValAddress) (*consortium.ValidatorAuth, error)
}
