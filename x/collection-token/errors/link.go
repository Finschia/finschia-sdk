package errors

import errorsmod "cosmossdk.io/errors"

const (
	linkCodespace = "link"
)

// NO additional errors allowed into this codespace
var (
	ErrInvalidPermission = errorsmod.Register(linkCodespace, 2, "invalid permission")
	ErrInvalidDenom      = errorsmod.Register(linkCodespace, 3, "invalid denom")
)
