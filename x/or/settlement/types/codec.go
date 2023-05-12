package types

import (
	"github.com/Finschia/finschia-sdk/codec"
	"github.com/Finschia/finschia-sdk/codec/legacy"
	cdctypes "github.com/Finschia/finschia-sdk/codec/types"

	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgAddTrieNode{}, "finschia-sdk/MsgAddTrieNode")
	legacy.RegisterAminoMsg(cdc, &MsgConfirmStateTransition{}, "finschia-sdk/MsgConfirmStateTransition")
	legacy.RegisterAminoMsg(cdc, &MsgDenyStateTransition{}, "finschia-sdk/MsgDenyStateTransition")
	legacy.RegisterAminoMsg(cdc, &MsgInitiateChallenge{}, "finschia-sdk/MsgInitiateChallenge")
	legacy.RegisterAminoMsg(cdc, &MsgProposeState{}, "finschia-sdk/MsgProposeState")
	legacy.RegisterAminoMsg(cdc, &MsgRespondState{}, "finschia-sdk/MsgRespondState")
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgAddTrieNode{},
		&MsgConfirmStateTransition{},
		&MsgDenyStateTransition{},
		&MsgInitiateChallenge{},
		&MsgProposeState{},
		&MsgRespondState{},
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
