package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultCodespace is the Module Name
const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	CodeInvalidDenom              sdk.CodeType = 101
	CodeNotExistLRC3              sdk.CodeType = 102
	CodeInvalidTokenId            sdk.CodeType = 103
	CodeInvalidOwner              sdk.CodeType = 104
	CodeNotExistApprove           sdk.CodeType = 105
	CodeNotExistOperatorApprovals sdk.CodeType = 106
	CodeNotExistOperator          sdk.CodeType = 107
	CodeNotExistNFT               sdk.CodeType = 108
	CodeAlreadyExistLRC3          sdk.CodeType = 109
	CodeInvalidPermission         sdk.CodeType = 110
)

func ErrInvalidDenom(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidDenom, "denom is invalid")
}

func ErrNotExistLRC3(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeNotExistLRC3, "LRC3 does not exist")
}

func ErrNotExistNFT(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeNotExistNFT, "nft does not exist")
}

func ErrNotExistApprove(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeNotExistApprove, "approve does not exist")
}

func ErrNotExistOperatorApprovals(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeNotExistOperatorApprovals, "operator approvals does not exist")
}

func ErrNotExistOperator(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeNotExistOperator, "operator does not exist")
}

func ErrInvalidTokenId(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidTokenId, "tokenId is invalid")
}

func ErrInvalidOwner(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidOwner, "owner is invalid")
}

func ErrAlreadyExistLRC3(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeAlreadyExistLRC3, "LRC-3 already exists.")
}

func ErrInvalidPermission(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidPermission, "Permission is invalid")
}
