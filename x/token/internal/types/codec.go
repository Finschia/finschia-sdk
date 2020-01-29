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
	cdc.RegisterConcrete(MsgModifyTokenURI{}, "link/MsgModifyTokenURI", nil)
	cdc.RegisterConcrete(MsgMint{}, "link/MsgMint", nil)
	cdc.RegisterConcrete(MsgBurn{}, "link/MsgBurn", nil)
	cdc.RegisterConcrete(MsgGrantPermission{}, "link/MsgGrantPermission", nil)
	cdc.RegisterConcrete(MsgRevokePermission{}, "link/MsgRevokePermission", nil)
	cdc.RegisterConcrete(MsgTransferFT{}, "link/MsgTransferFT", nil)
	cdc.RegisterConcrete(MsgTransferIDFT{}, "link/MsgTransferIDFT", nil)
	cdc.RegisterConcrete(MsgTransferNFT{}, "link/MsgTransferNFT", nil)
	cdc.RegisterConcrete(MsgTransferIDNFT{}, "link/MsgTransferIDNFT", nil)
	cdc.RegisterConcrete(MsgAttach{}, "link/MsgAttach", nil)
	cdc.RegisterConcrete(MsgDetach{}, "link/MsgDetach", nil)
	cdc.RegisterInterface((*PermissionI)(nil), nil)
	cdc.RegisterConcrete(&Permission{}, "link/TokenPermission", nil)

	cdc.RegisterInterface((*Token)(nil), nil)
	cdc.RegisterInterface((*FT)(nil), nil)
	cdc.RegisterInterface((*NFT)(nil), nil)
	cdc.RegisterConcrete(&BaseFT{}, "link/FT", nil)
	cdc.RegisterConcrete(&BaseNFT{}, "link/NFT", nil)
	cdc.RegisterConcrete(&BaseIDFT{}, "link/IDFT", nil)
	cdc.RegisterConcrete(&BaseIDNFT{}, "link/IDNFT", nil)
}

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}
