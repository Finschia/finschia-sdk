package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	//Token
	CodeTokenExist       sdk.CodeType = 100
	CodeTokenNotExist    sdk.CodeType = 101
	CodeTokenNotMintable sdk.CodeType = 102

	//Token invalidation
	CodeTokenInvalidTokenName   sdk.CodeType = 200
	CodeTokenInvalidTokenSymbol sdk.CodeType = 201
	CodeTokenInvalidTokenID     sdk.CodeType = 202
	CodeTokenInvalidDecimals    sdk.CodeType = 203
	CodeTokenInvalidFT          sdk.CodeType = 204

	//Collection
	CodeCollectionDenomExist    sdk.CodeType = 300
	CodeCollectionDenomNotExist sdk.CodeType = 301

	//Permission
	CodeTokenPermission sdk.CodeType = 400
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

func ErrInvalidTokenID(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeTokenInvalidTokenID, msg)
}

func ErrInvalidTokenDecimals(codespace sdk.CodespaceType, decimals sdk.Int) sdk.Error {
	return sdk.NewError(codespace, CodeTokenInvalidDecimals, "token decimals [%s] should be within the range in 0 ~ 18", decimals.String())
}

func ErrInvalidIssueFT(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeTokenInvalidFT, "Issuing token with amount[1], decimals[0], mintable[false] prohibited. Issue nft token instead.")
}

func ErrCollectionExist(codespace sdk.CodespaceType, symbol string) sdk.Error {
	return sdk.NewError(codespace, CodeCollectionDenomExist, "collection [%s] already exists", symbol)
}

func ErrCollectionNotExist(codespace sdk.CodespaceType, symbol string) sdk.Error {
	return sdk.NewError(codespace, CodeCollectionDenomNotExist, "collection [%s] does not exists", symbol)
}

func ErrTokenPermission(codespace sdk.CodespaceType, account sdk.AccAddress, permission PermissionI) sdk.Error {
	return sdk.NewError(codespace, CodeTokenPermission, "account [%s] does not have the permission [%s]", account.String(), permission.String())
}
