package types

import (
	sdkerrors "github.com/line/lbm-sdk/types/errors"
)

var (
	ErrTokenExist               = sdkerrors.Register(ModuleName, 1, "token already exists")
	ErrTokenNotExist            = sdkerrors.Register(ModuleName, 2, "token does not exist")
	ErrTokenNotMintable         = sdkerrors.Register(ModuleName, 3, "token is not mintable")
	ErrInvalidTokenName         = sdkerrors.Register(ModuleName, 4, "token name should not be empty")
	ErrInvalidTokenDecimals     = sdkerrors.Register(ModuleName, 5, "token decimals should be within the range in 0 ~ 18")
	ErrInvalidAmount            = sdkerrors.Register(ModuleName, 6, "invalid token amount")
	ErrInvalidImageURILength    = sdkerrors.Register(ModuleName, 7, "invalid token uri length")
	ErrInvalidNameLength        = sdkerrors.Register(ModuleName, 8, "invalid name length")
	ErrInvalidTokenSymbol       = sdkerrors.Register(ModuleName, 9, "invalid token symbol")
	ErrTokenNoPermission        = sdkerrors.Register(ModuleName, 10, "account does not have the permission")
	ErrAccountExist             = sdkerrors.Register(ModuleName, 11, "account already exists")
	ErrAccountNotExist          = sdkerrors.Register(ModuleName, 12, "account does not exists")
	ErrInsufficientBalance      = sdkerrors.Register(ModuleName, 13, "insufficient balance")
	ErrSupplyExist              = sdkerrors.Register(ModuleName, 14, "supply for token already exists")
	ErrInsufficientSupply       = sdkerrors.Register(ModuleName, 15, "insufficient supply")
	ErrInvalidChangesFieldCount = sdkerrors.Register(ModuleName, 16, "invalid count of field changes")
	ErrEmptyChanges             = sdkerrors.Register(ModuleName, 17, "changes is empty")
	ErrInvalidChangesField      = sdkerrors.Register(ModuleName, 18, "invalid field of changes")
	ErrDuplicateChangesField    = sdkerrors.Register(ModuleName, 19, "invalid field of changes")
	ErrInvalidMetaLength        = sdkerrors.Register(ModuleName, 20, "invalid meta length")
	ErrSupplyOverflow           = sdkerrors.Register(ModuleName, 21, "supply for token reached maximum")
	ErrApproverProxySame        = sdkerrors.Register(ModuleName, 22, "approver is same with proxy")
	ErrTokenNotApproved         = sdkerrors.Register(ModuleName, 23, "proxy is not approved on the token")
	ErrTokenAlreadyApproved     = sdkerrors.Register(ModuleName, 24, "proxy is already approved on the token")
)

func WrapIfOverflowPanic(r interface{}) error {
	if isOverflowPanic(r) {
		return ErrSupplyOverflow
	}
	// unknown panic, bubble up :(
	panic(r)
}

func isOverflowPanic(r interface{}) bool {
	return r == "Int overflow" || r == "negative coin amount"
}
