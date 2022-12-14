package class

import (
	sdkerrors "github.com/line/lbm-sdk/types/errors"
)

const contractCodespace = "contract"

var (
	ErrInvalidContractID = sdkerrors.Register(contractCodespace, 2, "invalid contractID")
	ErrContractNotFound  = sdkerrors.Register(contractCodespace, 3, "contract does not exist")
)
