package slashing_test

import (
	sdk "github.com/line/lbm-sdk/types"
)

var (
	// The default power validators are initialized to have within tests
	InitTokens = sdk.TokensFromConsensusPower(200, sdk.DefaultPowerReduction)
)
