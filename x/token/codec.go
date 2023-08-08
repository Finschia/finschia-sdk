package token

import (
	"github.com/Finschia/finschia-rdk/codec"
	"github.com/Finschia/finschia-rdk/codec/legacy"
	"github.com/Finschia/finschia-rdk/codec/types"
	cryptocodec "github.com/Finschia/finschia-rdk/crypto/codec"
	sdk "github.com/Finschia/finschia-rdk/types"
	"github.com/Finschia/finschia-rdk/types/msgservice"
	authzcodec "github.com/Finschia/finschia-rdk/x/authz/codec"
	fdncodec "github.com/Finschia/finschia-rdk/x/foundation/codec"
	govcodec "github.com/Finschia/finschia-rdk/x/gov/codec"
)

// RegisterLegacyAminoCodec registers concrete types on the LegacyAmino codec
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgSend{}, "lbm-sdk/MsgSend")
	legacy.RegisterAminoMsg(cdc, &MsgOperatorSend{}, "lbm-sdk/MsgOperatorSend")
	legacy.RegisterAminoMsg(cdc, &MsgRevokeOperator{}, "lbm-sdk/token/MsgRevokeOperator")       // Changed msgName due to conflict with `x/collection`
	legacy.RegisterAminoMsg(cdc, &MsgAuthorizeOperator{}, "lbm-sdk/token/MsgAuthorizeOperator") // Changed msgName due to conflict with `x/collection`
	legacy.RegisterAminoMsg(cdc, &MsgIssue{}, "lbm-sdk/MsgIssue")
	legacy.RegisterAminoMsg(cdc, &MsgGrantPermission{}, "lbm-sdk/token/MsgGrantPermission")   // Changed msgName due to conflict with `x/collection`
	legacy.RegisterAminoMsg(cdc, &MsgRevokePermission{}, "lbm-sdk/token/MsgRevokePermission") // Changed msgName due to conflict with `x/collection`
	legacy.RegisterAminoMsg(cdc, &MsgMint{}, "lbm-sdk/MsgMint")
	legacy.RegisterAminoMsg(cdc, &MsgBurn{}, "lbm-sdk/MsgBurn")
	legacy.RegisterAminoMsg(cdc, &MsgOperatorBurn{}, "lbm-sdk/MsgOperatorBurn")
	legacy.RegisterAminoMsg(cdc, &MsgModify{}, "lbm-sdk/token/MsgModify") // Changed msgName due to conflict with `x/collection`
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSend{},
		&MsgRevokeOperator{},
		&MsgIssue{},
		&MsgMint{},
		&MsgBurn{},
		&MsgModify{},
		&MsgOperatorSend{},
		&MsgAuthorizeOperator{},
		&MsgOperatorBurn{},
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
	RegisterLegacyAminoCodec(fdncodec.Amino)
}
