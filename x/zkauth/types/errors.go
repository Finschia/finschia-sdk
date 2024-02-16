package types

import (
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
)

const zkAuthCodespace = ModuleName

var (
	ErrInvalidZkAuthInputs = sdkerrors.Register(zkAuthCodespace, 2, "invalid zkauth inputs")
	ErrSample              = sdkerrors.Register(zkAuthCodespace, 1100, "sample error")
)
