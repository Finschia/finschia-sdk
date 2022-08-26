package types

import (
	sdkerrors "github.com/line/lbm-sdk/types/errors"
)

// ICA Controller sentinel errors
var (
	ErrControllerSubModuleDisabled = sdkerrors.Register(SubModuleName, 2, "controller submodule is disabled")
)
