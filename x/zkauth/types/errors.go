package types

import (
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
)

const zkAuthCodespace = ModuleName

var (
	ErrInvalidZKAuthSignature = sdkerrors.Register(zkAuthCodespace, 2, "invalid ZKAuthSignature")
	ErrInvalidMessage         = sdkerrors.Register(zkAuthCodespace, 3, "invalid message")
	ErrInvalidZkAuthInputs    = sdkerrors.Register(zkAuthCodespace, 4, "invalid zkauth inputs")
	ErrSample                 = sdkerrors.Register(zkAuthCodespace, 1100, "sample error")
)
