package token

import (
	sdk "github.com/line/lbm-sdk/types"
	authtypes "github.com/line/lbm-sdk/x/auth/types"
)

type (
	// AccountKeeper defines the contract required for account APIs.
	AccountKeeper interface {
		GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
		SetAccount(ctx sdk.Context, account authtypes.AccountI)

		NewAccountWithAddress(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI
	}

	// ClassKeeper defines the contract needed to be fulfilled for class dependencies.
	ClassKeeper interface {
		NewID(ctx sdk.Context) string
		HasID(ctx sdk.Context, id string) bool
		DeleteID(ctx sdk.Context, id string)

		InitGenesis(ctx sdk.Context, data *ClassGenesisState)
		ExportGenesis(ctx sdk.Context) *ClassGenesisState
	}
)
