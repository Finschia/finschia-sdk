package stakingplus

import (
	sdk "github.com/line/lbm-sdk/types"
)

// FoundationKeeper defines the expected foundation keeper
type FoundationKeeper interface {
	GetEnabled(ctx sdk.Context) bool
	Accept(ctx sdk.Context, grantee sdk.AccAddress, msg sdk.Msg) error
}
