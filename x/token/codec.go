package token

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
	legacy.RegisterAminoMsg(cdc, &MsgSend{}, "lbm-sdk/MsgSend")
	legacy.RegisterAminoMsg(cdc, &MsgTransferFrom{}, "lbm-sdk/MsgTransferFrom")
	legacy.RegisterAminoMsg(cdc, &MsgRevokeOperator{}, "lbm-sdk/MsgRevokeOperator")
	legacy.RegisterAminoMsg(cdc, &MsgApprove{}, "lbm-sdk/MsgTokenApprove") // Changed msgName due to conflict with `x/collection`
	legacy.RegisterAminoMsg(cdc, &MsgIssue{}, "lbm-sdk/MsgIssue")
	legacy.RegisterAminoMsg(cdc, &MsgGrantPermission{}, "lbm-sdk/MsgTokenGrantPermission")   // Changed msgName due to conflict with `x/collection`
	legacy.RegisterAminoMsg(cdc, &MsgRevokePermission{}, "lbm-sdk/MsgTokenRevokePermission") // Changed msgName due to conflict with `x/collection`
	legacy.RegisterAminoMsg(cdc, &MsgMint{}, "lbm-sdk/MsgMint")
	legacy.RegisterAminoMsg(cdc, &MsgBurn{}, "lbm-sdk/MsgBurn")
	legacy.RegisterAminoMsg(cdc, &MsgBurnFrom{}, "lbm-sdk/MsgBurnFrom")
	legacy.RegisterAminoMsg(cdc, &MsgModify{}, "lbm-sdk/MsgTokenModify") // Changed msgName due to conflict with `x/collection`

	cdc.RegisterConcrete(&Pair{}, "lbm-sdk/Pair", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSend{},
		&MsgRevokeOperator{},
		&MsgIssue{},
		&MsgMint{},
		&MsgBurn{},
		&MsgModify{},
		&MsgTransferFrom{},
		&MsgApprove{},
		&MsgBurnFrom{},
		&MsgGrantPermission{},
		&MsgRevokePermission{},
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
