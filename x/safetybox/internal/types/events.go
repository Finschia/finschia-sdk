package types

var (
	EventSafetyBoxCreate     = "safety_box_create"
	EventSafetyBoxSendCoin   = "safety_box_send_coin"
	EventSafetyBoxPermission = "safety_box_permission"

	MsgTypeSafetyBoxCreate                    = "safety_box_create"
	MsgTypeSafetyBoxAllocateCoin              = "safety_box_allocate_coin"
	MsgTypeSafetyBoxRecallCoin                = "safety_box_recall_coin"
	MsgTypeSafetyBoxIssueCoin                 = "safety_box_issue_coin"
	MsgTypeSafetyBoxReturnCoin                = "safety_box_return_coin"
	MsgTypeSafetyBoxGrantOperatorPermission   = "safety_box_grant_operator_permission"
	MsgTypeSafetyBoxGrantIssuerPermission     = "safety_box_grant_issuer_permission"
	MsgTypeSafetyBoxGrantReturnerPermission   = "safety_box_grant_returner_permission"
	MsgTypeSafetyBoxGrantAllocatorPermission  = "safety_box_grant_allocator_permission"
	MsgTypeSafetyBoxRevokeOperatorPermission  = "safety_box_revoke_operator_permission"
	MsgTypeSafetyBoxRevokeIssuerPermission    = "safety_box_revoke_issuer_permission"
	MsgTypeSafetyBoxRevokeReturnerPermission  = "safety_box_revoke_returner_permission"
	MsgTypeSafetyBoxRevokeAllocatorPermission = "safety_box_revoke_allocator_permission"

	AttributeKeySafetyBoxId                        = "safety_box_id"
	AttributeKeySafetyBoxOwner                     = "safety_box_owner"
	AttributeKeySafetyBoxAddress                   = "safety_box_address"
	AttributeKeySafetyBoxAllocatorAddress          = "safety_box_allocator_address"
	AttributeKeySafetyBoxIssueFromAddress          = "safety_box_issue_from_address"
	AttributeKeySafetyBoxIssueToAddress            = "safety_box_issue_to_address"
	AttributeKeySafetyBoxReturnerAddress           = "safety_box_returner_address"
	AttributeKeySafetyBoxAction                    = "safety_box_action"
	AttributeKeySafetyBoxGrantOperatorPermission   = "safety_box_grant_operator_permission"
	AttributeKeySafetyBoxRevokeOperatorPermission  = "safety_box_revoke_operator_permission"
	AttributeKeySafetyBoxGrantAllocatorPermission  = "safety_box_grant_allocator_permission"
	AttributeKeySafetyBoxRevokeAllocatorPermission = "safety_box_revoke_allocator_permission"
	AttributeKeySafetyBoxGrantIssuerPermission     = "safety_box_grant_issuer_permission"
	AttributeKeySafetyBoxRevokeIssuerPermission    = "safety_box_revoke_issuer_permission"
	AttributeKeySafetyBoxGrantReturnerPermission   = "safety_box_grant_returner_permission"
	AttributeKeySafetyBoxRevokeReturnerPermission  = "safety_box_revoke_returner_permission"

	AttributeKeySafetyBoxCoins    = "safety_box_coins"
	AttributeKeySafetyBoxOperator = "safety_box_operator"
	AttributeKeySafetyBoxTarget   = "safety_box_target"

	AttributeValueCategory = ModuleName
)
