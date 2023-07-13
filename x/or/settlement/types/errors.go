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
	ErrInvalidChallengeID   = sdkerrors.Register(ModuleName, 5, "invalid challenge id")
	ErrInvalidStateHashes   = sdkerrors.Register(ModuleName, 6, "invalid state hashes")
	ErrChallengeExists      = sdkerrors.Register(ModuleName, 7, "challenge already exist")
	ErrNotSearching         = sdkerrors.Register(ModuleName, 8, "challenge is not searching")
	ErrNotResponder         = sdkerrors.Register(ModuleName, 9, "not responder")
)
