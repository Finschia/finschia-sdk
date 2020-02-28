package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	CodeInvalidContractID sdk.CodeType = 100
	CodeContractNotExist  sdk.CodeType = 101
)

func ErrInvalidContractID(codespace sdk.CodespaceType, contractID string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidContractID, "invalid contractID: %s", contractID)
}

func ErrContractNotExist(codespace sdk.CodespaceType, contractID string) sdk.Error {
	return sdk.NewError(codespace, CodeContractNotExist, "contract[%s] does not exist", contractID)
}
