package collection

import (
	"github.com/Finschia/finschia-sdk/codec"
	"github.com/Finschia/finschia-sdk/codec/legacy"
	"github.com/Finschia/finschia-sdk/codec/types"
	cryptocodec "github.com/Finschia/finschia-sdk/crypto/codec"
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/types/msgservice"
	authzcodec "github.com/Finschia/finschia-sdk/x/authz/codec"
	fdncodec "github.com/Finschia/finschia-sdk/x/foundation/codec"
	govcodec "github.com/Finschia/finschia-sdk/x/gov/codec"
)

// RegisterLegacyAminoCodec registers concrete types on the LegacyAmino codec
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgSendFT{}, "lbm-sdk/MsgSendFT")
	legacy.RegisterAminoMsg(cdc, &MsgOperatorSendFT{}, "lbm-sdk/MsgOperatorSendFT")
	legacy.RegisterAminoMsg(cdc, &MsgSendNFT{}, "lbm-sdk/MsgSendNFT")
	legacy.RegisterAminoMsg(cdc, &MsgOperatorSendNFT{}, "lbm-sdk/MsgOperatorSendNFT")
	legacy.RegisterAminoMsg(cdc, &MsgAuthorizeOperator{}, "lbm-sdk/collection/MsgAuthorizeOperator") // Changed msgName due to conflict with `x/token`
	legacy.RegisterAminoMsg(cdc, &MsgRevokeOperator{}, "lbm-sdk/collection/MsgRevokeOperator")       // Changed msgName due to conflict with `x/token`
	legacy.RegisterAminoMsg(cdc, &MsgCreateContract{}, "lbm-sdk/MsgCreateContract")
	legacy.RegisterAminoMsg(cdc, &MsgIssueFT{}, "lbm-sdk/MsgIssueFT")
	legacy.RegisterAminoMsg(cdc, &MsgIssueNFT{}, "lbm-sdk/MsgIssueNFT")
	legacy.RegisterAminoMsg(cdc, &MsgMintFT{}, "lbm-sdk/MsgMintFT")
	legacy.RegisterAminoMsg(cdc, &MsgMintNFT{}, "lbm-sdk/MsgMintNFT")
	legacy.RegisterAminoMsg(cdc, &MsgBurnFT{}, "lbm-sdk/MsgBurnFT")
	legacy.RegisterAminoMsg(cdc, &MsgOperatorBurnFT{}, "lbm-sdk/MsgOperatorBurnFT")
	legacy.RegisterAminoMsg(cdc, &MsgBurnNFT{}, "lbm-sdk/MsgBurnNFT")
	legacy.RegisterAminoMsg(cdc, &MsgOperatorBurnNFT{}, "lbm-sdk/MsgOperatorBurnNFT")
	legacy.RegisterAminoMsg(cdc, &MsgModify{}, "lbm-sdk/collection/MsgModify")                     // Changed msgName due to conflict with `x/token`
	legacy.RegisterAminoMsg(cdc, &MsgGrantPermission{}, "lbm-sdk/collection/MsgGrantPermission")   // Changed msgName due to conflict with `x/token`
	legacy.RegisterAminoMsg(cdc, &MsgRevokePermission{}, "lbm-sdk/collection/MsgRevokePermission") // Changed msgName due to conflict with `x/token`
	legacy.RegisterAminoMsg(cdc, &MsgAttach{}, "lbm-sdk/MsgAttach")
	legacy.RegisterAminoMsg(cdc, &MsgDetach{}, "lbm-sdk/MsgDetach")
	legacy.RegisterAminoMsg(cdc, &MsgOperatorAttach{}, "lbm-sdk/MsgOperatorAttach")
	legacy.RegisterAminoMsg(cdc, &MsgOperatorDetach{}, "lbm-sdk/MsgOperatorDetach")
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
		&MsgSendFT{},
		&MsgOperatorSendFT{},
		&MsgSendNFT{},
		&MsgOperatorSendNFT{},
		&MsgAuthorizeOperator{},
		&MsgRevokeOperator{},
		&MsgBurnFT{},
		&MsgOperatorBurnFT{},
		&MsgBurnNFT{},
		&MsgOperatorBurnNFT{},
		&MsgModify{},
		&MsgGrantPermission{},
		&MsgRevokePermission{},
		&MsgOperatorAttach{},
		&MsgOperatorDetach{},
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
	RegisterLegacyAminoCodec(fdncodec.Amino)
}
