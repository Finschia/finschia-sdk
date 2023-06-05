package types

import (
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
)

var (
	ErrInvalidCompressedData = sdkerrors.Register(ModuleName, 1100, "this data cannot be decompressed.")
)
