package types

import sdkerrors "github.com/Finschia/finschia-sdk/types/errors"

var (
	ErrUnknownProposal  = sdkerrors.Register(ModuleName, 2, "unknown proposal")
	ErrInactiveProposal = sdkerrors.Register(ModuleName, 3, "inactive proposal")
	ErrUnknownVote      = sdkerrors.Register(ModuleName, 4, "unknown vote")
)
