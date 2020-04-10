package safetybox

import (
	"github.com/line/link/x/safetybox/internal/keeper"
	"github.com/line/link/x/safetybox/internal/types"
)

const (
	ModuleName = types.ModuleName
	StoreKey   = types.StoreKey
	RouterKey  = types.RouterKey
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
	ModuleCdc              = types.ModuleCdc
	RegisterCodec          = types.RegisterCodec
	NewKeeper              = keeper.NewKeeper
	NewMultiSafetyBoxHooks = types.NewMultiSafetyBoxHooks

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
	ErrSafetyBoxIncorrectDenom     = types.ErrSafetyBoxIncorrectDenom
	ErrSafetyBoxTooManyCoinDenoms  = types.ErrSafetyBoxTooManyCoinDenoms

	AttributeKeySafetyBoxID                        = types.AttributeKeySafetyBoxID
	AttributeKeySafetyBoxOwner                     = types.AttributeKeySafetyBoxOwner
	AttributeKeySafetyBoxAddress                   = types.AttributeKeySafetyBoxAddress
	AttributeKeySafetyBoxAllocatorAddress          = types.AttributeKeySafetyBoxAllocatorAddress
	AttributeKeySafetyBoxIssueFromAddress          = types.AttributeKeySafetyBoxIssueFromAddress
	AttributeKeySafetyBoxIssueToAddress            = types.AttributeKeySafetyBoxIssueToAddress
	AttributeKeySafetyBoxReturnerAddress           = types.AttributeKeySafetyBoxReturnerAddress
	AttributeKeySafetyBoxAction                    = types.AttributeKeySafetyBoxAction
	AttributeKeySafetyBoxGrantOperatorPermission   = types.AttributeKeySafetyBoxGrantOperatorPermission
	AttributeKeySafetyBoxRevokeOperatorPermission  = types.AttributeKeySafetyBoxRevokeOperatorPermission
	AttributeKeySafetyBoxGrantAllocatorPermission  = types.AttributeKeySafetyBoxGrantAllocatorPermission
	AttributeKeySafetyBoxRevokeAllocatorPermission = types.AttributeKeySafetyBoxRevokeAllocatorPermission
	AttributeKeySafetyBoxGrantIssuerPermission     = types.AttributeKeySafetyBoxGrantIssuerPermission
	AttributeKeySafetyBoxRevokeIssuerPermission    = types.AttributeKeySafetyBoxRevokeIssuerPermission
	AttributeKeySafetyBoxGrantReturnerPermission   = types.AttributeKeySafetyBoxGrantReturnerPermission
	AttributeKeySafetyBoxRevokeReturnerPermission  = types.AttributeKeySafetyBoxRevokeReturnerPermission
	AttributeKeySafetyBoxCoins                     = types.AttributeKeySafetyBoxCoins
	AttributeValueCategory                         = types.AttributeValueCategory

	AttributeKeySafetyBoxOperator = types.AttributeKeySafetyBoxOperator
	AttributeKeySafetyBoxTarget   = types.AttributeKeySafetyBoxTarget
)
