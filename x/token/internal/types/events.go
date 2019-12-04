package types

var (
	EventTypeIssueToken      = "issue_token"
	EventTypeMintToken       = "mint_token"
	EventTypeBurnToken       = "burn_token"
	EventTypeGrantPermToken  = "grant_perm_token"
	EventTypeRevokePermToken = "revoke_perm_token"
	EventTypeOccupySymbol    = "occupy_symbol"

	AttributeKeyName      = "name"
	AttributeKeySymbol    = "symbol"
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

	AttributeValueTokenTypeFT   = "ft"
	AttributeValueTokenTypeNFT  = "nft"
	AttributeValueTokenTypeCFT  = "cft"
	AttributeValueTokenTypeCNFT = "cnft"

	AttributeValueCategory = ModuleName
)
