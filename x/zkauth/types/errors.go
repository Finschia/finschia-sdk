package types

import (
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
)

var (
	ErrSample = sdkerrors.Register(ModuleName, 1100, "sample error")
)
