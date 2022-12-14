package token

import (
	sdkerrors "github.com/line/lbm-sdk/types/errors"
)

const tokenCodespace = ModuleName

var (
	ErrContractNotFound           = sdkerrors.Register(tokenCodespace, 3, "contract not found")
	ErrNotMintable                = sdkerrors.Register(tokenCodespace, 4, "not mintable")
	ErrInvalidDecimals            = sdkerrors.Register(tokenCodespace, 6, "invalid decimals")
	ErrInvalidAmount              = sdkerrors.Register(tokenCodespace, 7, "invalid amount")
	ErrInvalidImageURI            = sdkerrors.Register(tokenCodespace, 8, "invalid image_uri")
	ErrInvalidName                = sdkerrors.Register(tokenCodespace, 9, "invalid name")
	ErrInvalidSymbol              = sdkerrors.Register(tokenCodespace, 10, "invalid symbol")
	ErrGrantNotFound              = sdkerrors.Register(tokenCodespace, 11, "grant not found")
	ErrInsufficientTokens         = sdkerrors.Register(tokenCodespace, 14, "insufficient tokens")
	ErrInvalidChanges             = sdkerrors.Register(tokenCodespace, 17, "invalid changes")
	ErrInvalidMeta                = sdkerrors.Register(tokenCodespace, 21, "invalid meta")
	ErrOperatorIsHolder           = sdkerrors.Register(tokenCodespace, 23, "operator and holder should be different")
	ErrAuthorizationNotFound      = sdkerrors.Register(tokenCodespace, 24, "authorization not found")
	ErrAuthorizationAlreadyExists = sdkerrors.Register(tokenCodespace, 25, "authorization already exists")
)
