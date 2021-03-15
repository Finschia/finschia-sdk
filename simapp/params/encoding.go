package params

import (
	"github.com/line/lbm-sdk/v2/client"
	"github.com/line/lbm-sdk/v2/codec"
	"github.com/line/lbm-sdk/v2/codec/types"
)

// EncodingConfig specifies the concrete encoding types to use for a given app.
// This is provided for compatibility between protobuf and amino implementations.
type EncodingConfig struct {
	InterfaceRegistry types.InterfaceRegistry
	Marshaler         codec.Marshaler
	TxConfig          client.TxConfig
	Amino             *codec.LegacyAmino
}
