package types

var (
	EventTypePublishToken    = "publish_token"
	EventTypeMintToken       = "mint_token"
	EventTypeBurnToken       = "burn_token"
	EventTypeGrantPermToken  = "grant_perm"
	EventTypeRevokePermToken = "revoke_perm"

	AttributeKeyName     = "name"
	AttributeKeySymbol   = "symbol"
	AttributeKeyOwner    = "owner"
	AttributeKeyAmount   = "amount"
	AttributeKeyMintable = "mintable"

	AttributeKeyFrom = "from"
	AttributeKeyTo   = "to"

	AttributeKeyResource = "resource"
	AttributeKeyAction   = "action"

	AttributeValueCategory = ModuleName
)
