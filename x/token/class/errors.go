package class

import (
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
)

const contractCodespace = "contract"

var (
	ErrInvalidContractID = sdkerrors.Register(contractCodespace, 2, "invalid contractID")
	ErrContractNotExist  = sdkerrors.Register(contractCodespace, 3, "contract does not exist")
)
