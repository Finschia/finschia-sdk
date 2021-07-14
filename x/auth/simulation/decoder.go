package simulation

import (
	"bytes"
	"fmt"

	"github.com/line/lfb-sdk/codec"
	"github.com/line/lfb-sdk/types/kv"
	"github.com/line/lfb-sdk/x/auth/types"
)

type AuthUnmarshaler interface {
	UnmarshalAccount([]byte) (types.AccountI, error)
	GetCodec() codec.BinaryMarshaler
}

// NewDecodeStore returns a decoder function closure that unmarshals the KVPair's
// Value to the corresponding auth type.
func NewDecodeStore(ak AuthUnmarshaler) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key[:1], types.AddressStoreKeyPrefix):
			accA, err := ak.UnmarshalAccount(kvA.Value)
			if err != nil {
				panic(err)
			}

			accB, err := ak.UnmarshalAccount(kvB.Value)
			if err != nil {
				panic(err)
			}

			return fmt.Sprintf("%v\n%v", accA, accB)

		default:
			panic(fmt.Sprintf("unexpected %s key %X (%s)", types.ModuleName, kvA.Key, kvA.Key))
		}
	}
}
