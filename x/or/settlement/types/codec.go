package types

import (
	"github.com/Finschia/finschia-sdk/codec"
	"github.com/Finschia/finschia-sdk/codec/legacy"
	cdctypes "github.com/Finschia/finschia-sdk/codec/types"

	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgStartChallenge{}, "finschia-sdk/MsgStartChallenge")
	legacy.RegisterAminoMsg(cdc, &MsgNsectChallenge{}, "finschia-sdk/MsgNsectChallenge")
	legacy.RegisterAminoMsg(cdc, &MsgFinishChallenge{}, "finschia-sdk/MsgFinishChallenge")
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgStartChallenge{},
		&MsgNsectChallenge{},
		&MsgFinishChallenge{},
	)

	registry.RegisterInterface(
		"finschia.or.settlement.v1.Challenge",
		(*interface{})(nil),
		&Challenge{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
