package foundation

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type (
	// AuthKeeper defines the auth module interface contract needed by the
	// foundation module.
	AuthKeeper interface {
		GetModuleAccount(ctx context.Context, name string) sdk.ModuleAccountI
	}

	// BankKeeper defines the bank module interface contract needed by the
	// foundation module.
	BankKeeper interface {
		GetAllBalances(ctx context.Context, addr sdk.AccAddress) sdk.Coins

		SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
		SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	}
)
