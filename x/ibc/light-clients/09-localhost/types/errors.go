package types

import (
	sdkerrors "github.com/line/lbm-sdk/v2/types/errors"
)

// Localhost sentinel errors
var (
	ErrConsensusStatesNotStored = sdkerrors.Register(SubModuleName, 2, "localhost does not store consensus states")
)
