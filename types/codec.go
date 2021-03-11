package types

import (
	"github.com/line/lbm-sdk/codec"
	"github.com/line/lbm-sdk/codec/types"
)

const (
	// MsgInterfaceProtoName defines the protobuf name of the lbm Msg interface
	MsgInterfaceProtoName = "lbm.base.v1.Msg"
	// ServiceMsgInterfaceProtoName defines the protobuf name of the lbm MsgRequest interface
	ServiceMsgInterfaceProtoName = "lbm.base.v1.ServiceMsg"
)

// RegisterLegacyAminoCodec registers the sdk message type.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterInterface((*Msg)(nil), nil)
	cdc.RegisterInterface((*Tx)(nil), nil)
}

// RegisterInterfaces registers the sdk message type.
func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterInterface(MsgInterfaceProtoName, (*Msg)(nil))
	// the interface name for MsgRequest is ServiceMsg because this is most useful for clients
	// to understand - it will be the way for clients to introspect on available Msg service methods
	registry.RegisterInterface(ServiceMsgInterfaceProtoName, (*MsgRequest)(nil))
}
