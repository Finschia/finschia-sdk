package types

var (
	EventTypeIssueToken              = "issue_token"
	EventTypeMintToken               = "mint_token"
	EventTypeBurnToken               = "burn_token"
	EventTypeGrantPermToken          = "grant_perm_token"
	EventTypeModifyTokenURI          = "modify_token_uri_token"
	EventTypeModifyTokenURIPermToken = "modify_token_uri_perm_token"
	EventTypeRevokePermToken         = "revoke_perm_token"
	EventTypeOccupySymbol            = "occupy_symbol"
	EventTypeAttachToken             = "attach_token"
	EventTypeDetachToken             = "detach_token"
	EventTypeTransfer                = "transfer"
	EventTypeTransferIDFT            = "transfer_idft"
	EventTypeTransferNFT             = "transfer_nft"
	EventTypeTransferIDNFT           = "transfer_idnft"

	AttributeKeyName      = "name"
	AttributeKeySymbol    = "symbol"
	AttributeKeyTokenID   = "token_id"
	AttributeKeyDenom     = "denom"
	AttributeKeyOwner     = "owner"
	AttributeKeyAmount    = "amount"
	AttributeKeyDecimals  = "decimals"
	AttributeKeyTokenURI  = "token_uri"
	AttributeKeyMintable  = "mintable"
	AttributeKeyTokenType = "token_type"
	AttributeKeyFrom      = "from"
	AttributeKeyTo        = "to"
	AttributeKeyResource  = "perm_resource"
	AttributeKeyAction    = "perm_action"
	AttributeKeyToTokenID = "to_token_id"

	AttributeValueTokenTypeFT    = "ft"
	AttributeValueTokenTypeNFT   = "nft"
	AttributeValueTokenTypeIDFT  = "idft"
	AttributeValueTokenTypeIDNFT = "idnft"

	AttributeValueCategory = ModuleName
)
