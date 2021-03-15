package types

import (
	sdk "github.com/line/lbm-sdk/v2/types"
)

// BankKeeper defines the contract needed for supply related APIs (noalias)
type BankKeeper interface {
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
}
