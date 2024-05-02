package types

// DONTCOVER

import (
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
)

// x/fswap module sentinel errors
var (
	ErrInvalidState                = sdkerrors.Register(ModuleName, 1100, "swap module invalid state")
	ErrCanNotHaveMoreSwap          = sdkerrors.Register(ModuleName, 1101, "no more swap allowed")
	ErrSwappedNotFound             = sdkerrors.Register(ModuleName, 1102, "swapped does not exist")
	ErrExceedSwappableToCoinAmount = sdkerrors.Register(ModuleName, 1103, "exceed swappable to-coin amount")
)
