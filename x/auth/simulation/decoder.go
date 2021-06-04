package simulation

import (
	"bytes"
	"fmt"

	gogotypes "github.com/gogo/protobuf/types"
	types3 "github.com/line/lbm-sdk/v2/store/types"
	types2 "github.com/line/lbm-sdk/v2/types"
	"github.com/line/lbm-sdk/v2/x/auth/keeper"

	"github.com/line/lbm-sdk/v2/codec"
	"github.com/line/lbm-sdk/v2/x/auth/types"
)

type AuthUnmarshaler interface {
	UnmarshalAccount([]byte) (types.AccountI, error)
	GetCodec() codec.BinaryMarshaler
}

// NewDecodeStore returns a decoder function closure that unmarshals the KVPair's
// Value to the corresponding auth type.
func NewDecodeStore(ak AuthUnmarshaler) types2.StoreDecoder {
	return types2.StoreDecoder{
		Marshal: func(key []byte) func(obj interface{}) []byte {
			switch {
			case bytes.Equal(key[:1], types.AddressStoreKeyPrefix):
				return types.GetAccountMarshalFunc(ak.GetCodec())
			case bytes.Equal(key, types.GlobalAccountNumberKey):
				return keeper.GetUInt64ValueMarshalFunc(ak.GetCodec())
			default:
				panic(fmt.Sprintf("unexpected %s key %X (%s)", types.ModuleName, key, key))
			}
		},
		Unmarshal: func(key []byte) func(value []byte) interface{} {
			switch {
			case bytes.Equal(key[:1], types.AddressStoreKeyPrefix):
				return types.GetAccountUnmarshalFunc(ak.GetCodec())
			case bytes.Equal(key, types.GlobalAccountNumberKey):
				return keeper.GetUInt64ValueUnmarshalFunc(ak.GetCodec())
			default:
				panic(fmt.Sprintf("unexpected %s key %X (%s)", types.ModuleName, key, key))
			}
		},
		LogPair: func(kvA, kvB types3.KOPair) string {
			switch {
			case bytes.Equal(kvA.Key[:1], types.AddressStoreKeyPrefix):
				accA := kvA.Value.(types.AccountI)
				accB := kvB.Value.(types.AccountI)

				return fmt.Sprintf("%v\n%v", accA, accB)

			case bytes.Equal(kvA.Key, types.GlobalAccountNumberKey):
				globalAccNumberA := *kvA.Value.(*gogotypes.UInt64Value)
				globalAccNumberB := *kvB.Value.(*gogotypes.UInt64Value)

				return fmt.Sprintf("GlobalAccNumberA: %d\nGlobalAccNumberB: %d", globalAccNumberA, globalAccNumberB)

			default:
				panic(fmt.Sprintf("unexpected %s key %X (%s)", types.ModuleName, kvA.Key, kvA.Key))
			}
		},
	}
}
