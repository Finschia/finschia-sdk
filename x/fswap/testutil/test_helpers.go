package testutil

import (
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/fswap/types"
)

var MockOldCoin = sdk.NewCoin(types.DefaultConfig().OldCoinDenom, sdk.NewInt(1000))
