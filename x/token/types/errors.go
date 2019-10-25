package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	CodeTokenExist          sdk.CodeType = 1
	CodeTokenNotExist       sdk.CodeType = 2
	CodeTokenNotMintable    sdk.CodeType = 3
	CodeTokenPermission     sdk.CodeType = 4
	CodeTokenPermissionMint sdk.CodeType = 5
	CodeTokenPermissionBurn sdk.CodeType = 6
	CodeTokenSymbolLength   sdk.CodeType = 7
)

func ErrTokenExist(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeTokenExist, "token symbol already exist")
}

func ErrTokenNotExist(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeTokenNotExist, "token symbol is not exist")
}

func ErrTokenNotMintable(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeTokenNotMintable, "token is not mintable")
}

func ErrTokenPermission(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeTokenPermission, "account does not have permissions")
}

func ErrTokenPermissionMint(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeTokenPermissionMint, "account does not have permissions to mint tokens")
}

func ErrTokenPermissionBurn(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeTokenPermissionBurn, "account does not have permissions to burn tokens")
}

func ErrTokenSymbolLength(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeTokenSymbolLength, "the length of token symbol should be 6+")
}
