package collection

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterLegacyAminoCodec registers concrete types on the LegacyAmino codec
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgSendNFT{}, "lbm-sdk/MsgSendNFT")
	legacy.RegisterAminoMsg(cdc, &MsgOperatorSendNFT{}, "lbm-sdk/MsgOperatorSendNFT")
	legacy.RegisterAminoMsg(cdc, &MsgAuthorizeOperator{}, "lbm-sdk/collection/MsgAuthorizeOperator") // Changed msgName due to conflict with `x/token`
	legacy.RegisterAminoMsg(cdc, &MsgRevokeOperator{}, "lbm-sdk/collection/MsgRevokeOperator")       // Changed msgName due to conflict with `x/token`
	legacy.RegisterAminoMsg(cdc, &MsgCreateContract{}, "lbm-sdk/MsgCreateContract")
	legacy.RegisterAminoMsg(cdc, &MsgIssueNFT{}, "lbm-sdk/MsgIssueNFT")
	legacy.RegisterAminoMsg(cdc, &MsgMintNFT{}, "lbm-sdk/MsgMintNFT")
	legacy.RegisterAminoMsg(cdc, &MsgBurnNFT{}, "lbm-sdk/MsgBurnNFT")
	legacy.RegisterAminoMsg(cdc, &MsgOperatorBurnNFT{}, "lbm-sdk/MsgOperatorBurnNFT")
	legacy.RegisterAminoMsg(cdc, &MsgModify{}, "lbm-sdk/collection/MsgModify")                     // Changed msgName due to conflict with `x/token`
	legacy.RegisterAminoMsg(cdc, &MsgGrantPermission{}, "lbm-sdk/collection/MsgGrantPermission")   // Changed msgName due to conflict with `x/token`
	legacy.RegisterAminoMsg(cdc, &MsgRevokePermission{}, "lbm-sdk/collection/MsgRevokePermission") // Changed msgName due to conflict with `x/token`
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateContract{},
		&MsgIssueNFT{},
		&MsgMintNFT{},
		&MsgSendNFT{},
		&MsgOperatorSendNFT{},
		&MsgAuthorizeOperator{},
		&MsgRevokeOperator{},
		&MsgBurnNFT{},
		&MsgOperatorBurnNFT{},
		&MsgModify{},
		&MsgGrantPermission{},
		&MsgRevokePermission{},
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
		&OwnerNFT{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	cryptocodec.RegisterCrypto(amino)
	// TODO: check, Can i remove commented out code?

	//// Register all Amino interfaces and concrete types on the authz and gov Amino codec so that this can later be
	//// used to properly serialize MsgGrant, MsgExec and MsgSubmitProposal instances
	// RegisterLegacyAminoCodec(authzcodec.Amino)
	// RegisterLegacyAminoCodec(govcodec.Amino)
	// RegisterLegacyAminoCodec(fdncodec.Amino)
}
