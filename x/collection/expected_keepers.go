package collection

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type (
	// ClassKeeper defines the contract needed to be fulfilled for class dependencies.
	ClassKeeper interface {
		NewID(ctx sdk.Context) string
		HasID(ctx sdk.Context, id string) bool
	}
)
