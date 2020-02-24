package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
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
