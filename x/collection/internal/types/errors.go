package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrTokenExist                    = sdkerrors.Register(ModuleName, 1, "token symbol, token-id already exists")
	ErrTokenNotExist                 = sdkerrors.Register(ModuleName, 2, "token symbol, token-id does not exist")
	ErrTokenNotMintable              = sdkerrors.Register(ModuleName, 3, "token symbol, token-id is not mintable")
	ErrInvalidTokenName              = sdkerrors.Register(ModuleName, 4, "token name should not be empty")
	ErrInvalidTokenID                = sdkerrors.Register(ModuleName, 5, "invalid token id")
	ErrInvalidTokenDecimals          = sdkerrors.Register(ModuleName, 6, "token decimals should be within the range in 0 ~ 18")
	ErrInvalidIssueFT                = sdkerrors.Register(ModuleName, 7, "Issuing token with amount[1], decimals[0], mintable[false] prohibited. Issue nft token instead.")
	ErrInvalidAmount                 = sdkerrors.Register(ModuleName, 8, "invalid token amount")
	ErrInvalidBaseImgURILength       = sdkerrors.Register(ModuleName, 9, "invalid base_img_uri length")
	ErrInvalidNameLength             = sdkerrors.Register(ModuleName, 10, "invalid name length")
	ErrInvalidTokenType              = sdkerrors.Register(ModuleName, 11, "invalid token type pattern found")
	ErrInvalidTokenIndex             = sdkerrors.Register(ModuleName, 12, "invalid token index pattern found")
	ErrCollectionExist               = sdkerrors.Register(ModuleName, 13, "collection already exists")
	ErrCollectionNotExist            = sdkerrors.Register(ModuleName, 14, "collection does not exists")
	ErrTokenTypeExist                = sdkerrors.Register(ModuleName, 15, "token type for contract_id, token-type already exists")
	ErrTokenTypeNotExist             = sdkerrors.Register(ModuleName, 16, "token type for contract_id, token-type does not exist")
	ErrTokenTypeFull                 = sdkerrors.Register(ModuleName, 17, "all token type for contract_id are occupied")
	ErrTokenIndexFull                = sdkerrors.Register(ModuleName, 18, "all non-fungible token index for contract_id, token-type are occupied")
	ErrTokenIDFull                   = sdkerrors.Register(ModuleName, 19, "all fungible token-id for contract_id are occupied")
	ErrTokenNoPermission             = sdkerrors.Register(ModuleName, 20, "account does not have the permission")
	ErrTokenAlreadyAChild            = sdkerrors.Register(ModuleName, 21, "token is already a child of some other")
	ErrTokenNotAChild                = sdkerrors.Register(ModuleName, 22, "token is not a child of some other")
	ErrTokenNotOwnedBy               = sdkerrors.Register(ModuleName, 23, "token is being not owned by")
	ErrTokenCannotTransferChildToken = sdkerrors.Register(ModuleName, 24, "cannot transfer a child token")
	ErrTokenNotNFT                   = sdkerrors.Register(ModuleName, 25, "token is not a NFT")
	ErrCannotAttachToItself          = sdkerrors.Register(ModuleName, 26, "cannot attach token to itself")
	ErrCannotAttachToADescendant     = sdkerrors.Register(ModuleName, 27, "cannot attach token to a descendant")
	ErrApproverProxySame             = sdkerrors.Register(ModuleName, 28, "approver is same with proxy")
	ErrCollectionNotApproved         = sdkerrors.Register(ModuleName, 29, "proxy is not approved on the collection")
	ErrCollectionAlreadyApproved     = sdkerrors.Register(ModuleName, 30, "proxy is already approved on the collection")
	ErrAccountExist                  = sdkerrors.Register(ModuleName, 31, "account already exists")
	ErrAccountNotExist               = sdkerrors.Register(ModuleName, 32, "account does not exists")
	ErrInsufficientSupply            = sdkerrors.Register(ModuleName, 33, "insufficient supply")
	ErrInvalidCoin                   = sdkerrors.Register(ModuleName, 34, "invalid coin")
	ErrInvalidChangesFieldCount      = sdkerrors.Register(ModuleName, 35, "invalid count of field changes")
	ErrEmptyChanges                  = sdkerrors.Register(ModuleName, 36, "changes is empty")
	ErrInvalidChangesField           = sdkerrors.Register(ModuleName, 37, "invalid field of changes")
	ErrTokenIndexWithoutType         = sdkerrors.Register(ModuleName, 38, "There is a token index but no token type")
	ErrTokenTypeFTWithoutIndex       = sdkerrors.Register(ModuleName, 39, "There is a token type of ft but no token index")
	ErrInsufficientToken             = sdkerrors.Register(ModuleName, 40, "insufficient token")
	ErrDuplicateChangesField         = sdkerrors.Register(ModuleName, 41, "duplicate field of changes")
	ErrInvalidMetaLength             = sdkerrors.Register(ModuleName, 42, "invalid meta length")
	ErrSupplyOverflow                = sdkerrors.Register(ModuleName, 43, "supply for collection reached maximum")
	ErrEmptyField                    = sdkerrors.Register(ModuleName, 44, "required field cannot be empty")
	ErrCompositionTooDeep            = sdkerrors.Register(ModuleName, 45, "cannot attach token (composition too deep)")
	ErrCompositionTooWide            = sdkerrors.Register(ModuleName, 46, "cannot attach token (composition too wide)")
	ErrBurnNonRootNFT                = sdkerrors.Register(ModuleName, 47, "cannot burn non-root NFTs")
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
