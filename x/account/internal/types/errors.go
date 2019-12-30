package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	CodeAccountAlreadyExist sdk.CodeType = 1
)

func ErrAccountAlreadyExist(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeAccountAlreadyExist, "Target account already exists")
}
