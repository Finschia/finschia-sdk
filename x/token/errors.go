package token

import (
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
)

const tokenCodespace = ModuleName

var (
	ErrTokenNotExist            = sdkerrors.Register(tokenCodespace, 2, "token does not exist")
	ErrTokenNotMintable         = sdkerrors.Register(tokenCodespace, 3, "token is not mintable") // Currently never used
	ErrInvalidTokenName         = sdkerrors.Register(tokenCodespace, 4, "token name should not be empty")
	ErrInvalidTokenDecimals     = sdkerrors.Register(tokenCodespace, 5, "token decimals should be within the range in 0 ~ 18")
	ErrInvalidAmount            = sdkerrors.Register(tokenCodespace, 6, "invalid token amount")
	ErrInvalidImageURILength    = sdkerrors.Register(tokenCodespace, 7, "invalid token uri length")
	ErrInvalidNameLength        = sdkerrors.Register(tokenCodespace, 8, "invalid name length")
	ErrInvalidTokenSymbol       = sdkerrors.Register(tokenCodespace, 9, "invalid token symbol")
	ErrTokenNoPermission        = sdkerrors.Register(tokenCodespace, 10, "account does not have the permission")
	ErrAccountExist             = sdkerrors.Register(tokenCodespace, 11, "account already exists")  // Currently never used
	ErrAccountNotExist          = sdkerrors.Register(tokenCodespace, 12, "account does not exists") // Currently never used
	ErrInsufficientBalance      = sdkerrors.Register(tokenCodespace, 13, "insufficient balance")
	ErrSupplyExist              = sdkerrors.Register(tokenCodespace, 14, "supply for token already exists") // Currently never used
	ErrInsufficientSupply       = sdkerrors.Register(tokenCodespace, 15, "insufficient supply")             // Currently never used
	ErrInvalidChangesFieldCount = sdkerrors.Register(tokenCodespace, 16, "invalid count of field changes")  // Currently never used
	ErrEmptyChanges             = sdkerrors.Register(tokenCodespace, 17, "changes is empty")
	ErrInvalidChangesField      = sdkerrors.Register(tokenCodespace, 18, "invalid field of changes")
	ErrDuplicateChangesField    = sdkerrors.Register(tokenCodespace, 19, "invalid field of changes")
	ErrInvalidMetaLength        = sdkerrors.Register(tokenCodespace, 20, "invalid meta length")
	ErrSupplyOverflow           = sdkerrors.Register(tokenCodespace, 21, "supply for token reached maximum") // Currently never used
	ErrApproverProxySame        = sdkerrors.Register(tokenCodespace, 22, "approver is same with proxy")
	ErrTokenNotApproved         = sdkerrors.Register(tokenCodespace, 23, "proxy is not approved on the token")
	ErrTokenAlreadyApproved     = sdkerrors.Register(tokenCodespace, 24, "proxy is already approved on the token")
)
