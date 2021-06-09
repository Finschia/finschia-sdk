package types

import (
	"github.com/line/lfb-sdk/codec"
	"github.com/line/lfb-sdk/codec/types"
)

// RegisterLegacyAminoCodec registers the sdk message type.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterInterface((*Msg)(nil), nil)
	cdc.RegisterInterface((*Tx)(nil), nil)
}

// RegisterInterfaces registers the sdk message type.
func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterInterface("lfb.base.v1beta1.Msg", (*Msg)(nil))
	// the interface name for MsgRequest is ServiceMsg because this is most useful for clients
	// to understand - it will be the way for clients to introspect on available Msg service methods
	registry.RegisterInterface("lfb.base.v1beta1.ServiceMsg", (*MsgRequest)(nil))
}
