package collection

import (
	errorsmod "cosmossdk.io/errors"
)

const collectionCodespace = ModuleName

var (
	ErrTokenNotExist                 = errorsmod.Register(collectionCodespace, 2, "token symbol, token-id does not exist")
	ErrTokenNotMintable              = errorsmod.Register(collectionCodespace, 3, "token symbol, token-id is not mintable")
	ErrInvalidTokenName              = errorsmod.Register(collectionCodespace, 4, "token name should not be empty")
	ErrInvalidTokenID                = errorsmod.Register(collectionCodespace, 5, "invalid token id")
	ErrInvalidTokenDecimals          = errorsmod.Register(collectionCodespace, 6, "token decimals should be within the range in 0 ~ 18")
	ErrInvalidIssueFT                = errorsmod.Register(collectionCodespace, 7, "Issuing token with amount[1], decimals[0], mintable[false] prohibited. Issue nft token instead.")
	ErrInvalidAmount                 = errorsmod.Register(collectionCodespace, 8, "invalid token amount")
	ErrInvalidBaseImgURILength       = errorsmod.Register(collectionCodespace, 9, "invalid base_img_uri length")
	ErrInvalidNameLength             = errorsmod.Register(collectionCodespace, 10, "invalid name length")
	ErrInvalidTokenType              = errorsmod.Register(collectionCodespace, 11, "invalid token type pattern found")
	ErrInvalidTokenIndex             = errorsmod.Register(collectionCodespace, 12, "invalid token index pattern found")
	ErrCollectionNotExist            = errorsmod.Register(collectionCodespace, 14, "collection does not exists")
	ErrTokenTypeNotExist             = errorsmod.Register(collectionCodespace, 16, "token type for contract_id, token-type does not exist")
	ErrTokenNoPermission             = errorsmod.Register(collectionCodespace, 20, "account does not have the permission")
	ErrTokenAlreadyAChild            = errorsmod.Register(collectionCodespace, 21, "token is already a child of some other")
	ErrTokenNotAChild                = errorsmod.Register(collectionCodespace, 22, "token is not a child of some other")
	ErrTokenNotOwnedBy               = errorsmod.Register(collectionCodespace, 23, "token is being not owned by")
	ErrTokenCannotTransferChildToken = errorsmod.Register(collectionCodespace, 24, "cannot transfer a child token")
	ErrTokenNotNFT                   = errorsmod.Register(collectionCodespace, 25, "token is not a NFT")
	ErrCannotAttachToItself          = errorsmod.Register(collectionCodespace, 26, "cannot attach token to itself")
	ErrCannotAttachToADescendant     = errorsmod.Register(collectionCodespace, 27, "cannot attach token to a descendant")
	ErrApproverProxySame             = errorsmod.Register(collectionCodespace, 28, "approver is same with proxy")
	ErrCollectionNotApproved         = errorsmod.Register(collectionCodespace, 29, "proxy is not approved on the collection")
	ErrCollectionAlreadyApproved     = errorsmod.Register(collectionCodespace, 30, "proxy is already approved on the collection")
	ErrInvalidChangesFieldCount      = errorsmod.Register(collectionCodespace, 35, "invalid count of field changes")
	ErrEmptyChanges                  = errorsmod.Register(collectionCodespace, 36, "changes is empty")
	ErrInvalidChangesField           = errorsmod.Register(collectionCodespace, 37, "invalid field of changes")
	ErrTokenIndexWithoutType         = errorsmod.Register(collectionCodespace, 38, "There is a token index but no token type")
	ErrTokenTypeFTWithoutIndex       = errorsmod.Register(collectionCodespace, 39, "There is a token type of ft but no token index")
	ErrInsufficientToken             = errorsmod.Register(collectionCodespace, 40, "insufficient token")
	ErrDuplicateChangesField         = errorsmod.Register(collectionCodespace, 41, "duplicate field of changes")
	ErrInvalidMetaLength             = errorsmod.Register(collectionCodespace, 42, "invalid meta length")
	ErrEmptyField                    = errorsmod.Register(collectionCodespace, 44, "required field cannot be empty")
	ErrCompositionTooDeep            = errorsmod.Register(collectionCodespace, 45, "cannot attach token (composition too deep)")
	ErrCompositionTooWide            = errorsmod.Register(collectionCodespace, 46, "cannot attach token (composition too wide)")
	ErrBurnNonRootNFT                = errorsmod.Register(collectionCodespace, 47, "cannot burn non-root NFTs")
)

// Legacy Token module error starts from 101
var (
	ErrInvalidPermission = errorsmod.Register(collectionCodespace, 101, "invalid permission")      // "link", 2
	ErrInvalidContractID = errorsmod.Register(collectionCodespace, 102, "invalid contractID")      // "contract", 2
	ErrContractNotExist  = errorsmod.Register(collectionCodespace, 103, "contract does not exist") // "contract", 3
)

// Deprecated: do not use from v0.50.x
var (
	ErrCollectionExist    = errorsmod.Register(collectionCodespace, 13, "collection already exists")
	ErrTokenTypeExist     = errorsmod.Register(collectionCodespace, 15, "token type for contract_id, token-type already exists")
	ErrTokenTypeFull      = errorsmod.Register(collectionCodespace, 17, "all token type for contract_id are occupied")
	ErrTokenIndexFull     = errorsmod.Register(collectionCodespace, 18, "all non-fungible token index for contract_id, token-type are occupied")
	ErrTokenIDFull        = errorsmod.Register(collectionCodespace, 19, "all fungible token-id for contract_id are occupied")
	ErrAccountExist       = errorsmod.Register(collectionCodespace, 31, "account already exists")
	ErrAccountNotExist    = errorsmod.Register(collectionCodespace, 32, "account does not exists")
	ErrInsufficientSupply = errorsmod.Register(collectionCodespace, 33, "insufficient supply")
	ErrInvalidCoin        = errorsmod.Register(collectionCodespace, 34, "invalid coin")
	ErrSupplyOverflow     = errorsmod.Register(collectionCodespace, 43, "supply for collection reached maximum")
)
