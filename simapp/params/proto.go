// +build !test_amino

package params

import (
	"github.com/line/lbm-sdk/v2/codec"
	"github.com/line/lbm-sdk/v2/codec/types"
	"github.com/line/lbm-sdk/v2/x/auth/tx"
)

// MakeTestEncodingConfig creates an EncodingConfig for a non-amino based test configuration.
// This function should be used only internally (in the SDK).
// App user should'nt create new codecs - use the app.AppCodec instead.
// [DEPRECATED]
func MakeTestEncodingConfig() EncodingConfig {
	cdc := codec.NewLegacyAmino()
	interfaceRegistry := types.NewInterfaceRegistry()
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	return EncodingConfig{
		InterfaceRegistry: interfaceRegistry,
		Marshaler:         marshaler,
		TxConfig:          tx.NewTxConfig(marshaler, tx.DefaultSignModes),
		Amino:             cdc,
	}
}
