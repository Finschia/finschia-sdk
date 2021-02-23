package types

import (
	sdkerrors "github.com/line/lbm-sdk/types/errors"
)

var (
	ErrAccountAlreadyExist = sdkerrors.Register(ModuleName, 1, "Target account already exists")
)
