package types

var (
	EventTypeIssueToken      = "issue"
	EventTypeMintToken       = "mint"
	EventTypeBurnToken       = "burn"
	EventTypeModifyTokenURI  = "modify_token_uri_token" /* #nosec */
	EventTypeGrantPermToken  = "grant_perm"
	EventTypeRevokePermToken = "revoke_perm"
	EventTypeTransfer        = "transfer"

	AttributeKeyName     = "name"
	AttributeKeySymbol   = "symbol"
	AttributeKeyOwner    = "owner"
	AttributeKeyAmount   = "amount"
	AttributeKeyDecimals = "decimals"
	AttributeKeyTokenURI = "token_uri"
	AttributeKeyMintable = "mintable"
	AttributeKeyFrom     = "from"
	AttributeKeyTo       = "to"
	AttributeKeyResource = "perm_resource"
	AttributeKeyAction   = "perm_action"

	AttributeValueCategory = ModuleName
)
