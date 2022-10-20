package collection

import (
	"github.com/line/lbm-sdk/codec"
	"github.com/line/lbm-sdk/codec/legacy"
	"github.com/line/lbm-sdk/codec/types"
	cryptocodec "github.com/line/lbm-sdk/crypto/codec"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/types/msgservice"
	authzcodec "github.com/line/lbm-sdk/x/authz/codec"
	govcodec "github.com/line/lbm-sdk/x/gov/codec"
)

// RegisterLegacyAminoCodec registers concrete types on the LegacyAmino codec
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgTransferFT{}, "lbm-sdk/MsgTransferFT")
	legacy.RegisterAminoMsg(cdc, &MsgTransferFTFrom{}, "lbm-sdk/MsgTransferFTFrom")
	legacy.RegisterAminoMsg(cdc, &MsgTransferNFT{}, "lbm-sdk/MsgTransferNFT")
	legacy.RegisterAminoMsg(cdc, &MsgTransferNFTFrom{}, "lbm-sdk/MsgTransferNFTFrom")
	legacy.RegisterAminoMsg(cdc, &MsgApprove{}, "lbm-sdk/MsgCollectionApprove") // Changed msgName due to conflict with `x/token`
	legacy.RegisterAminoMsg(cdc, &MsgDisapprove{}, "lbm-sdk/MsgDisapprove")
	legacy.RegisterAminoMsg(cdc, &MsgCreateContract{}, "lbm-sdk/MsgCreateContract")
	legacy.RegisterAminoMsg(cdc, &MsgIssueFT{}, "lbm-sdk/MsgIssueFT")
	legacy.RegisterAminoMsg(cdc, &MsgIssueNFT{}, "lbm-sdk/MsgIssueNFT")
	legacy.RegisterAminoMsg(cdc, &MsgMintFT{}, "lbm-sdk/MsgMintFT")
	legacy.RegisterAminoMsg(cdc, &MsgMintNFT{}, "lbm-sdk/MsgMintNFT")
	legacy.RegisterAminoMsg(cdc, &MsgBurnFT{}, "lbm-sdk/MsgBurnFT")
	legacy.RegisterAminoMsg(cdc, &MsgBurnFTFrom{}, "lbm-sdk/MsgBurnFTFrom")
	legacy.RegisterAminoMsg(cdc, &MsgBurnNFT{}, "lbm-sdk/MsgBurnNFT")
	legacy.RegisterAminoMsg(cdc, &MsgBurnNFTFrom{}, "lbm-sdk/MsgBurnNFTFrom")
	legacy.RegisterAminoMsg(cdc, &MsgModify{}, "lbm-sdk/MsgCollectionModify")                     // Changed msgName due to conflict with `x/token`
	legacy.RegisterAminoMsg(cdc, &MsgGrantPermission{}, "lbm-sdk/MsgCollectionGrantPermission")   // Changed msgName due to conflict with `x/token`
	legacy.RegisterAminoMsg(cdc, &MsgRevokePermission{}, "lbm-sdk/MsgCollectionRevokePermission") // Changed msgName due to conflict with `x/token`
	legacy.RegisterAminoMsg(cdc, &MsgAttach{}, "lbm-sdk/MsgAttach")
	legacy.RegisterAminoMsg(cdc, &MsgDetach{}, "lbm-sdk/MsgDetach")
	legacy.RegisterAminoMsg(cdc, &MsgAttachFrom{}, "lbm-sdk/MsgAttachFrom")
	legacy.RegisterAminoMsg(cdc, &MsgDetachFrom{}, "lbm-sdk/MsgDetachFrom")

	cdc.RegisterConcrete(&MintNFTParam{}, "lbm-sdk/MintNFTParam", nil)
	cdc.RegisterConcrete(&Change{}, "lbm-sdk/Change", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateContract{},
		&MsgIssueFT{},
		&MsgIssueNFT{},
		&MsgMintFT{},
		&MsgMintNFT{},
		&MsgAttach{},
		&MsgDetach{},
		&MsgTransferFT{},
		&MsgTransferFTFrom{},
		&MsgTransferNFT{},
		&MsgTransferNFTFrom{},
		&MsgApprove{},
		&MsgDisapprove{},
		&MsgBurnFT{},
		&MsgBurnFTFrom{},
		&MsgBurnNFT{},
		&MsgBurnNFTFrom{},
		&MsgModify{},
		&MsgGrantPermission{},
		&MsgRevokePermission{},
		&MsgAttachFrom{},
		&MsgDetachFrom{},
	)

	registry.RegisterInterface(
		"lbm.collection.v1.TokenClass",
		(*TokenClass)(nil),
		&FTClass{},
		&NFTClass{},
	)

	registry.RegisterInterface(
		"lbm.collection.v1.Token",
		(*Token)(nil),
		&FT{},
		&OwnerNFT{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	sdk.RegisterLegacyAminoCodec(amino)

	// Register all Amino interfaces and concrete types on the authz and gov Amino codec so that this can later be
	// used to properly serialize MsgGrant, MsgExec and MsgSubmitProposal instances
	RegisterLegacyAminoCodec(authzcodec.Amino)
	RegisterLegacyAminoCodec(govcodec.Amino)
}
