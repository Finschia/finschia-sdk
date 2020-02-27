package types

import (
	"fmt"
	"unicode/utf8"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	// Token
	CodeTokenExist       sdk.CodeType = 100
	CodeTokenNotExist    sdk.CodeType = 101
	CodeTokenNotMintable sdk.CodeType = 102

	// Token invalidation
	CodeTokenInvalidTokenName      sdk.CodeType = 200
	CodeTokenInvalidTokenSymbol    sdk.CodeType = 201
	CodeTokenInvalidDecimals       sdk.CodeType = 202
	CodeTokenInvalidAmount         sdk.CodeType = 203
	CodeTokenInvalidTokenURILength sdk.CodeType = 204
	CodeTokenInvalidNameLength     sdk.CodeType = 205

	// Permission
	CodePermission sdk.CodeType = 300

	// Account
	CodeAccountExist    sdk.CodeType = 400
	CodeAccountNotExist sdk.CodeType = 401

	// Bank
	CodeInsufficientBalance sdk.CodeType = 500
	CodeInvalidAmount       sdk.CodeType = 501

	// Supply
	CodeSupplyExist        sdk.CodeType = 600
	CodeInsufficientSupply sdk.CodeType = 601

	// Modify
	CodeInvalidChangesFieldCount sdk.CodeType = 701
	CodeEmptyChanges             sdk.CodeType = 702
	CodeTokenInvalidChangesField sdk.CodeType = 703
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

func ErrInvalidChangesFieldCount(codespace sdk.CodespaceType, changesFieldCount int) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidChangesFieldCount,
		"You can not change fields more than [%d] at once, current count: [%d]", MaxChangeFieldsCount, changesFieldCount)
}

func ErrEmptyChanges(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeEmptyChanges, "changes is empty")
}

func ErrInvalidTokenURILength(codespace sdk.CodespaceType, tokenURI string) sdk.Error {
	return sdk.NewError(codespace, CodeTokenInvalidTokenURILength, "invalid token uri [%s] should be shorter than [%d] UTF-8 characters, current length: [%d]", tokenURI, MaxTokenURILength, utf8.RuneCountInString(tokenURI))
}

func ErrInvalidNameLength(codespace sdk.CodespaceType, name string) sdk.Error {
	return sdk.NewError(codespace, CodeTokenInvalidNameLength,
		"invalid name [%s] should be shorter than [%d] UTF-8 characters, current length: [%d]", name,
		MaxTokenNameLength, utf8.RuneCountInString(name))
}

func ErrInvalidChangesField(codespace sdk.CodespaceType, field string) sdk.Error {
	return sdk.NewError(codespace, CodeTokenInvalidChangesField, "[%s] is invalid field of changes", field)
}

func ErrTokenNoPermission(codespace sdk.CodespaceType, account fmt.Stringer, permission fmt.Stringer) sdk.Error {
	return sdk.NewError(codespace, CodePermission, "account [%s] does not have the permission [%s]", account.String(), permission.String())
}
func ErrAccountExist(codespace sdk.CodespaceType, acc sdk.AccAddress) sdk.Error {
	return sdk.NewError(codespace, CodeAccountExist, "account [%s] already exists", acc.String())
}

func ErrAccountNotExist(codespace sdk.CodespaceType, acc sdk.AccAddress) sdk.Error {
	return sdk.NewError(codespace, CodeAccountNotExist, "account [%s] does not exists", acc.String())
}

func ErrInsufficientBalance(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInsufficientBalance, msg)
}

func ErrSupplyExist(codespace sdk.CodespaceType, symbol string) sdk.Error {
	return sdk.NewError(codespace, CodeSupplyExist, "supply for token [%s] already exists", symbol)
}

func ErrInsufficientSupply(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInsufficientSupply, msg)
}
