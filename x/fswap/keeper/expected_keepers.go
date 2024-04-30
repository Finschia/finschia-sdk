package keeper

import (
	sdk "github.com/Finschia/finschia-sdk/types"
)

type (
	AccountKeeper interface {
		GetModuleAddress(name string) sdk.AccAddress
	}

	BankKeeper interface {
		HasBalance(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coin) bool
		GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
		SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
		SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
		MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
		BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	}
)
