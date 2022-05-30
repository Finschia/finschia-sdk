package collection

import (
	"github.com/line/lbm-sdk/codec/types"
	// sdk "github.com/line/lbm-sdk/types"
	// "github.com/line/lbm-sdk/types/msgservice"
)

func RegisterInterfaces(registry types.InterfaceRegistry) {
	// registry.RegisterImplementations((*sdk.Msg)(nil))

	registry.RegisterInterface(
		"lbm.collection.v1.TokenClass",
		(*TokenClass)(nil),
		&FTClass{},
		&NFTClass{},
	)

	// msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
