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
	cdc.RegisterConcrete(MsgCreateCollection{}, "collection/MsgCreate", nil)
	cdc.RegisterConcrete(MsgIssueCFT{}, "collection/MsgIssueFT", nil)
	cdc.RegisterConcrete(MsgIssueCNFT{}, "collection/MsgIssueNFT", nil)
	cdc.RegisterConcrete(MsgMintCNFT{}, "collection/MsgMintNFT", nil)
	cdc.RegisterConcrete(MsgBurnCNFT{}, "collection/MsgBurnNFT", nil)
	cdc.RegisterConcrete(MsgBurnCNFTFrom{}, "collection/MsgBurnNFTFrom", nil)
	cdc.RegisterConcrete(MsgModifyTokenURI{}, "collection/MsgModifyTokenURI", nil)
	cdc.RegisterConcrete(MsgMintCFT{}, "collection/MsgMintFT", nil)
	cdc.RegisterConcrete(MsgBurnCFT{}, "collection/MsgBurnFT", nil)
	cdc.RegisterConcrete(MsgBurnCFTFrom{}, "collection/MsgBurnFTFrom", nil)
	cdc.RegisterConcrete(MsgGrantPermission{}, "collection/MsgGrantPermission", nil)
	cdc.RegisterConcrete(MsgRevokePermission{}, "collection/MsgRevokePermission", nil)
	cdc.RegisterConcrete(MsgTransferCFT{}, "collection/MsgTransferFT", nil)
	cdc.RegisterConcrete(MsgTransferCNFT{}, "collection/MsgTransferNFT", nil)
	cdc.RegisterConcrete(MsgTransferCFTFrom{}, "collection/MsgTransferFTFrom", nil)
	cdc.RegisterConcrete(MsgTransferCNFTFrom{}, "collection/MsgTransferNFTFrom", nil)
	cdc.RegisterConcrete(MsgAttach{}, "collection/MsgAttach", nil)
	cdc.RegisterConcrete(MsgDetach{}, "collection/MsgDetach", nil)
	cdc.RegisterConcrete(MsgAttachFrom{}, "collection/MsgAttachFrom", nil)
	cdc.RegisterConcrete(MsgDetachFrom{}, "collection/MsgDetachFrom", nil)
	cdc.RegisterConcrete(MsgApprove{}, "collection/MsgApprove", nil)
	cdc.RegisterConcrete(MsgDisapprove{}, "collection/MsgDisapprove", nil)
	cdc.RegisterInterface((*PermissionI)(nil), nil)
	cdc.RegisterConcrete(&Permission{}, "collection/Permission", nil)

	cdc.RegisterInterface((*Token)(nil), nil)
	cdc.RegisterInterface((*FT)(nil), nil)

	cdc.RegisterInterface((*Collection)(nil), nil)
	cdc.RegisterConcrete(&BaseCollection{}, "collection/Collection", nil)
	cdc.RegisterConcrete(&BaseFT{}, "collection/FT", nil)
	cdc.RegisterConcrete(&BaseNFT{}, "collection/NFT", nil)
}
