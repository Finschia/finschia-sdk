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
	cdc.RegisterConcrete(MsgIssue{}, "token/MsgIssue", nil)
	cdc.RegisterConcrete(MsgModify{}, "token/MsgModify", nil)
	cdc.RegisterConcrete(MsgMint{}, "token/MsgMint", nil)
	cdc.RegisterConcrete(MsgBurn{}, "token/MsgBurn", nil)
	cdc.RegisterConcrete(MsgGrantPermission{}, "token/MsgGrantPermission", nil)
	cdc.RegisterConcrete(MsgRevokePermission{}, "token/MsgRevokePermission", nil)
	cdc.RegisterConcrete(MsgTransfer{}, "token/MsgTransfer", nil)
	cdc.RegisterInterface((*PermissionI)(nil), nil)
	cdc.RegisterConcrete(&Permission{}, "token/TokenPermission", nil)
	cdc.RegisterInterface((*Token)(nil), nil)
	cdc.RegisterConcrete(&BaseToken{}, "token/Token", nil)
	cdc.RegisterInterface((*Supply)(nil), nil)
	cdc.RegisterConcrete(&BaseSupply{}, "token/Supply", nil)
	cdc.RegisterInterface((*Account)(nil), nil)
	cdc.RegisterConcrete(&BaseAccount{}, "token/Account", nil)
}
