package types

var (
	EventTypeIssueToken      = "issue"
	EventTypeMintToken       = "mint"
	EventTypeBurnToken       = "burn"
	EventTypeModifyToken     = "modify_token"
	EventTypeGrantPermToken  = "grant_perm"
	EventTypeRevokePermToken = "revoke_perm"
	EventTypeTransfer        = "transfer"

	AttributeKeyName       = "name"
	AttributeKeySymbol     = "symbol"
	AttributeKeyContractID = "contract_id"
	AttributeKeyOwner      = "owner"
	AttributeKeyAmount     = "amount"
	AttributeKeyDecimals   = "decimals"
	AttributeKeyMeta       = "meta"
	AttributeKeyImageURI   = "img_uri"
	AttributeKeyMintable   = "mintable"
	AttributeKeyFrom       = "from"
	AttributeKeyTo         = "to"
	AttributeKeyPerm       = "perm"
	AttributeValueCategory = ModuleName
)
