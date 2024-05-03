package types

// DONTCOVER

import (
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
)

// x/fswap module sentinel errors
var (
	ErrInvalidState                = sdkerrors.Register(ModuleName, 2, "swap module invalid state")
	ErrCanNotHaveMoreSwap          = sdkerrors.Register(ModuleName, 3, "no more swap allowed")
	ErrSwappedNotFound             = sdkerrors.Register(ModuleName, 4, "swapped does not exist")
	ErrExceedSwappableToCoinAmount = sdkerrors.Register(ModuleName, 5, "exceed swappable to-coin amount")
)
