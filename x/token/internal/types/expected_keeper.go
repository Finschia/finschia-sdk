package types

import (
	sdk "github.com/line/lbm-sdk/types"
	auth "github.com/line/lbm-sdk/x/auth/exported"
)

type AccountKeeper interface {
	NewAccountWithAddress(ctx sdk.Context, addr sdk.AccAddress) auth.Account
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) auth.Account
	SetAccount(ctx sdk.Context, acc auth.Account)
}
