package types

import (
	"fmt"
	"unicode/utf8"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	//Token
	CodeTokenExist       sdk.CodeType = 100
	CodeTokenNotExist    sdk.CodeType = 101
	CodeTokenNotMintable sdk.CodeType = 102

	//Token invalidation
	CodeTokenInvalidTokenName      sdk.CodeType = 200
	CodeTokenInvalidTokenSymbol    sdk.CodeType = 201
	CodeTokenInvalidDecimals       sdk.CodeType = 202
	CodeTokenInvalidAmount         sdk.CodeType = 203
	CodeTokenInvalidTokenURILength sdk.CodeType = 204

	//Permission
	CodeTokenPermission sdk.CodeType = 300
)

func ErrTokenExist(codespace sdk.CodespaceType, symbol string) sdk.Error {
	return sdk.NewError(codespace, CodeTokenExist, "token [%s] already exists", symbol)
}

func ErrTokenNotExist(codespace sdk.CodespaceType, symbol string) sdk.Error {
	return sdk.NewError(codespace, CodeTokenNotExist, "token [%s] does not exist", symbol)
}

func ErrTokenNotMintable(codespace sdk.CodespaceType, symbol string) sdk.Error {
	return sdk.NewError(codespace, CodeTokenNotMintable, "token [%s] is not mintable", symbol)
}

func ErrInvalidTokenName(codespace sdk.CodespaceType, name string) sdk.Error {
	return sdk.NewError(codespace, CodeTokenInvalidTokenName, "token name [%s] should not be empty", name)
}

func ErrInvalidTokenSymbol(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeTokenInvalidTokenSymbol, msg)
}

func ErrInvalidTokenDecimals(codespace sdk.CodespaceType, decimals sdk.Int) sdk.Error {
	return sdk.NewError(codespace, CodeTokenInvalidDecimals, "token decimals [%s] should be within the range in 0 ~ 18", decimals.String())
}

func ErrInvalidAmount(codespace sdk.CodespaceType, amount string) sdk.Error {
	return sdk.NewError(codespace, CodeTokenInvalidAmount, "invalid token amount [%s]", amount)
}

func ErrInvalidTokenURILength(codespace sdk.CodespaceType, tokenURI string) sdk.Error {
	return sdk.NewError(codespace, CodeTokenInvalidTokenURILength, "invalid token uri [%s] should be shorter than [%d] UTF-8 characters, current length: [%d]", tokenURI, TokenURIMaxLength, utf8.RuneCountInString(tokenURI))
}

func ErrTokenNoPermission(codespace sdk.CodespaceType, account fmt.Stringer, permission fmt.Stringer) sdk.Error {
	return sdk.NewError(codespace, CodeTokenPermission, "account [%s] does not have the permission [%s]", account.String(), permission.String())
}
