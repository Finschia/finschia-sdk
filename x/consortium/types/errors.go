package types

import (
	sdkerrors "github.com/line/lbm-sdk/types/errors"
)

// x/consortium module sentinel errors
var (
	ErrInvalidParams = sdkerrors.Register(ModuleName, 1, "invalid params")

	ErrInvalidProposalValidator = sdkerrors.Register(ModuleName, 2, "invalid proposal validator")
)
