package exported

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type PermissionI interface {
	GetResource() string
	GetAction() string
	Equal(string, string) bool
}

type IamKeeper interface {
	GetPermissions(sdk.Context, sdk.AccAddress) []PermissionI
	InheritPermission(sdk.Context, sdk.AccAddress, sdk.AccAddress)
	GrantPermission(sdk.Context, sdk.AccAddress, PermissionI)
	RevokePermission(sdk.Context, sdk.AccAddress, PermissionI)
	HasPermission(sdk.Context, sdk.AccAddress, PermissionI) bool
	WithPrefix(string) IamKeeper
}
