package types

var (
	EventTypeIssueToken     = "issue_token"
	EventTypeMintToken      = "mint_token"
	EventTypeBurnToken      = "burn_token"
	EventTypeGrantPermToken = "grant_perm_token"
	//EventTypeModifyTokenURI #nosec
	EventTypeModifyTokenURI = "modify_token_uri_token"
	//EventTypeModifyTokenURIPermToken #nosec
	EventTypeModifyTokenURIPermToken = "modify_token_uri_perm_token"
	//EventTypeRevokePermToken #nosec
	EventTypeRevokePermToken  = "revoke_perm_token"
	EventTypeCreateCollection = "create_collection"
	EventTypeAttachToken      = "attach_token"
	EventTypeDetachToken      = "detach_token"
	EventTypeTransfer         = "transfer"
	EventTypeTransferCFT      = "transfer_cft"
	EventTypeTransferNFT      = "transfer_nft"
	EventTypeTransferCNFT     = "transfer_cnft"

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

	AttributeValueTokenTypeFT   = "ft"
	AttributeValueTokenTypeNFT  = "nft"
	AttributeValueTokenTypeCFT  = "cft"
	AttributeValueTokenTypeCNFT = "cnft"

	AttributeValueCategory = ModuleName
)
