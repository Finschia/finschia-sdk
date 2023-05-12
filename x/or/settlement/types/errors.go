package types

// DONTCOVER

import (
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
)

// x/or/settlement module sentinel errors
var (
	ErrChallengeNotExist = sdkerrors.Register(ModuleName, 1, "challenge does not exist")
)
