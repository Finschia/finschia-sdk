package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Bank errors reserve 100 ~ 199.
const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	CodeSendDisabled            sdk.CodeType = 101
	CodeInvalidInputsOutputs    sdk.CodeType = 102
	CodeInvalidRequestGetsLimit sdk.CodeType = 103
)

// ErrNoInputs is an error
func ErrNoInputs(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInputsOutputs, "no inputs to send transaction")
}

// ErrNoOutputs is an error
func ErrNoOutputs(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInputsOutputs, "no outputs to send transaction")
}

// ErrInputOutputMismatch is an error
func ErrInputOutputMismatch(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInputsOutputs, "sum inputs != sum outputs")
}

// ErrSendDisabled is an error
func ErrSendDisabled(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeSendDisabled, "send transactions are currently disabled")
}

// ErrRequestGetsLimit
func ErrRequestGetsLimit(codespace sdk.CodespaceType, limit int) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidRequestGetsLimit, fmt.Sprintf("the gets limit should be %d", limit))
}
