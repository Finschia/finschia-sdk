package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

var ModuleCdc = codec.New()

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(SafetyBox{}, "safetybox", nil)
	cdc.RegisterConcrete(MsgSafetyBoxCreate{}, "safetybox/MsgCreate", nil)
	cdc.RegisterConcrete(MsgSafetyBoxAllocateCoins{}, "safetybox/MsgAllocate", nil)
	cdc.RegisterConcrete(MsgSafetyBoxRecallCoins{}, "safetybox/MsgRecall", nil)
	cdc.RegisterConcrete(MsgSafetyBoxIssueCoins{}, "safetybox/MsgIssue", nil)
	cdc.RegisterConcrete(MsgSafetyBoxReturnCoins{}, "safetybox/MsgReturn", nil)
	cdc.RegisterConcrete(MsgSafetyBoxRegisterIssuer{}, "safetybox/MsgGrantIssuePermission", nil)
	cdc.RegisterConcrete(MsgSafetyBoxDeregisterIssuer{}, "safetybox/MsgRevokeIssuePermission", nil)
	cdc.RegisterConcrete(MsgSafetyBoxRegisterReturner{}, "safetybox/MsgGrantReturnPermission", nil)
	cdc.RegisterConcrete(MsgSafetyBoxDeregisterReturner{}, "safetybox/MsgRevokeReturnPermission", nil)
	cdc.RegisterConcrete(MsgSafetyBoxRegisterAllocator{}, "safetybox/MsgGrantAllocatePermission", nil)
	cdc.RegisterConcrete(MsgSafetyBoxDeregisterAllocator{}, "safetybox/MsgRevokeAllocatePermission", nil)
	cdc.RegisterConcrete(MsgSafetyBoxRegisterOperator{}, "safetybox/MsgGrantRecallPermission", nil)
	cdc.RegisterConcrete(MsgSafetyBoxDeregisterOperator{}, "safetybox/MsgRevokeRecallPermission", nil)
	cdc.RegisterInterface((*PermissionI)(nil), nil)
	cdc.RegisterConcrete(&Permission{}, "safetybox/perms", nil)
}

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}
