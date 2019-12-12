package safetybox

import (
	"github.com/line/link/x/safetybox/internal/keeper"
	"github.com/line/link/x/safetybox/internal/types"
)

const (
	ModuleName       = types.ModuleName
	StoreKey         = types.StoreKey
	RouterKey        = types.RouterKey
	DefaultCodespace = types.DefaultCodespace
)

type (
	Keeper    = keeper.Keeper
	SafetyBox = types.SafetyBox

	MsgSafetyBoxCreate              = types.MsgSafetyBoxCreate
	MsgSafetyBoxAllocateCoins       = types.MsgSafetyBoxAllocateCoins
	MsgSafetyBoxRecallCoins         = types.MsgSafetyBoxRecallCoins
	MsgSafetyBoxIssueCoins          = types.MsgSafetyBoxIssueCoins
	MsgSafetyBoxReturnCoins         = types.MsgSafetyBoxReturnCoins
	MsgSafetyBoxRegisterIssuer      = types.MsgSafetyBoxRegisterIssuer
	MsgSafetyBoxDeregisterIssuer    = types.MsgSafetyBoxDeregisterIssuer
	MsgSafetyBoxRegisterReturner    = types.MsgSafetyBoxRegisterReturner
	MsgSafetyBoxDeregisterReturner  = types.MsgSafetyBoxDeregisterReturner
	MsgSafetyBoxRegisterAllocator   = types.MsgSafetyBoxRegisterAllocator
	MsgSafetyBoxDeregisterAllocator = types.MsgSafetyBoxDeregisterAllocator
	MsgSafetyBoxRegisterOperator    = types.MsgSafetyBoxRegisterOperator
	MsgSafetyBoxDeregisterOperator  = types.MsgSafetyBoxDeregisterOperator
	MsgSafetyBoxRoleResponse        = types.MsgSafetyBoxRoleResponse
)

var (
	ModuleCdc     = types.ModuleCdc
	RegisterCodec = types.RegisterCodec
	NewKeeper     = keeper.NewKeeper

	RoleOwner     = types.RoleOwner
	RoleOperator  = types.RoleOperator
	RoleAllocator = types.RoleAllocator
	RoleIssuer    = types.RoleIssuer
	RoleReturner  = types.RoleReturner

	ActionAllocate = types.ActionAllocate
	ActionRecall   = types.ActionRecall
	ActionIssue    = types.ActionIssue
	ActionReturn   = types.ActionReturn

	RegisterRole   = types.RegisterRole
	DeregisterRole = types.DeregisterRole

	EventSafetyBoxCreate     = types.EventSafetyBoxCreate
	EventSafetyBoxSendCoin   = types.EventSafetyBoxSendCoin
	EventSafetyBoxPermission = types.EventSafetyBoxPermission

	ErrSafetyBoxInvalidMsgType     = types.ErrSafetyBoxInvalidMsgType
	ErrSafetyBoxPermissionAllocate = types.ErrSafetyBoxPermissionAllocate
	ErrSafetyBoxPermissionRecall   = types.ErrSafetyBoxPermissionRecall
	ErrSafetyBoxPermissionIssue    = types.ErrSafetyBoxPermissionIssue
	ErrSafetyBoxPermissionReturn   = types.ErrSafetyBoxPermissionReturn

	AttributeKeySafetyBoxId                       = types.AttributeKeySafetyBoxId
	AttributeKeySafetyBoxOwner                    = types.AttributeKeySafetyBoxOwner
	AttributeKeySafetyBoxAddress                  = types.AttributeKeySafetyBoxAddress
	AttributeKeySafetyBoxAllocatorAddress         = types.AttributeKeySafetyBoxAllocatorAddress
	AttributeKeySafetyBoxIssueFromAddress         = types.AttributeKeySafetyBoxIssueFromAddress
	AttributeKeySafetyBoxIssueToAddress           = types.AttributeKeySafetyBoxIssueToAddress
	AttributeKeySafetyBoxReturnerAddress          = types.AttributeKeySafetyBoxReturnerAddress
	AttributeKeySafetyBoxAction                   = types.AttributeKeySafetyBoxAction
	AttributeKeySafetyBoxGrantAllocatePermission  = types.AttributeKeySafetyBoxGrantAllocatePermission
	AttributeKeySafetyBoxRevokeAllocatePermission = types.AttributeKeySafetyBoxRevokeAllocatePermission
	AttributeKeySafetyBoxGrantRecallPermission    = types.AttributeKeySafetyBoxGrantRecallPermission
	AttributeKeySafetyBoxRevokeRecallPermission   = types.AttributeKeySafetyBoxRevokeRecallPermission
	AttributeKeySafetyBoxGrantIssuePermission     = types.AttributeKeySafetyBoxGrantIssuePermission
	AttributeKeySafetyBoxRevokeIssuePermission    = types.AttributeKeySafetyBoxRevokeIssuePermission
	AttributeKeySafetyBoxGrantReturnPermission    = types.AttributeKeySafetyBoxGrantReturnPermission
	AttributeKeySafetyBoxRevokeReturnPermission   = types.AttributeKeySafetyBoxRevokeReturnPermission
	AttributeKeySafetyBoxCoins                    = types.AttributeKeySafetyBoxCoins
	AttributeValueCategory                        = types.AttributeValueCategory

	AttributeKeySafetyBoxOperator = types.AttributeKeySafetyBoxOperator
	AttributeKeySafetyBoxTarget   = types.AttributeKeySafetyBoxTarget
)
