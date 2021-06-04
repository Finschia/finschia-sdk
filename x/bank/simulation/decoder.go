package simulation

import (
	"bytes"
	"fmt"

	"github.com/line/lbm-sdk/v2/codec"
	types3 "github.com/line/lbm-sdk/v2/store/types"
	types2 "github.com/line/lbm-sdk/v2/types"
	"github.com/line/lbm-sdk/v2/x/bank/exported"
	keeper2 "github.com/line/lbm-sdk/v2/x/bank/keeper"
	"github.com/line/lbm-sdk/v2/x/bank/types"
)

// SupplyUnmarshaler defines the expected encoding store functions.
type SupplyUnmarshaler interface {
	UnmarshalSupply([]byte) (exported.SupplyI, error)
}

// NewDecodeStore returns a function closure that unmarshals the KVPair's values
// to the corresponding types.
func NewDecodeStore(cdc codec.Marshaler) types2.StoreDecoder {
	return types2.StoreDecoder{
		Marshal: func(key []byte) func(obj interface{}) []byte {
			switch {
			case bytes.Equal(key[:1], types.SupplyKey):
				return keeper2.GetSupplyMarshalFunc(cdc)
			default:
				panic(fmt.Sprintf("unexpected %s key %X (%s)", types.ModuleName, key, key))
			}
		},
		Unmarshal: func(key []byte) func(value []byte) interface{} {
			switch {
			case bytes.Equal(key[:1], types.SupplyKey):
				return keeper2.GetSupplyUnmarshalFunc(cdc)
			default:
				panic(fmt.Sprintf("unexpected %s key %X (%s)", types.ModuleName, key, key))
			}
		},
		LogPair: func(kvA, kvB types3.KOPair) string {
			switch {
			case bytes.Equal(kvA.Key[:1], types.SupplyKey):
				supplyA := kvA.Value.(exported.SupplyI)
				supplyB := kvB.Value.(exported.SupplyI)

				return fmt.Sprintf("%v\n%v", supplyA, supplyB)

			default:
				panic(fmt.Sprintf("unexpected %s key %X (%s)", types.ModuleName, kvA.Key, kvA.Key))
			}
		},
	}
}
