package types

import (
	sdkerrors "github.com/line/lbm-sdk/types/errors"
)

// x/consortium module sentinel errors
var (
	ErrInvalidProposalValidator = sdkerrors.Register(ModuleName, 1, "invalid proposal validator")
)
