package foundation

import (
	sdk "github.com/line/lbm-sdk/types"
	authtypes "github.com/line/lbm-sdk/x/auth/types"
	stakingtypes "github.com/line/lbm-sdk/x/staking/types"
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

		SendCoinsFromModuleToModule(ctx sdk.Context, senderModule string, recipientModule string, amt sdk.Coins) error
		SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
		SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	}

	// StakingKeeper defines the staking module interface contract needed by the
	// foundation module.
	StakingKeeper interface {
		// iterate through validators by operator address, execute func for each validator
		IterateValidators(ctx sdk.Context, fn func(index int64, validator stakingtypes.ValidatorI) (stop bool))
	}
)
