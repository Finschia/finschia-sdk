package types

import (
	sdkerrors "github.com/line/lbm-sdk/v2/types/errors"
)

const StoreCodespace = "store"

var (
	ErrInvalidProof = sdkerrors.Register(StoreCodespace, 2, "invalid proof")
)
