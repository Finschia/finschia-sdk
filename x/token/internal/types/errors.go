package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	CodeTokenExist              sdk.CodeType = 1
	CodeTokenNotExist           sdk.CodeType = 2
	CodeTokenNotMintable        sdk.CodeType = 3
	CodeTokenPermission         sdk.CodeType = 4
	CodeTokenPermissionMint     sdk.CodeType = 5
	CodeTokenInvalidTokenName   sdk.CodeType = 6
	CodeTokenInvalidTokenSymbol sdk.CodeType = 7
	CodeTokenInvalidTokenID     sdk.CodeType = 8
	CodeTokenInvalidDecimals    sdk.CodeType = 9
	CodeTokenNFTExist           sdk.CodeType = 10
	CodeCollectionDenomExist    sdk.CodeType = 11
	CodeCollectionDenomNotExist sdk.CodeType = 12
	CodeTokenInvalidFT          sdk.CodeType = 13
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

func ErrTokenInvalidDecimals(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeTokenInvalidDecimals, "invalid decimals")
}

func ErrInvalidTokenName(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeTokenInvalidTokenName, "invalid token name")
}

func ErrInvalidTokenSymbol(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeTokenInvalidTokenSymbol, "invalid symbol format")
}

func ErrInvalidTokenID(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeTokenInvalidTokenID, "invalid token-id format")
}

func ErrTokenNFTExist(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeTokenNFTExist, "the supply of the nft is not 0")
}

func ErrCollectionExist(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeCollectionDenomExist, "collection already exist")
}

func ErrCollectionNotExist(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeCollectionDenomNotExist, "collection is not exist")
}

func ErrInvalidFTExist(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeTokenInvalidFT, "invalid ft setting. issue nft instead")
}
