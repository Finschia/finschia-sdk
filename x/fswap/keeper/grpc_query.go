package keeper

import (
	"github.com/Finschia/finschia-sdk/x/fswap/types"
)

var _ types.QueryServer = Keeper{}
