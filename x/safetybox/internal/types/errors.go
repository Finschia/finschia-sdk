package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrSafetyBoxIDRequired              = sdkerrors.Register(ModuleName, 1, "Safety box ID is required")
	ErrSafetyBoxOwnerRequired           = sdkerrors.Register(ModuleName, 2, "Safety box owner is required")
	ErrSafetyBoxIDExist                 = sdkerrors.Register(ModuleName, 3, "Safety box with the ID exists")
	ErrSafetyBoxNotExist                = sdkerrors.Register(ModuleName, 4, "Safety box ID does not exist")
	ErrSafetyBoxPermissionWhitelist     = sdkerrors.Register(ModuleName, 5, "The account does not have a permission to grant/revoke")
	ErrSafetyBoxPermissionAllocate      = sdkerrors.Register(ModuleName, 6, "The account does not have a permission to allocate")
	ErrSafetyBoxPermissionRecall        = sdkerrors.Register(ModuleName, 7, "The account does not have a permission to recall")
	ErrSafetyBoxPermissionIssue         = sdkerrors.Register(ModuleName, 8, "The account does not have a permission to issue or get issued")
	ErrSafetyBoxPermissionReturn        = sdkerrors.Register(ModuleName, 9, "The account does not have a permission to return")
	ErrSafetyBoxInvalidAction           = sdkerrors.Register(ModuleName, 10, "Invalid action")
	ErrSafetyBoxInvalidMsgType          = sdkerrors.Register(ModuleName, 11, "Invalid msg type")
	ErrSafetyBoxSelfPermission          = sdkerrors.Register(ModuleName, 12, "Can not grant/revoke a permission to itself")
	ErrSafetyBoxHasPermissionAlready    = sdkerrors.Register(ModuleName, 13, "The account already has the permission")
	ErrSafetyBoxHasOtherPermission      = sdkerrors.Register(ModuleName, 14, "The account has other permission(s)")
	ErrSafetyBoxDoesNotHavePermission   = sdkerrors.Register(ModuleName, 15, "The account does not have the permission")
	ErrSafetyBoxAccountExist            = sdkerrors.Register(ModuleName, 16, "The safety box id exists - please try different safety box id")
	ErrSafetyBoxReturnMoreThanIssued    = sdkerrors.Register(ModuleName, 17, "Can not return more than issued")
	ErrSafetyBoxIssuerAddressRequired   = sdkerrors.Register(ModuleName, 18, "Issuer address is required to issue")
	ErrSafetyBoxRecallMoreThanAllocated = sdkerrors.Register(ModuleName, 19, "Can not recall more than allocated")
	ErrSafetyBoxCoinsRequired           = sdkerrors.Register(ModuleName, 20, "Coins required to send coins")
	ErrSafetyBoxInvalidRole             = sdkerrors.Register(ModuleName, 21, "Invalid role")
	ErrSafetyBoxDenomRequired           = sdkerrors.Register(ModuleName, 22, "Must specify a denom to create a safety box")
	ErrSafetyBoxTooManyCoinDenoms       = sdkerrors.Register(ModuleName, 23, "Only one coin denom is allowed")
	ErrSafetyBoxIncorrectDenom          = sdkerrors.Register(ModuleName, 24, "The safety box doesn't accept the denom")
)
