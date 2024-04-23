package types

// DONTCOVER

import (
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
)

// x/fswap module sentinel errors
var (
	ErrSwappedNotFound     = sdkerrors.Register(ModuleName, 1100, "swapped does not exist")
	ErrTotalSupplyNotFound = sdkerrors.Register(ModuleName, 1101, "swappable new coin amount does not exist")
)
