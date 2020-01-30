package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	CodeError             sdk.CodeType = 1
	CodeInvalidPermission sdk.CodeType = 2
	CodeInvalidDenom      sdk.CodeType = 3

	CodespaceLink sdk.CodespaceType = "link"
)

func ErrError(msg string) sdk.Error {
	return newErrorWithLinkCodespace(CodeError, msg)
}

func newErrorWithLinkCodespace(code sdk.CodeType, format string, args ...interface{}) sdk.Error {
	return sdk.NewError(CodespaceLink, code, format, args...)
}

func ErrInvalidPermission(msg string) sdk.Error {
	return newErrorWithLinkCodespace(CodeInvalidPermission, msg)
}
func ErrInvalidDenom(msg string) sdk.Error {
	return newErrorWithLinkCodespace(CodeInvalidDenom, msg)
}
