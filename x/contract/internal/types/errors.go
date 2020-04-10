package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrInvalidContractID = sdkerrors.Register(ModuleName, 1, "invalid contractID")
	ErrContractNotExist  = sdkerrors.Register(ModuleName, 2, "contract does not exist")
)
