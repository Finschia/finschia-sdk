package types

import sdkerrors "github.com/Finschia/finschia-sdk/types/errors"

var (
	ErrUnknownProposal = sdkerrors.Register(ModuleName, 2, "unknown proposal")
	ErrUnknownVote     = sdkerrors.Register(ModuleName, 3, "unknown vote")
	ErrInactiveBridge  = sdkerrors.Register(ModuleName, 4, "the bridge has halted")
)
