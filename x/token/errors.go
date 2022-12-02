package token

import (
	sdkerrors "github.com/line/lbm-sdk/types/errors"
)

const tokenCodespace = ModuleName

var (
	ErrInvalidContractID          = sdkerrors.Register(tokenCodespace, 2, "invalid contract id")
	ErrContractNotFound           = sdkerrors.Register(tokenCodespace, 3, "contract not found")
	ErrInvalidPermission          = sdkerrors.Register(tokenCodespace, 4, "invalid permission")
	ErrGrantNotFound              = sdkerrors.Register(tokenCodespace, 5, "grant not found")
	ErrGrantAlreadyExists         = sdkerrors.Register(tokenCodespace, 6, "grant already exists")
	ErrOperatorIsHolder           = sdkerrors.Register(tokenCodespace, 7, "operator and holder should be different")
	ErrAuthorizationNotFound      = sdkerrors.Register(tokenCodespace, 8, "authorization not found")
	ErrAuthorizationAlreadyExists = sdkerrors.Register(tokenCodespace, 9, "authorization already exists")
	ErrInvalidAmount              = sdkerrors.Register(tokenCodespace, 10, "invalid amount")
	ErrInsufficientTokens         = sdkerrors.Register(tokenCodespace, 11, "insufficient tokens")
	ErrInvalidName                = sdkerrors.Register(tokenCodespace, 12, "invalid name")
	ErrInvalidSymbol              = sdkerrors.Register(tokenCodespace, 13, "invalid symbol")
	ErrInvalidImageURI            = sdkerrors.Register(tokenCodespace, 14, "invalid image_uri")
	ErrInvalidMeta                = sdkerrors.Register(tokenCodespace, 15, "invalid meta")
	ErrInvalidDecimals            = sdkerrors.Register(tokenCodespace, 16, "invalid decimals")
	ErrNotMintable                = sdkerrors.Register(tokenCodespace, 17, "not mintable")
	ErrInvalidChanges             = sdkerrors.Register(tokenCodespace, 18, "invalid changes")
)
