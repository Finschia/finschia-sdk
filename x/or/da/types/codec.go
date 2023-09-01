package types

import (
	"github.com/Finschia/finschia-sdk/codec"
	"github.com/Finschia/finschia-sdk/codec/legacy"
	cdctypes "github.com/Finschia/finschia-sdk/codec/types"
	cryptocodec "github.com/Finschia/finschia-sdk/crypto/codec"
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/types/msgservice"
	authzcodec "github.com/Finschia/finschia-sdk/x/authz/codec"
	fdncodec "github.com/Finschia/finschia-sdk/x/foundation/codec"
	govcodec "github.com/Finschia/finschia-sdk/x/gov/codec"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "finschia-sdk/or/da/MsgUpdateParams")
	legacy.RegisterAminoMsg(cdc, &MsgEnqueue{}, "finshcia-sdk/or/MsgEnqueue")
	legacy.RegisterAminoMsg(cdc, &MsgAppendCCBatch{}, "finschia-sdk/or/MsgAppendCCBatch")
	legacy.RegisterAminoMsg(cdc, &MsgAppendSCCBatch{}, "finshcia-sdk/or/MsgAppendSCCBatch")

	cdc.RegisterConcrete(Params{}, "finschia-sdk/x/or/da/Params", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgUpdateParams{},
		&MsgEnqueue{},
		&MsgAppendCCBatch{},
		&MsgAppendSCCBatch{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)

func init() {
	RegisterLegacyAminoCodec(Amino)
	cryptocodec.RegisterCrypto(Amino)
	sdk.RegisterLegacyAminoCodec(Amino)

	// Register all Amino interfaces and concrete types on the authz and gov Amino codec so that this can later be
	// used to properly serialize MsgGrant, MsgExec and MsgSubmitProposal instances
	RegisterLegacyAminoCodec(authzcodec.Amino)
	RegisterLegacyAminoCodec(govcodec.Amino)
	RegisterLegacyAminoCodec(fdncodec.Amino)
}
