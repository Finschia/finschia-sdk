package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

var ModuleCdc = codec.New()

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgIssue{}, "link/MsgIssue", nil)
	cdc.RegisterConcrete(MsgIssueNFT{}, "link/MsgIssueNFT", nil)
	cdc.RegisterConcrete(MsgIssueCollection{}, "link/MsgIssueCollection", nil)
	cdc.RegisterConcrete(MsgIssueNFTCollection{}, "link/MsgIssueNFTCollection", nil)
	cdc.RegisterConcrete(MsgMint{}, "link/MsgMint", nil)
	cdc.RegisterConcrete(MsgBurn{}, "link/MsgBurn", nil)
	cdc.RegisterConcrete(MsgGrantPermission{}, "link/MsgGrantPermission", nil)
	cdc.RegisterConcrete(MsgRevokePermission{}, "link/MsgRevokePermission", nil)
	cdc.RegisterInterface((*PermissionI)(nil), nil)
	cdc.RegisterConcrete(&Permission{}, "link/TokenPermission", nil)
}

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}
