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
	CodeTokenInvalidTokenName        sdk.CodeType = 200
	CodeTokenInvalidTokenSymbol      sdk.CodeType = 201
	CodeTokenInvalidTokenID          sdk.CodeType = 202
	CodeTokenInvalidDecimals         sdk.CodeType = 203
	CodeTokenInvalidFT               sdk.CodeType = 204
	CodeTokenInvalidAmount           sdk.CodeType = 205
	CodeTokenInvalidBaseImgURILength sdk.CodeType = 206
	CodeTokenInvalidNameLength       sdk.CodeType = 207
	CodeTokenInvalidTokenType        sdk.CodeType = 208
	CodeTokenInvalidTokenIndex       sdk.CodeType = 209

	//Collection
	CodeCollectionExist             sdk.CodeType = 300
	CodeCollectionNotExist          sdk.CodeType = 301
	CodeCollectionTokenTypeExist    sdk.CodeType = 302
	CodeCollectionTokenTypeNotExist sdk.CodeType = 303
	CodeCollectionTokenTypeFull     sdk.CodeType = 304
	CodeCollectionTokenIndexFull    sdk.CodeType = 305
	CodeCollectionTokenIDFull       sdk.CodeType = 306

	//Permission
	CodeTokenPermission sdk.CodeType = 400

	// Composability
	CodeTokenAlreadyAChild             sdk.CodeType = 500
	CodeTokenNotAChild                 sdk.CodeType = 501
	CodeTokenNotOwnedBy                sdk.CodeType = 502
	CodeTokenChildNotTransferable      sdk.CodeType = 503
	CodeTokenNotIDNF                   sdk.CodeType = 504
	CodeTokenCannotAttachToItself      sdk.CodeType = 505
	CodeTokenCannotAttachToADescendant sdk.CodeType = 506

	// Proxy
	CodeTokenApproverProxySame sdk.CodeType = 600
	CodeTokenNotApproved       sdk.CodeType = 601
	CodeTokenAlreadyApproved   sdk.CodeType = 602

	//Account
	CodeAccountExist    sdk.CodeType = 700
	CodeAccountNotExist sdk.CodeType = 701

	//Bank
	CodeInsufficientSupply sdk.CodeType = 800
	CodeInvalidCoin        sdk.CodeType = 801

	// Modify
	CodeInvalidChangesFieldCount sdk.CodeType = 901
	CodeEmptyChanges             sdk.CodeType = 902
	CodeTokenInvalidChangesField sdk.CodeType = 903
	CodeTokenIndexWithoutType    sdk.CodeType = 904
)

func ErrTokenNotMintable(codespace sdk.CodespaceType, symbol, tokenID string) sdk.Error {
	return sdk.NewError(codespace, CodeTokenNotMintable, "token symbol[%s] token-id[%s] is not mintable", symbol, tokenID)
}

func ErrInvalidTokenName(codespace sdk.CodespaceType, name string) sdk.Error {
	return sdk.NewError(codespace, CodeTokenInvalidTokenName, "token name [%s] should not be empty", name)
}

func ErrInvalidTokenSymbol(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeTokenInvalidTokenSymbol, "invalid symbol pattern found %s", msg)
}

func ErrInvalidTokenID(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeTokenInvalidTokenID, "invalid token id pattern found %s", msg)
}

func ErrInvalidTokenType(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeTokenInvalidTokenType, "invalid token type pattern found %s", msg)
}

func ErrInvalidTokenIndex(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeTokenInvalidTokenIndex, "invalid token index pattern found %s", msg)
}

func ErrInvalidTokenDecimals(codespace sdk.CodespaceType, decimals sdk.Int) sdk.Error {
	return sdk.NewError(codespace, CodeTokenInvalidDecimals, "token decimals [%s] should be within the range in 0 ~ 18", decimals.String())
}

func ErrInvalidIssueFT(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeTokenInvalidFT, "Issuing token with amount[1], decimals[0], mintable[false] prohibited. Issue nft token instead.")
}

func ErrInvalidAmount(codespace sdk.CodespaceType, amount string) sdk.Error {
	return sdk.NewError(codespace, CodeTokenInvalidAmount, "invalid token amount [%s]", amount)
}

func ErrInvalidChangesFieldCount(codespace sdk.CodespaceType, changesFieldCount int) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidChangesFieldCount,
		"You can not change fields more than [%d], current count: [%d]", MaxChangeFieldsCount, changesFieldCount)
}

func ErrEmptyChanges(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeEmptyChanges, "changes is empty")
}

func ErrInvalidBaseImgURILength(codespace sdk.CodespaceType, baseImgURI string) sdk.Error {
	return sdk.NewError(codespace, CodeTokenInvalidBaseImgURILength,
		"invalid base_img_uri [%s] should be shorter than [%d] UTF-8 characters, current length: [%d]", baseImgURI,
		MaxBaseImgURILength, utf8.RuneCountInString(baseImgURI))
}

func ErrInvalidNameLength(codespace sdk.CodespaceType, name string) sdk.Error {
	return sdk.NewError(codespace, CodeTokenInvalidNameLength,
		"invalid name [%s] should be shorter than [%d] UTF-8 characters, current length: [%d]", name,
		MaxTokenNameLength, utf8.RuneCountInString(name))
}

func ErrInvalidChangesField(codespace sdk.CodespaceType, field string) sdk.Error {
	return sdk.NewError(codespace, CodeTokenInvalidChangesField, "[%s] is invalid field of changes", field)
}

func ErrTokenIndexWithoutType(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeTokenIndexWithoutType, "There is a token index but no token type")
}

