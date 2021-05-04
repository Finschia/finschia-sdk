package types

const (
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
