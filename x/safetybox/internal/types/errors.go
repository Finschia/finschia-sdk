package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	// Validations
	CodeSafetyBoxIdExist               sdk.CodeType = 100
	CodeSafetyBoxIdNotExist            sdk.CodeType = 101
	CodeSafetyBoxAccountExist          sdk.CodeType = 102
	CodeSafetyBoxInvalidAction         sdk.CodeType = 103
	CodeSafetyBoxInvalidMsgType        sdk.CodeType = 104
	CodeSafetyBoxInvalidRole           sdk.CodeType = 105
	CodeSafetyBoxIssuerAddressRequired sdk.CodeType = 106
	CodeSafetyBoxCoinsRequired         sdk.CodeType = 107
	CodeSafetyBoxDenomRequired         sdk.CodeType = 108
	CodeSafetyBoxIncorrectDenom        sdk.CodeType = 109

	// Permissions
	CodeSafetyBoxPermissionWhitelist   sdk.CodeType = 201
	CodeSafetyBoxPermissionAllocate    sdk.CodeType = 202
	CodeSafetyBoxPermissionRecall      sdk.CodeType = 203
	CodeSafetyBoxPermissionIssue       sdk.CodeType = 204
	CodeSafetyBoxPermissionReturn      sdk.CodeType = 205
	CodeSafetyBoxSelfPermission        sdk.CodeType = 206
	CodeSafetyBoxHasPermissionAlready  sdk.CodeType = 207
	CodeSafetyBoxHasOtherPermission    sdk.CodeType = 208
	CodeSafetyBoxDoesNotHavePermission sdk.CodeType = 209

	// Restrictions
	CodeSafetyBoxReturnMoreThanIssued    sdk.CodeType = 300
	CodeSafetyBoxRecallMoreThanAllocated sdk.CodeType = 301
	CodeSafetyBoxOnlyOneDenomAllowed     sdk.CodeType = 302
)

func ErrSafetyBoxIdRequired(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeSafetyBoxIdExist, "Safety box ID is required")
}

func ErrSafetyBoxOwnerRequired(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeSafetyBoxIdExist, "Safety box owner is required")
}

func ErrSafetyBoxIdExist(codespace sdk.CodespaceType, safetyBoxId string) sdk.Error {
	return sdk.NewError(codespace, CodeSafetyBoxIdExist, "Safety box with the ID (%s) exists", safetyBoxId)
}

func ErrSafetyBoxNotExist(codespace sdk.CodespaceType, safetyBoxId string) sdk.Error {
	return sdk.NewError(codespace, CodeSafetyBoxIdNotExist, "Safety box (ID: %s) does not exist", safetyBoxId)
}

func ErrSafetyBoxPermissionWhitelist(codespace sdk.CodespaceType, account string) sdk.Error {
	return sdk.NewError(codespace, CodeSafetyBoxPermissionWhitelist, "The account (%s) does not have a permission to grant/revoke", account)
}

func ErrSafetyBoxPermissionAllocate(codespace sdk.CodespaceType, account string) sdk.Error {
	return sdk.NewError(codespace, CodeSafetyBoxPermissionAllocate, "The account (%s) does not have a permission to allocate", account)
}

func ErrSafetyBoxPermissionRecall(codespace sdk.CodespaceType, account string) sdk.Error {
	return sdk.NewError(codespace, CodeSafetyBoxPermissionRecall, "The account (%s) does not have a permission to recall", account)
}

func ErrSafetyBoxPermissionIssue(codespace sdk.CodespaceType, account string) sdk.Error {
	return sdk.NewError(codespace, CodeSafetyBoxPermissionIssue, "The account (%s) does not have a permission to issue or get issued", account)
}

func ErrSafetyBoxPermissionReturn(codespace sdk.CodespaceType, account string) sdk.Error {
	return sdk.NewError(codespace, CodeSafetyBoxPermissionReturn, "The account (%s) does not have a permission to return", account)
}

func ErrSafetyBoxInvalidAction(codespace sdk.CodespaceType, action string) sdk.Error {
	return sdk.NewError(codespace, CodeSafetyBoxInvalidAction, "Invalid action: %s", action)
}

func ErrSafetyBoxInvalidMsgType(codespace sdk.CodespaceType, msgType string) sdk.Error {
	return sdk.NewError(codespace, CodeSafetyBoxInvalidMsgType, "Invalid msg type: %s", msgType)
}

func ErrSafetyBoxSelfPermission(codespace sdk.CodespaceType, account string) sdk.Error {
	return sdk.NewError(codespace, CodeSafetyBoxSelfPermission, "Can not grant/revoke a permission to itself (%s)", account)
}

func ErrSafetyBoxHasPermissionAlready(codespace sdk.CodespaceType, account string) sdk.Error {
	return sdk.NewError(codespace, CodeSafetyBoxHasPermissionAlready, "The account (%s) already has the permission", account)
}

func ErrSafetyBoxHasOtherPermission(codespace sdk.CodespaceType, account string) sdk.Error {
	return sdk.NewError(codespace, CodeSafetyBoxHasOtherPermission, "The account (%s) has other permission(s)", account)
}

func ErrSafetyBoxDoesNotHavePermission(codespace sdk.CodespaceType, account string) sdk.Error {
	return sdk.NewError(codespace, CodeSafetyBoxDoesNotHavePermission, "The account (%s) does not have the permission", account)
}

func ErrSafetyBoxAccountExist(codespace sdk.CodespaceType, safetyBoxId string) sdk.Error {
	return sdk.NewError(codespace, CodeSafetyBoxAccountExist, "The safety box id (%s) exists - please try different safety box id", safetyBoxId)
}

func ErrSafetyBoxReturnMoreThanIssued(codespace sdk.CodespaceType, has, toReturn sdk.Coins) sdk.Error {
	return sdk.NewError(codespace, CodeSafetyBoxReturnMoreThanIssued, "Can not return more than issued. Has: %v, Requested: %v", has, toReturn)
}

func ErrSafetyBoxIssuerAddressRequired(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeSafetyBoxIssuerAddressRequired, "Issuer address is required to issue")
}

func ErrSafetyBoxRecallMoreThanAllocated(codespace sdk.CodespaceType, has, toRecall sdk.Coins) sdk.Error {
	return sdk.NewError(codespace, CodeSafetyBoxRecallMoreThanAllocated, "Can not recall more than allocated. Has: %v, Requested: %v", has, toRecall)
}

func ErrSafetyBoxCoinsRequired(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeSafetyBoxCoinsRequired, "Coins required to send coins")
}

func ErrSafetyBoxInvalidRole(codespace sdk.CodespaceType, role string) sdk.Error {
	return sdk.NewError(codespace, CodeSafetyBoxInvalidRole, "Invalid role: %s", role)
}

func ErrSafetyBoxDenomRequired(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeSafetyBoxDenomRequired, "Must specify a denom to create a safety box")
}

func ErrSafetyBoxTooManyCoinDenoms(codespace sdk.CodespaceType, denoms []string) sdk.Error {
	return sdk.NewError(codespace, CodeSafetyBoxOnlyOneDenomAllowed, "Only one coin denom is allowed. Requested: %v", denoms)
}

func ErrSafetyBoxIncorrectDenom(codespace sdk.CodespaceType, expectedDenom, givenDenom string) sdk.Error {
	return sdk.NewError(codespace, CodeSafetyBoxIncorrectDenom, "The safety box doesn't accept the denom. Expected: %s, Requested: %s", expectedDenom, givenDenom)
}
