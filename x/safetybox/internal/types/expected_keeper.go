package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/token"
)

type BankKeeper interface {
	GetBalance(ctx sdk.Context, contractID string, addr sdk.AccAddress) sdk.Int
	SetBalance(ctx sdk.Context, contractID string, addr sdk.AccAddress, amt sdk.Int) error
	HasBalance(ctx sdk.Context, contractID string, addr sdk.AccAddress, amt sdk.Int) bool

	SubtractBalance(ctx sdk.Context, contractID string, addr sdk.AccAddress, amt sdk.Int) (sdk.Int, error)
	AddBalance(ctx sdk.Context, contractID string, addr sdk.AccAddress, amt sdk.Int) (sdk.Int, error)
}

type AccountKeeper interface {
	NewAccountWithAddress(ctx sdk.Context, contractID string, addr sdk.AccAddress) (acc token.Account, err error)
	GetAccount(ctx sdk.Context, contractID string, addr sdk.AccAddress) (acc token.Account, err error)
	SetAccount(ctx sdk.Context, acc token.Account) error
}

type TokenKeeper interface {
	BankKeeper
	AccountKeeper

	GetToken(ctx sdk.Context, contractID string) (token.Token, error)
}

type SafetyBoxHooks interface {
	AfterSafetyBoxCreated(ctx sdk.Context, sbAddress sdk.AccAddress) // Must be called when a safety box is created
}
