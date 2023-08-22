package types

import (
	sdkerrors "github.com/Finschia/finschia-rdk/types/errors"
)

const StoreCodespace = "store"

var ErrInvalidProof = sdkerrors.Register(StoreCodespace, 2, "invalid proof")
