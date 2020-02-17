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
	cdc.RegisterConcrete(MsgCreateCollection{}, "link/MsgCreateCollection", nil)
	cdc.RegisterConcrete(MsgIssue{}, "link/MsgIssue", nil)
	cdc.RegisterConcrete(MsgIssueCFT{}, "link/MsgIssueCFT", nil)
	cdc.RegisterConcrete(MsgIssueCNFT{}, "link/MsgIssueCNFT", nil)
	cdc.RegisterConcrete(MsgMintCNFT{}, "link/MsgMintCNFT", nil)
	cdc.RegisterConcrete(MsgBurnCNFT{}, "link/MsgBurnCNFT", nil)
	cdc.RegisterConcrete(MsgBurnCNFTFrom{}, "link/MsgBurnCNFTFrom", nil)
	cdc.RegisterConcrete(MsgModifyTokenURI{}, "link/MsgModifyTokenURI", nil)
	cdc.RegisterConcrete(MsgMint{}, "link/MsgMint", nil)
	cdc.RegisterConcrete(MsgBurn{}, "link/MsgBurn", nil)
	cdc.RegisterConcrete(MsgMintCFT{}, "link/MsgMintCFT", nil)
	cdc.RegisterConcrete(MsgBurnCFT{}, "link/MsgBurnCFT", nil)
	cdc.RegisterConcrete(MsgBurnCFTFrom{}, "link/MsgBurnCFTFrom", nil)
	cdc.RegisterConcrete(MsgGrantPermission{}, "link/MsgGrantPermission", nil)
	cdc.RegisterConcrete(MsgRevokePermission{}, "link/MsgRevokePermission", nil)
	cdc.RegisterConcrete(MsgTransferFT{}, "link/MsgTransferFT", nil)
	cdc.RegisterConcrete(MsgTransferCFT{}, "link/MsgTransferCFT", nil)
	cdc.RegisterConcrete(MsgTransferCNFT{}, "link/MsgTransferCNFT", nil)
	cdc.RegisterConcrete(MsgTransferCFTFrom{}, "link/MsgTransferCFTFrom", nil)
	cdc.RegisterConcrete(MsgTransferCNFTFrom{}, "link/MsgTransferCNFTFrom", nil)
	cdc.RegisterConcrete(MsgAttach{}, "link/MsgAttach", nil)
	cdc.RegisterConcrete(MsgDetach{}, "link/MsgDetach", nil)
	cdc.RegisterConcrete(MsgAttachFrom{}, "link/MsgAttachFrom", nil)
	cdc.RegisterConcrete(MsgDetachFrom{}, "link/MsgDetachFrom", nil)
	cdc.RegisterConcrete(MsgApproveCollection{}, "link/MsgApproveCollection", nil)
	cdc.RegisterConcrete(MsgDisapproveCollection{}, "link/MsgDisapproveCollection", nil)
	cdc.RegisterInterface((*PermissionI)(nil), nil)
	cdc.RegisterConcrete(&Permission{}, "link/TokenPermission", nil)

	cdc.RegisterInterface((*Token)(nil), nil)
	cdc.RegisterInterface((*FT)(nil), nil)
	cdc.RegisterConcrete(&BaseFT{}, "link/FT", nil)

	cdc.RegisterInterface((*Collection)(nil), nil)
	cdc.RegisterConcrete(&BaseCollection{}, "link/Collection", nil)
	cdc.RegisterInterface((*CollectiveToken)(nil), nil)
	cdc.RegisterConcrete(&BaseCollectiveFT{}, "link/CFT", nil)
	cdc.RegisterConcrete(&BaseCollectiveNFT{}, "link/CNFT", nil)
}
