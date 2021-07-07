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
	// WasmModuleEventType is stored with any contract TX
	WasmModuleEventType = "wasm"
	// CustomContractEventPrefix contracts can create custom events. To not mix them with other system events they got the `wasm-` prefix.
	CustomContractEventPrefix = "wasm-"
	EventTypePinCode          = "pin_code"
	EventTypeUnpinCode        = "unpin_code"
)

const ( // event attributes
	AttributeKeyContractAddr   = "contract_address"
	AttributeKeyCodeID         = "code_id"
	AttributeKeySigner         = "signer"
	AttributeResultDataHex     = "result"
	AttributeKeyCodeIDs        = "code_ids"
	AttributeKeyContractStatus = "contract_status"
)
