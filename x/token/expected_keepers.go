package token

import (
	sdk "github.com/Finschia/finschia-sdk/types"
)

type (
	// ClassKeeper defines the contract needed to be fulfilled for class dependencies.
	ClassKeeper interface {
		NewID(ctx sdk.Context) string
		HasID(ctx sdk.Context, id string) bool

		InitGenesis(ctx sdk.Context, data *ClassGenesisState)
		ExportGenesis(ctx sdk.Context) *ClassGenesisState
	}
)
