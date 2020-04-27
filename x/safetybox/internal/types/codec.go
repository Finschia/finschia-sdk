package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(SafetyBox{}, "safetybox", nil)
	cdc.RegisterConcrete(MsgSafetyBoxCreate{}, "safetybox/MsgCreate", nil)
	cdc.RegisterConcrete(MsgSafetyBoxAllocateToken{}, "safetybox/MsgAllocate", nil)
	cdc.RegisterConcrete(MsgSafetyBoxRecallToken{}, "safetybox/MsgRecall", nil)
	cdc.RegisterConcrete(MsgSafetyBoxIssueToken{}, "safetybox/MsgIssue", nil)
	cdc.RegisterConcrete(MsgSafetyBoxReturnToken{}, "safetybox/MsgReturn", nil)
	cdc.RegisterConcrete(MsgSafetyBoxRegisterIssuer{}, "safetybox/MsgGrantIssuerPermission", nil)
	cdc.RegisterConcrete(MsgSafetyBoxDeregisterIssuer{}, "safetybox/MsgRevokeIssuerPermission", nil)
	cdc.RegisterConcrete(MsgSafetyBoxRegisterReturner{}, "safetybox/MsgGrantReturnerPermission", nil)
	cdc.RegisterConcrete(MsgSafetyBoxDeregisterReturner{}, "safetybox/MsgRevokeReturnerPermission", nil)
	cdc.RegisterConcrete(MsgSafetyBoxRegisterAllocator{}, "safetybox/MsgGrantAllocatorPermission", nil)
	cdc.RegisterConcrete(MsgSafetyBoxDeregisterAllocator{}, "safetybox/MsgRevokeAllocatorPermission", nil)
	cdc.RegisterConcrete(MsgSafetyBoxRegisterOperator{}, "safetybox/MsgGrantOperatorPermission", nil)
	cdc.RegisterConcrete(MsgSafetyBoxDeregisterOperator{}, "safetybox/MsgRevokeOperatorPermission", nil)
	cdc.RegisterInterface((*PermissionI)(nil), nil)
	cdc.RegisterConcrete(&Permission{}, "safetybox/perms", nil)
}
