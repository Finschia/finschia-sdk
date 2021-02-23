package types

import (
	"github.com/line/lbm-sdk/codec"
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
	cdc.RegisterConcrete(MsgIssueFT{}, "collection/MsgIssueFT", nil)
	cdc.RegisterConcrete(MsgIssueNFT{}, "collection/MsgIssueNFT", nil)
	cdc.RegisterConcrete(MsgMintNFT{}, "collection/MsgMintNFT", nil)
	cdc.RegisterConcrete(MsgBurnNFT{}, "collection/MsgBurnNFT", nil)
	cdc.RegisterConcrete(MsgBurnNFTFrom{}, "collection/MsgBurnNFTFrom", nil)
	cdc.RegisterConcrete(MsgModify{}, "collection/MsgModify", nil)
	cdc.RegisterConcrete(MsgMintFT{}, "collection/MsgMintFT", nil)
	cdc.RegisterConcrete(MsgBurnFT{}, "collection/MsgBurnFT", nil)
	cdc.RegisterConcrete(MsgBurnFTFrom{}, "collection/MsgBurnFTFrom", nil)
	cdc.RegisterConcrete(MsgGrantPermission{}, "collection/MsgGrantPermission", nil)
	cdc.RegisterConcrete(MsgRevokePermission{}, "collection/MsgRevokePermission", nil)
	cdc.RegisterConcrete(MsgTransferFT{}, "collection/MsgTransferFT", nil)
	cdc.RegisterConcrete(MsgTransferNFT{}, "collection/MsgTransferNFT", nil)
	cdc.RegisterConcrete(MsgTransferFTFrom{}, "collection/MsgTransferFTFrom", nil)
	cdc.RegisterConcrete(MsgTransferNFTFrom{}, "collection/MsgTransferNFTFrom", nil)
	cdc.RegisterConcrete(MsgAttach{}, "collection/MsgAttach", nil)
	cdc.RegisterConcrete(MsgDetach{}, "collection/MsgDetach", nil)
	cdc.RegisterConcrete(MsgAttachFrom{}, "collection/MsgAttachFrom", nil)
	cdc.RegisterConcrete(MsgDetachFrom{}, "collection/MsgDetachFrom", nil)
	cdc.RegisterConcrete(MsgApprove{}, "collection/MsgApprove", nil)
	cdc.RegisterConcrete(MsgDisapprove{}, "collection/MsgDisapprove", nil)

	cdc.RegisterInterface((*Token)(nil), nil)
	cdc.RegisterInterface((*FT)(nil), nil)

	cdc.RegisterInterface((*Collection)(nil), nil)
	cdc.RegisterConcrete(&BaseCollection{}, "collection/Collection", nil)
	cdc.RegisterConcrete(&BaseFT{}, "collection/FT", nil)
	cdc.RegisterConcrete(&BaseNFT{}, "collection/NFT", nil)

	cdc.RegisterInterface((*Account)(nil), nil)
	cdc.RegisterConcrete(&BaseAccount{}, "collection/Account", nil)

	cdc.RegisterInterface((*Supply)(nil), nil)
	cdc.RegisterConcrete(&BaseSupply{}, "collection/Supply", nil)

	cdc.RegisterInterface((*TokenType)(nil), nil)
	cdc.RegisterConcrete(&BaseTokenType{}, "collection/TokenType", nil)

	cdc.RegisterInterface((*AccountPermissionI)(nil), nil)
	cdc.RegisterConcrete(&AccountPermission{}, "collection/AccountPermission", nil)
}
