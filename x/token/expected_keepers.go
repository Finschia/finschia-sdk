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

	// ClassKeeper defines the contract needed to be fullfilled for class dependencies.
	ClassKeeper interface {
		NewId(ctx sdk.Context) string
		HasId(ctx sdk.Context, id string) bool
		DeleteId(ctx sdk.Context, id string)
	}
)
