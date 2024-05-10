package keeper

import (
	sdk "github.com/Finschia/finschia-sdk/types"
	banktypes "github.com/Finschia/finschia-sdk/x/bank/types"
)

type (
	BankKeeper interface {
		GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
		SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
		SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
		IsSendEnabledCoins(ctx sdk.Context, coins ...sdk.Coin) error
		GetDenomMetaData(ctx sdk.Context, denom string) (banktypes.Metadata, bool)
		SetDenomMetaData(ctx sdk.Context, denomMetaData banktypes.Metadata)
		HasSupply(ctx sdk.Context, denom string) bool
		MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
		BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	}
)
