package types

import (
	sdk "github.com/line/lbm-sdk/v2/types"
)

// UpgradeHandler specifies the type of function that is called when an upgrade is applied
type UpgradeHandler func(ctx sdk.Context, plan Plan)
