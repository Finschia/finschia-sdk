package types

import (
	sdkerrors "github.com/line/lbm-sdk/v2/types/errors"
)

// IBC port sentinel errors
var (
	ErrPortExists   = sdkerrors.Register(SubModuleName, 2, "port is already binded")
	ErrPortNotFound = sdkerrors.Register(SubModuleName, 3, "port not found")
	ErrInvalidPort  = sdkerrors.Register(SubModuleName, 4, "invalid port")
	ErrInvalidRoute = sdkerrors.Register(SubModuleName, 5, "route not found")
)
