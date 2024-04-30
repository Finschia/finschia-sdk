package types

// DONTCOVER

import (
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
)

// x/fswap module sentinel errors
var (
	ErrFswapInitNotFound  = sdkerrors.Register(ModuleName, 1100, "fswap init does not exist")
	ErrSwappedNotFound    = sdkerrors.Register(ModuleName, 1101, "swapped does not exist")
	ErrExceedSwappable    = sdkerrors.Register(ModuleName, 1102, "exceed swappable coin amount")
	ErrFswapNotInitilized = sdkerrors.Register(ModuleName, 1103, "fswap not initilized")
)
