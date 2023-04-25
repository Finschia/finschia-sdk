package foundation

import (
	sdk "github.com/Finschia/finschia-sdk/types"
	authtypes "github.com/Finschia/finschia-sdk/x/auth/types"
)

type (
	// AuthKeeper defines the auth module interface contract needed by the
	// foundation module.
	AuthKeeper interface {
		GetModuleAccount(ctx sdk.Context, name string) authtypes.ModuleAccountI
	}

	// BankKeeper defines the bank module interface contract needed by the
	// foundation module.
	BankKeeper interface {
		GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins

		SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
		SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	}
)
