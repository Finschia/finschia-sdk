package types

var (
	EventSafetyBoxCreate     = "safety_box_create"
	EventSafetyBoxSendCoin   = "safety_box_send_coin"
	EventSafetyBoxPermission = "safety_box_permission"

	MsgTypeSafetyBoxCreate                   = "safety_box_create"
	MsgTypeSafetyBoxAllocateCoin             = "safety_box_allocate_coin"
	MsgTypeSafetyBoxRecallCoin               = "safety_box_recall_coin"
	MsgTypeSafetyBoxIssueCoin                = "safety_box_issue_coin"
	MsgTypeSafetyBoxReturnCoin               = "safety_box_return_coin"
	MsgTypeSafetyBoxGrantIssuePermission     = "safety_box_grant_issue_permission"
	MsgTypeSafetyBoxGrantReturnPermission    = "safety_box_grant_return_permission"
	MsgTypeSafetyBoxGrantAllocatePermission  = "safety_box_grant_allocate_permission"
	MsgTypeSafetyBoxGrantRecallPermission    = "safety_box_grant_recall_permission"
	MsgTypeSafetyBoxRevokeIssuePermission    = "safety_box_revoke_issue_permission"
	MsgTypeSafetyBoxRevokeReturnPermission   = "safety_box_revoke_return_permission"
	MsgTypeSafetyBoxRevokeAllocatePermission = "safety_box_revoke_allocate_permission"
	MsgTypeSafetyBoxRevokeRecallPermission   = "safety_box_revoke_recall_permission"

	AttributeKeySafetyBoxId                       = "safety_box_id"
	AttributeKeySafetyBoxOwner                    = "safety_box_owner"
	AttributeKeySafetyBoxAddress                  = "safety_box_address"
	AttributeKeySafetyBoxAllocatorAddress         = "safety_box_allocator_address"
	AttributeKeySafetyBoxIssueFromAddress         = "safety_box_issue_from_address"
	AttributeKeySafetyBoxIssueToAddress           = "safety_box_issue_to_address"
	AttributeKeySafetyBoxReturnerAddress          = "safety_box_returner_address"
	AttributeKeySafetyBoxAction                   = "safety_box_action"
	AttributeKeySafetyBoxGrantAllocatePermission  = "safety_box_grant_allocate_permission"
	AttributeKeySafetyBoxRevokeAllocatePermission = "safety_box_revoke_allocate_permission"
	AttributeKeySafetyBoxGrantRecallPermission    = "safety_box_grant_recall_permission"
	AttributeKeySafetyBoxRevokeRecallPermission   = "safety_box_revoke_recall_permission"
	AttributeKeySafetyBoxGrantIssuePermission     = "safety_box_grant_issue_permission"
	AttributeKeySafetyBoxRevokeIssuePermission    = "safety_box_revoke_issue_permission"
	AttributeKeySafetyBoxGrantReturnPermission    = "safety_box_grant_return_permission"
	AttributeKeySafetyBoxRevokeReturnPermission   = "safety_box_revoke_return_permission"

	AttributeKeySafetyBoxCoins = "safety_box_coins"

	AttributeKeySafetyBoxOperator = "safety_box_operator"
	AttributeKeySafetyBoxTarget   = "safety_box_target"

	AttributeValueCategory = ModuleName
)
