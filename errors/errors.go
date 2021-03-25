package errors

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	CodespaceLink = "link"
)

var (
	ErrError             = sdkerrors.Register(CodespaceLink, 1, "error")
	ErrInvalidPermission = sdkerrors.Register(CodespaceLink, 2, "invalid permission")
	ErrInvalidDenom      = sdkerrors.Register(CodespaceLink, 3, "invalid denom")
)
