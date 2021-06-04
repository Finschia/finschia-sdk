package simulation

import (
	"bytes"
	"fmt"

	"github.com/line/lbm-sdk/v2/codec"
	types3 "github.com/line/lbm-sdk/v2/store/types"
	types2 "github.com/line/lbm-sdk/v2/types"
	keeper2 "github.com/line/lbm-sdk/v2/x/ibc/applications/transfer/keeper"
	"github.com/line/lbm-sdk/v2/x/ibc/applications/transfer/types"
)

// TransferUnmarshaler defines the expected encoding store functions.
type TransferUnmarshaler interface {
	MustUnmarshalDenomTrace([]byte) types.DenomTrace
}

// NewDecodeStore returns a decoder function closure that unmarshals the KVPair's
// Value to the corresponding DenomTrace type.
func NewDecodeStore(cdc codec.Marshaler) types2.StoreDecoder {
	return types2.StoreDecoder{
		Marshal: func(key []byte) func(obj interface{}) []byte {
			switch {
			case bytes.Equal(key[:1], types.PortKey):
				return types3.GetBytesMarshalFunc()
			case bytes.Equal(key[:1], types.DenomTraceKey):
				return keeper2.GetDenomTraceMarshalFunc(cdc)
			default:
				panic(fmt.Sprintf("unexpected %s key %X (%s)", types.ModuleName, key, key))
			}
		},
		Unmarshal: func(key []byte) func(value []byte) interface{} {
			switch {
			case bytes.Equal(key[:1], types.PortKey):
				return types3.GetBytesUnmarshalFunc()
			case bytes.Equal(key[:1], types.DenomTraceKey):
				return keeper2.GetDenomTraceUnmarshalFunc(cdc)
			default:
				panic(fmt.Sprintf("unexpected %s key %X (%s)", types.ModuleName, key, key))
			}
		},
		LogPair: func(kvA, kvB types3.KOPair) string {
			switch {
			case bytes.Equal(kvA.Key[:1], types.PortKey):
				return fmt.Sprintf("Port A: %s\nPort B: %s", string(kvA.Value.([]byte)),
					string(kvB.Value.([]byte)))

			case bytes.Equal(kvA.Key[:1], types.DenomTraceKey):
				denomTraceA := *kvA.Value.(*types.DenomTrace)
				denomTraceB := *kvB.Value.(*types.DenomTrace)
				return fmt.Sprintf("DenomTrace A: %s\nDenomTrace B: %s",
					denomTraceA.IBCDenom(), denomTraceB.IBCDenom())

			default:
				panic(fmt.Sprintf("invalid %s key prefix %X", types.ModuleName, kvA.Key[:1]))
			}
		},
	}
}
