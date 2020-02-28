package types

var (
	EventCreateAccount = "create_account"
	EventEmpty         = "empty"

	MsgTypeCreateAccount = EventCreateAccount
	MsgTypeEmpty         = EventEmpty

	AttributeKeyCreateAccountFrom   = "create_account_from"
	AttributeKeyCreateAccountTarget = "create_account_target"

	AttributeKeyFrom = "from"

	AttributeValueCategory = ModuleName
)
