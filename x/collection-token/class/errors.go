package class

import (
	errorsmod "cosmossdk.io/errors"
)

const contractCodespace = "contract"

var (
	ErrInvalidContractID = errorsmod.Register(contractCodespace, 2, "invalid contractID")
	ErrContractNotExist  = errorsmod.Register(contractCodespace, 3, "contract does not exist")
)
