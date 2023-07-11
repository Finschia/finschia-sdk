package types

// DONTCOVER

import (
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
)

// x/or/settlement module sentinel errors
var (
	ErrChallengeNotExist    = sdkerrors.Register(ModuleName, 1, "challenge does not exist")
	ErrInvalidRollupName    = sdkerrors.Register(ModuleName, 2, "invalid rollup name")
	ErrInvalidL2BlockHeight = sdkerrors.Register(ModuleName, 3, "invalid L2 block height")
	ErrInvalidStepCount     = sdkerrors.Register(ModuleName, 4, "invalid step count")
	ErrChallengeExists      = sdkerrors.Register(ModuleName, 5, "challenge already exist")
)
