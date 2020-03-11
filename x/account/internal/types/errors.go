package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrAccountAlreadyExist = sdkerrors.Register(ModuleName, 1, "Target account already exists")
)