func ErrCollectionExist(codespace sdk.CodespaceType, symbol string) sdk.Error {
	return sdk.NewError(codespace, CodeCollectionExist, "collection [%s] already exists", symbol)
}

func ErrCollectionNotExist(codespace sdk.CodespaceType, symbol string) sdk.Error {
	return sdk.NewError(codespace, CodeCollectionNotExist, "collection [%s] does not exists", symbol)
}

func ErrTokenNoPermission(codespace sdk.CodespaceType, account fmt.Stringer, permission fmt.Stringer) sdk.Error {
	return sdk.NewError(codespace, CodeTokenPermission, "account [%s] does not have the permission [%s]", account.String(), permission.String())
}

func ErrTokenExist(codespace sdk.CodespaceType, symbol, tokenID string) sdk.Error {
	return sdk.NewError(codespace, CodeTokenExist, "token symbol[%s] token-id[%s] already exists", symbol, tokenID)
}

func ErrTokenNotExist(codespace sdk.CodespaceType, symbol, tokenID string) sdk.Error {
	return sdk.NewError(codespace, CodeTokenNotExist, "token symbol[%s] token-id[%s] does not exist", symbol, tokenID)
}

func ErrTokenTypeExist(codespace sdk.CodespaceType, symbol, tokenType string) sdk.Error {
	return sdk.NewError(codespace, CodeCollectionTokenTypeExist, "token type for symbol[%s] token-type[%s] already exists", symbol, tokenType)
}

func ErrTokenTypeNotExist(codespace sdk.CodespaceType, symbol, tokenType string) sdk.Error {
	return sdk.NewError(codespace, CodeCollectionTokenTypeNotExist, "token type for symbol[%s] token-type[%s] does not exist", symbol, tokenType)
}

func ErrTokenTypeFull(codespace sdk.CodespaceType, symbol string) sdk.Error {
	return sdk.NewError(codespace, CodeCollectionTokenTypeFull, "all token type for symbol[%s] are occupied", symbol)
}

func ErrTokenIndexFull(codespace sdk.CodespaceType, symbol, tokenType string) sdk.Error {
	return sdk.NewError(codespace, CodeCollectionTokenIndexFull, "all non-fungible token index for symbol[%s] token-type[%s] are occupied", symbol, tokenType)
}

func ErrTokenIDFull(codespace sdk.CodespaceType, symbol string) sdk.Error {
	return sdk.NewError(codespace, CodeCollectionTokenIDFull, "all fungible token-id for symbol[%s] are occupied", symbol)
}

func ErrTokenAlreadyAChild(codespace sdk.CodespaceType, tokenID string) sdk.Error {
	return sdk.NewError(codespace, CodeTokenAlreadyAChild, "token [%s] is already a child of some other", tokenID)
}

func ErrTokenNotAChild(codespace sdk.CodespaceType, tokenID string) sdk.Error {
	return sdk.NewError(codespace, CodeTokenNotAChild, "token [%s] is not a child of some other", tokenID)
}

func ErrTokenNotOwnedBy(codespace sdk.CodespaceType, tokenID string, owner fmt.Stringer) sdk.Error {
	return sdk.NewError(codespace, CodeTokenNotOwnedBy, "token is being not owned by [%s]", tokenID, owner.String())
}

func ErrTokenNotNFT(codespace sdk.CodespaceType, tokenID string) sdk.Error {
	return sdk.NewError(codespace, CodeTokenNotIDNF, "token [%s] is not a NFT", tokenID)
}

func ErrTokenCannotTransferChildToken(codespace sdk.CodespaceType, tokenID string) sdk.Error {
	return sdk.NewError(codespace, CodeTokenChildNotTransferable, "cannot transfer a child token [%s]", tokenID)
}

func ErrCannotAttachToItself(codespace sdk.CodespaceType, tokenID string) sdk.Error {
	return sdk.NewError(codespace, CodeTokenCannotAttachToItself, "cannot attach token [%s] to itself", tokenID)
}

func ErrCannotAttachToADescendant(codespace sdk.CodespaceType, tokenID string, tokenIDdesc string) sdk.Error {
	return sdk.NewError(codespace, CodeTokenCannotAttachToADescendant, "cannot attach token [%s] to a descendant [%s]", tokenID, tokenIDdesc)
}

func ErrApproverProxySame(codespace sdk.CodespaceType, approver string) sdk.Error {
	return sdk.NewError(codespace, CodeTokenApproverProxySame, "approver[%s] is same with proxy", approver)
}

func ErrCollectionNotApproved(codespace sdk.CodespaceType, proxy string, approver string, symbol string) sdk.Error {
	return sdk.NewError(codespace, CodeTokenNotApproved, "proxy[%s] is not approved by %s on the collection[%s]", proxy, approver, symbol)
}

func ErrCollectionAlreadyApproved(codespace sdk.CodespaceType, proxy string, approver string, symbol string) sdk.Error {
	return sdk.NewError(codespace, CodeTokenAlreadyApproved, "proxy[%s] is already approved by %s on the collection[%s]", proxy, approver, symbol)
}

func ErrAccountExist(codespace sdk.CodespaceType, acc sdk.AccAddress) sdk.Error {
	return sdk.NewError(codespace, CodeAccountExist, "account [%s] already exists", acc.String())
}

func ErrAccountNotExist(codespace sdk.CodespaceType, acc sdk.AccAddress) sdk.Error {
	return sdk.NewError(codespace, CodeAccountNotExist, "account [%s] does not exists", acc.String())
}

func ErrInsufficientSupply(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInsufficientSupply, msg)
}

func ErrInvalidCoin(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidCoin, msg)
}
