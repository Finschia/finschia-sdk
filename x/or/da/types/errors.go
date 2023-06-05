package types

import (
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
)

var (
	ErrInvalidCompressedData = sdkerrors.Register(ModuleName, 1100, "this data cannot be decompressed.")
	ErrInvalidCCBatch        = sdkerrors.Register(ModuleName, 1101, "invalid cc batch")
	ErrCCStateNotFound       = sdkerrors.Register(ModuleName, 1102, "cc state not found")
	ErrCCRefNotFound         = sdkerrors.Register(ModuleName, 1103, "cc reference not found")
)
