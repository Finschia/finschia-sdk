package types

const (
	EventTypeStoreCode            = "store_code"
	EventTypeInstantiateContract  = "instantiate_contract"
	EventTypeExecuteContract      = "execute_contract"
	EventTypeMigrateContract      = "migrate_contract"
	EventTypeUpdateAdmin          = "update_admin"
	EventTypeClearAdmin           = "clear_admin"
	EventTypePinCode              = "pin_code"
	EventTypeUnpinCode            = "unpin_code"
	EventTypeUpdateContractStatus = "update_contract_status"
)
const ( // event attributes
	AttributeKeyContract       = "contract_address"
	AttributeKeyCodeID         = "code_id"
	AttributeKeyCodeIDs        = "code_ids"
	AttributeKeySigner         = "signer"
	AttributeKeyContractStatus = "contract_status"
)
