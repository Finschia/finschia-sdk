package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	CodeSafetyBoxIdExist                 sdk.CodeType = 1
	CodeSafetyBoxIdNotExist              sdk.CodeType = 2
	CodeSafetyBoxPermissionWhitelist     sdk.CodeType = 3
	CodeSafetyBoxPermissionAllocate      sdk.CodeType = 4
	CodeSafetyBoxPermissionRecall        sdk.CodeType = 5
	CodeSafetyBoxPermissionIssue         sdk.CodeType = 6
	CodeSafetyBoxPermissionReturn        sdk.CodeType = 7
	CodeSafetyBoxInvalidAction           sdk.CodeType = 8
	CodeSafetyBoxInvalidMsgType          sdk.CodeType = 9
	CodeSafetyBoxSelfPermission          sdk.CodeType = 11
	CodeSafetyBoxHasPermissionAlready    sdk.CodeType = 12
	CodeSafetyBoxHasOtherPermission      sdk.CodeType = 13
	CodeSafetyBoxDoesNotHavePermission   sdk.CodeType = 14
	CodeSafetyBoxAccountExist            sdk.CodeType = 15
	CodeSafetyBoxReturnMoreThanIssued    sdk.CodeType = 16
	CodeSafetyBoxNeedsIssuerAddress      sdk.CodeType = 17
	CodeSafetyBoxRecallMoreThanAllocated sdk.CodeType = 18
	CodeSafetyBoxCoinsRequired           sdk.CodeType = 19
	CodeSafetyBoxInvalidRole             sdk.CodeType = 20
	CodeSafetyBoxDenomRequired           sdk.CodeType = 21
	CodeSafetyBoxOnlyOneDenomAllowed     sdk.CodeType = 22
	CodeSafetyBoxIncorrectDenom          sdk.CodeType = 23
)

func ErrSafetyBoxIdRequired(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeSafetyBoxIdExist, "SafetyBox ID is required")
}

func ErrSafetyBoxOwnerRequired(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeSafetyBoxIdExist, "SafetyBox owner is required")
}

func ErrSafetyBoxIdExist(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeSafetyBoxIdExist, "SafetyBox with the ID exists")
}

func ErrSafetyBoxNotExist(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeSafetyBoxIdNotExist, "SafetyBox does not exist")
}

func ErrSafetyBoxPermissionWhitelist(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeSafetyBoxPermissionWhitelist, "The account does not have a permission to whitelist")
}

func ErrSafetyBoxPermissionAllocate(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeSafetyBoxPermissionAllocate, "The account does not have a permission to allocate")
}

func ErrSafetyBoxPermissionRecall(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeSafetyBoxPermissionRecall, "The account does not have a permission to recall")
}

func ErrSafetyBoxPermissionIssue(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeSafetyBoxPermissionIssue, "The account does not have a permission to issue or get issued")
}

func ErrSafetyBoxPermissionReturn(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeSafetyBoxPermissionReturn, "The account does not have a permission to return")
}

func ErrSafetyBoxInvalidAction(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeSafetyBoxInvalidAction, "Invalid action")
}

func ErrSafetyBoxInvalidMsgType(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeSafetyBoxInvalidMsgType, "Invalid msg type")
}

func ErrSafetyBoxSelfPermission(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeSafetyBoxSelfPermission, "Can not grant/revoke a permission to itself")
}

func ErrSafetyBoxHasPermissionAlready(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeSafetyBoxHasPermissionAlready, "The account already has the permission")
}

func ErrSafetyBoxHasOtherPermission(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeSafetyBoxHasOtherPermission, "The account has other permission(s)")
}

func ErrSafetyBoxDoesNotHavePermission(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeSafetyBoxDoesNotHavePermission, "The account does not have the permission")
}

func ErrSafetyBoxAccountExist(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeSafetyBoxAccountExist, "The safety id exists - please try different safety box id")
}

func ErrSafetyBoxReturnMoreThanIssued(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeSafetyBoxReturnMoreThanIssued, "Can not return more than issued")
}

func ErrSafetyBoxNeedsIssuerAddress(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeSafetyBoxNeedsIssuerAddress, "Issuer address is required to issue")
}

func ErrSafetyBoxRecallMoreThanAllocated(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeSafetyBoxRecallMoreThanAllocated, "Can not recall more than allocated")
}

func ErrSafetyBoxCoinsRequired(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeSafetyBoxCoinsRequired, "Coins required to send coins")
}

func ErrSafetyBoxInvalidRole(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeSafetyBoxInvalidRole, "Invalid role")
}

func ErrSafetyBoxDenomRequired(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeSafetyBoxDenomRequired, "Must specify a denom to create a safety box")
}

func ErrSafetyBoxTooManyCoinDenoms(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeSafetyBoxOnlyOneDenomAllowed, "Only one coin denom is allowed")
}

func ErrSafetyBoxIncorrectDenom(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeSafetyBoxIncorrectDenom, "The safety box doesn't accept the denom")
}
