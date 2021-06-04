package simulation

import (
	"bytes"
	"fmt"

	"github.com/line/lbm-sdk/v2/codec"
	types3 "github.com/line/lbm-sdk/v2/store/types"
	types2 "github.com/line/lbm-sdk/v2/types"
	"github.com/line/lbm-sdk/v2/x/mint/keeper"
	"github.com/line/lbm-sdk/v2/x/mint/types"
)

// NewDecodeStore returns a decoder function closure that umarshals the KVPair's
// Value to the corresponding mint type.
func NewDecodeStore(cdc codec.Marshaler) types2.StoreDecoder {
	return types2.StoreDecoder{
		Marshal: func(key []byte) func(obj interface{}) []byte {
			switch {
			case bytes.Equal(key, types.MinterKey):
				return keeper.GetMinterMarshalFunc(cdc)

			default:
				panic(fmt.Sprintf("unexpected %s key %X (%s)", types.ModuleName, key, key))
			}
		},
		Unmarshal: func(key []byte) func(value []byte) interface{} {
			switch {
			case bytes.Equal(key, types.MinterKey):
				return keeper.GetMinterUnmarshalFunc(cdc)

			default:
				panic(fmt.Sprintf("unexpected %s key %X (%s)", types.ModuleName, key, key))
			}
		},
		LogPair: func(kvA, kvB types3.KOPair) string {
			switch {
			case bytes.Equal(kvA.Key, types.MinterKey):
				minterA := *kvA.Value.(*types.Minter)
				minterB := *kvB.Value.(*types.Minter)
				return fmt.Sprintf("%v\n%v", minterA, minterB)
			default:
				panic(fmt.Sprintf("invalid mint key %X", kvA.Key))
			}
		},
	}
}
