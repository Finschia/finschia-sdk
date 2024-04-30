package types

// DONTCOVER

import (
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
)

// x/fswap module sentinel errors
var (
	ErrSwapNotInitialized           = sdkerrors.Register(ModuleName, 1100, "swap not initialized")
	ErrSwapCanNotBeInitializedTwice = sdkerrors.Register(ModuleName, 1101, "swap cannot be initialized twice")
	ErrSwappedNotFound              = sdkerrors.Register(ModuleName, 1102, "swapped does not exist")
	ErrExceedSwappableToCoinAmount  = sdkerrors.Register(ModuleName, 1103, "exceed swappable to-coin amount")
)
