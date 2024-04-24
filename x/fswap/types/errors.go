package types

// DONTCOVER

import (
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
)

// x/fswap module sentinel errors
var (
	ErrParamsNotFound                 = sdkerrors.Register(ModuleName, 1100, "params does not exist")
	ErrSwappedNotFound                = sdkerrors.Register(ModuleName, 1101, "swapped does not exist")
	ErrSwappableNewCoinAmountNotFound = sdkerrors.Register(ModuleName, 1102, "swappable new coin amount does not exist")
)
