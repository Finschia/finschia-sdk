package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

var ModuleCdc = codec.New()

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreateCollection{}, "link/MsgCreateCollection", nil)
	cdc.RegisterConcrete(MsgIssue{}, "link/MsgIssue", nil)
	cdc.RegisterConcrete(MsgIssueNFT{}, "link/MsgIssueNFT", nil)
	cdc.RegisterConcrete(MsgIssueCollection{}, "link/MsgIssueCollection", nil)
	cdc.RegisterConcrete(MsgIssueNFTCollection{}, "link/MsgIssueNFTCollection", nil)
	cdc.RegisterConcrete(MsgModifyTokenURI{}, "link/MsgModifyTokenURI", nil)
	cdc.RegisterConcrete(MsgMint{}, "link/MsgMint", nil)
	cdc.RegisterConcrete(MsgBurn{}, "link/MsgBurn", nil)
	cdc.RegisterConcrete(MsgMintCollection{}, "link/MsgMintCollection", nil)
	cdc.RegisterConcrete(MsgBurnCollection{}, "link/MsgBurnCollection", nil)
	cdc.RegisterConcrete(MsgGrantPermission{}, "link/MsgGrantPermission", nil)
	cdc.RegisterConcrete(MsgRevokePermission{}, "link/MsgRevokePermission", nil)
	cdc.RegisterConcrete(MsgTransferFT{}, "link/MsgTransferFT", nil)
	cdc.RegisterConcrete(MsgTransferCFT{}, "link/MsgTransferCFT", nil)
	cdc.RegisterConcrete(MsgTransferNFT{}, "link/MsgTransferNFT", nil)
	cdc.RegisterConcrete(MsgTransferCNFT{}, "link/MsgTransferCNFT", nil)
	cdc.RegisterConcrete(MsgAttach{}, "link/MsgAttach", nil)
	cdc.RegisterConcrete(MsgDetach{}, "link/MsgDetach", nil)
	cdc.RegisterInterface((*PermissionI)(nil), nil)
	cdc.RegisterConcrete(&Permission{}, "link/TokenPermission", nil)

	cdc.RegisterInterface((*Token)(nil), nil)
	cdc.RegisterInterface((*FT)(nil), nil)
	cdc.RegisterInterface((*NFT)(nil), nil)
	cdc.RegisterConcrete(&BaseFT{}, "link/FT", nil)
	cdc.RegisterConcrete(&BaseNFT{}, "link/NFT", nil)

	cdc.RegisterInterface((*Collection)(nil), nil)
	cdc.RegisterConcrete(&BaseCollection{}, "link/Collection", nil)
	cdc.RegisterInterface((*CollectiveToken)(nil), nil)
	cdc.RegisterConcrete(&BaseCollectiveFT{}, "link/CFT", nil)
	cdc.RegisterConcrete(&BaseCollectiveNFT{}, "link/CNFT", nil)
}
