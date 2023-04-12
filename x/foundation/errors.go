package foundation

import (
	sdkerrors "github.com/line/lbm-sdk/types/errors"
)

// x/foundation module sentinel errors
var (
	ErrInvalidGovMintLeftCount = sdkerrors.Register(ModuleName, 2, "invalid gov mint left count")
)
