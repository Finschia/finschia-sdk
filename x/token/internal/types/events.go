package types

var (
	EventTypeIssueToken      = "issue"
	EventTypeMintToken       = "mint"
	EventTypeBurnToken       = "burn"
	EventTypeBurnTokenFrom   = "burn_from"
	EventTypeModifyToken     = "modify_token"
	EventTypeGrantPermToken  = "grant_perm"
	EventTypeRevokePermToken = "revoke_perm"
	EventTypeTransfer        = "transfer"
	EventTypeTransferFrom    = "transfer_from"
	EventTypeApproveToken    = "approve_token"

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
	AttributeKeyApprover   = "approver"
	AttributeKeyProxy      = "proxy"
	AttributeValueCategory = ModuleName
)
