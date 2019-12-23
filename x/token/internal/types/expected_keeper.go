package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	auth "github.com/cosmos/cosmos-sdk/x/auth/exported"
	supply "github.com/cosmos/cosmos-sdk/x/supply/exported"
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

type SupplyKeeper interface {
	GetModuleAddress(string) sdk.AccAddress
	MintCoins(sdk.Context, string, sdk.Coins) sdk.Error
	BurnCoins(sdk.Context, string, sdk.Coins) sdk.Error
	GetSupply(sdk.Context) supply.SupplyI
	SendCoinsFromModuleToAccount(sdk.Context, string, sdk.AccAddress, sdk.Coins) sdk.Error
	SendCoinsFromAccountToModule(sdk.Context, sdk.AccAddress, string, sdk.Coins) sdk.Error
}

type AccountKeeper interface {
	NewAccountWithAddress(ctx sdk.Context, addr sdk.AccAddress) auth.Account

	GetAccount(ctx sdk.Context, addr sdk.AccAddress) auth.Account
	GetAllAccounts(ctx sdk.Context) []auth.Account
	SetAccount(ctx sdk.Context, acc auth.Account)

	IterateAccounts(ctx sdk.Context, process func(auth.Account) bool)
}
