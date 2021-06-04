package simulation

import (
	"bytes"
	"fmt"

	"github.com/line/lbm-sdk/v2/codec"
	types3 "github.com/line/lbm-sdk/v2/store/types"
	types2 "github.com/line/lbm-sdk/v2/types"
	"github.com/line/lbm-sdk/v2/x/evidence/exported"
	keeper2 "github.com/line/lbm-sdk/v2/x/evidence/keeper"
	"github.com/line/lbm-sdk/v2/x/evidence/types"
)

// NewDecodeStore returns a decoder function closure that unmarshals the KVPair's
// Value to the corresponding evidence type.
func NewDecodeStore(cdc codec.Marshaler) types2.StoreDecoder {
	return types2.StoreDecoder{
		Marshal: func(key []byte) func(obj interface{}) []byte {
			switch {
			case bytes.Equal(key[:1], types.KeyPrefixEvidence):
				return keeper2.GetEvidenceMarshalFunc(cdc)
			default:
				panic(fmt.Sprintf("unexpected %s key %X (%s)", types.ModuleName, key, key))
			}
		},
		Unmarshal: func(key []byte) func(value []byte) interface{} {
			switch {
			case bytes.Equal(key[:1], types.KeyPrefixEvidence):
				return keeper2.GetEvidenceUnmarshalFunc(cdc)
			default:
				panic(fmt.Sprintf("unexpected %s key %X (%s)", types.ModuleName, key, key))
			}
		},
		LogPair: func(kvA, kvB types3.KOPair) string {
			switch {
			case bytes.Equal(kvA.Key[:1], types.KeyPrefixEvidence):
				evidenceA := kvA.Value.(exported.Evidence)
				evidenceB := kvB.Value.(exported.Evidence)

				return fmt.Sprintf("%v\n%v", evidenceA, evidenceB)
			default:
				panic(fmt.Sprintf("invalid %s key prefix %X", types.ModuleName, kvA.Key[:1]))
			}
		},
	}
}
