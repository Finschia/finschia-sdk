package keeper

import sdk "github.com/Finschia/finschia-sdk/types"

type (
	AccountKeeper interface {
		GetModuleAddress(name string) sdk.AccAddress
	}

	BankKeeper interface {
		HasBalance(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coin) bool
		GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
		SendCoins(ctx sdk.Context, fromAddr, toAddr sdk.AccAddress, amt sdk.Coins) error
		MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
		BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	}
)
