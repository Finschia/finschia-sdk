package simulation

import (
	"bytes"
	"fmt"

	"github.com/line/lbm-sdk/v2/codec"
	types2 "github.com/line/lbm-sdk/v2/store/types"
	sdk "github.com/line/lbm-sdk/v2/types"
	"github.com/line/lbm-sdk/v2/x/capability/keeper"
	"github.com/line/lbm-sdk/v2/x/capability/types"
)

// NewDecodeStore returns a decoder function closure that unmarshals the KVPair's
// Value to the corresponding capaility type.
func NewDecodeStore(cdc codec.Marshaler) sdk.StoreDecoder {
	return sdk.StoreDecoder{
		Marshal: func(key []byte) func(obj interface{}) []byte {
			switch {
			case bytes.Equal(key, types.KeyIndex):
				return types2.GetBytesMarshalFunc()
			case bytes.HasPrefix(key, types.KeyPrefixIndexCapability):
				return keeper.GetCapabilityOwnersMarshalFunc(cdc)
			default:
				panic(fmt.Sprintf("unexpected %s key %X (%s)", types.ModuleName, key, key))
			}
		},
		Unmarshal: func(key []byte) func(value []byte) interface{} {
			switch {
			case bytes.Equal(key, types.KeyIndex):
				return types2.GetBytesUnmarshalFunc()
			case bytes.HasPrefix(key, types.KeyPrefixIndexCapability):
				return keeper.GetCapabilityOwnersUnmarshalFunc(cdc)
			default:
				panic(fmt.Sprintf("unexpected %s key %X (%s)", types.ModuleName, key, key))
			}
		},
		LogPair: func(kvA, kvB types2.KOPair) string {
			switch {
			case bytes.Equal(kvA.Key, types.KeyIndex):
				idxA := sdk.BigEndianToUint64(kvA.Value.([]byte))
				idxB := sdk.BigEndianToUint64(kvB.Value.([]byte))
				return fmt.Sprintf("Index A: %d\nIndex B: %d\n", idxA, idxB)

			case bytes.HasPrefix(kvA.Key, types.KeyPrefixIndexCapability):
				capOwnersA := *kvA.Value.(*types.CapabilityOwners)
				capOwnersB := *kvA.Value.(*types.CapabilityOwners)
				return fmt.Sprintf("CapabilityOwners A: %v\nCapabilityOwners B: %v\n", capOwnersA, capOwnersB)

			default:
				panic(fmt.Sprintf("invalid %s key prefix %X (%s)", types.ModuleName, kvA.Key, string(kvA.Key)))
			}
		},
	}
}
