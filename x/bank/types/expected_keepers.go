package types

import (
	sdk "github.com/line/lfb-sdk/types"
	"github.com/line/lfb-sdk/x/auth/types"
)

// AccountKeeper defines the account contract that must be fulfilled when
// creating a x/bank keeper.
type AccountKeeper interface {
	NewAccountWithAddress(ctx sdk.Context, addr sdk.AccAddress) types.AccountI

	GetAccount(ctx sdk.Context, addr sdk.AccAddress) types.AccountI
	GetAllAccounts(ctx sdk.Context) []types.AccountI
	SetAccount(ctx sdk.Context, acc types.AccountI)

	IterateAccounts(ctx sdk.Context, process func(types.AccountI) bool)

	ValidatePermissions(macc types.ModuleAccountI) error

	GetModuleAddress(moduleName string) sdk.AccAddress
	GetModuleAddressAndPermissions(moduleName string) (addr sdk.AccAddress, permissions []string)
	GetModuleAccountAndPermissions(ctx sdk.Context, moduleName string) (types.ModuleAccountI, []string)
	GetModuleAccount(ctx sdk.Context, moduleName string) types.ModuleAccountI
	SetModuleAccount(ctx sdk.Context, macc types.ModuleAccountI)
}
