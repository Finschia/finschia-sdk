package types

import (
	"github.com/Finschia/finschia-sdk/codec"
	cdctypes "github.com/Finschia/finschia-sdk/codec/types"
	// this line is used by starport scaffolding # 1
	"github.com/Finschia/finschia-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
