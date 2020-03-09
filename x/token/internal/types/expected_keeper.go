package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	auth "github.com/cosmos/cosmos-sdk/x/auth/exported"
	iam "github.com/line/link/x/iam/exported"
)

type IamKeeper interface {
	GetPermissions(sdk.Context, sdk.AccAddress) []iam.PermissionI
	InheritPermission(sdk.Context, sdk.AccAddress, sdk.AccAddress)
	GrantPermission(sdk.Context, sdk.AccAddress, iam.PermissionI)
	RevokePermission(sdk.Context, sdk.AccAddress, iam.PermissionI)
	HasPermission(sdk.Context, sdk.AccAddress, iam.PermissionI) bool
	WithPrefix(string) iam.IamKeeper
}

type AccountKeeper interface {
	NewAccountWithAddress(ctx sdk.Context, addr sdk.AccAddress) auth.Account
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) auth.Account
	SetAccount(ctx sdk.Context, acc auth.Account)
}
