package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type BankKeeper interface {
	GetCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	SetCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) error
	HasCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) bool

	SubtractCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, error)
	AddCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, error)
}

type SafetyBoxHooks interface {
	AfterSafetyBoxCreated(ctx sdk.Context, sbAddress sdk.AccAddress) // Must be called when a safety box is created
}
