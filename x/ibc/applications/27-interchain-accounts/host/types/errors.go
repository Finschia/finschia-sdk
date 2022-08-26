package types

import (
	sdkerrors "github.com/line/lbm-sdk/types/errors"
)

// ICA Host sentinel errors
var (
	ErrHostSubModuleDisabled = sdkerrors.Register(SubModuleName, 2, "host submodule is disabled")
)
