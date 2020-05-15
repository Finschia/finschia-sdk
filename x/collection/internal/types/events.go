package types

var (
	EventTypeIssueFT              = "issue_ft"
	EventTypeIssueNFT             = "issue_nft"
	EventTypeMintFT               = "mint_ft"
	EventTypeBurnFT               = "burn_ft"
	EventTypeMintNFT              = "mint_nft"
	EventTypeModifyCollection     = "modify_collection"
	EventTypeModifyTokenType      = "modify_token_type" /* #nosec */
	EventTypeModifyToken          = "modify_token"
	EventTypeGrantPermToken       = "grant_perm"
	EventTypeRevokePermToken      = "revoke_perm"
	EventTypeCreateCollection     = "create_collection"
	EventTypeAttachToken          = "attach" /* #nosec */
	EventTypeDetachToken          = "detach" /* #nosec */
	EventTypeAttachFrom           = "attach_from"
	EventTypeDetachFrom           = "detach_from"
	EventTypeTransfer             = "transfer"
	EventTypeTransferFT           = "transfer_ft"
	EventTypeTransferNFT          = "transfer_nft"
	EventTypeTransferFTFrom       = "transfer_ft_from"
	EventTypeTransferNFTFrom      = "transfer_nft_from"
	EventTypeOperationTransferNFT = "operation_transfer_nft"
	EventTypeApproveCollection    = "approve_collection"
	EventTypeDisapproveCollection = "disapprove_collection"
	EventTypeBurnNFT              = "burn_nft"
	EventTypeBurnFTFrom           = "burn_ft_from"
	EventTypeBurnNFTFrom          = "burn_nft_from"
	EventTypeOperationBurnNFT     = "operation_burn_nft"
	EventTypeOperationRootChanged = "operation_root_changed"

	AttributeKeyName        = "name"
	AttributeKeyMeta        = "meta"
	AttributeKeyContractID  = "contract_id"
	AttributeKeyTokenID     = "token_id"
	AttributeKeyOwner       = "owner"
	AttributeKeyAmount      = "amount"
	AttributeKeyDecimals    = "decimals"
	AttributeKeyBaseImgURI  = "base_img_uri"
	AttributeKeyMintable    = "mintable"
	AttributeKeyTokenType   = "token_type"
	AttributeKeyFrom        = "from"
	AttributeKeyTo          = "to"
	AttributeKeyPerm        = "perm"
	AttributeKeyToTokenID   = "to_token_id"
	AttributeKeyFromTokenID = "from_token_id"
	AttributeKeyApprover    = "approver"
	AttributeKeyProxy       = "proxy"
	AttributeKeyOldRoot     = "old_root_token_id"
	AttributeKeyNewRoot     = "new_root_token_id"

	AttributeValueCategory = ModuleName
)
