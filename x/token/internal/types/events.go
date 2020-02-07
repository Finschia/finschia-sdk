package types

var (
	EventTypeIssueToken       = "issue"
	EventTypeIssueCFT         = "issue_cft"
	EventTypeIssueCNFT        = "issue_cnft"
	EventTypeMintToken        = "mint"
	EventTypeBurnToken        = "burn"
	EventTypeMintCFT          = "mint_cft"
	EventTypeBurnCFT          = "burn_cft"
	EventTypeMintCNFT         = "mint_cnft"
	EventTypeModifyTokenURI   = "modify_token_uri_token" /* #nosec */
	EventTypeGrantPermToken   = "grant_perm"
	EventTypeRevokePermToken  = "revoke_perm"
	EventTypeCreateCollection = "create_collection"
	EventTypeAttachToken      = "attach_cnft" /* #nosec */
	EventTypeDetachToken      = "detach_cnft" /* #nosec */
	EventTypeTransfer         = "transfer"
	EventTypeTransferCFT      = "transfer_cft"
	EventTypeTransferCNFT     = "transfer_cnft"

	AttributeKeyName      = "name"
	AttributeKeySymbol    = "symbol"
	AttributeKeyTokenID   = "token_id"
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

	AttributeValueCategory = ModuleName
)
