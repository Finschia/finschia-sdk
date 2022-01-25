package types

import (
	sdk "github.com/line/lbm-sdk/types"

	consortiumtypes "github.com/line/lbm-sdk/x/consortium/types"
)

// ConsortiumKeeper defines the expected consortium keeper
type ConsortiumKeeper interface {
	GetEnabled(ctx sdk.Context) bool
	GetValidatorAuth(ctx sdk.Context, valAddr sdk.ValAddress) (*consortiumtypes.ValidatorAuth, error)
}
