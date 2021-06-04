package simulation

import (
	"bytes"
	"fmt"

	gogotypes "github.com/gogo/protobuf/types"
	types4 "github.com/line/lbm-sdk/v2/store/types"
	types2 "github.com/line/lbm-sdk/v2/types"
	keeper2 "github.com/line/lbm-sdk/v2/x/slashing/keeper"

	"github.com/line/lbm-sdk/v2/codec"
	"github.com/line/lbm-sdk/v2/x/slashing/types"
)

func GetStringValueValueUnmarshalFunc(cdc codec.BinaryMarshaler) func (value []byte) interface{} {
	return func (value []byte) interface{} {
		var val gogotypes.StringValue
		cdc.MustUnmarshalBinaryBare(value, &val)
		return &val
	}
}

func GetStringValueValueMarshalFunc(cdc codec.BinaryMarshaler) func (obj interface{}) []byte {
	return func (obj interface{}) []byte {
		return cdc.MustMarshalBinaryBare(obj.(*gogotypes.StringValue))
	}
}

// NewDecodeStore returns a decoder function closure that unmarshals the KVPair's
// Value to the corresponding slashing type.
func NewDecodeStore(cdc codec.Marshaler) types2.StoreDecoder {
	return types2.StoreDecoder{
		Marshal: func(key []byte) func(obj interface{}) []byte {
			switch {
			case bytes.Equal(key[:1], types.ValidatorSigningInfoKeyPrefix):
				return keeper2.GetValidatorSigningInfoMarshalFunc(cdc)
			case bytes.Equal(key[:1], types.ValidatorMissedBlockBitArrayKeyPrefix):
				return keeper2.GetBoolValueMarshalFunc(cdc)
			case bytes.Equal(key[:1], types.AddrPubkeyRelationKeyPrefix):
				return GetStringValueValueMarshalFunc(cdc)
			default:
				panic(fmt.Sprintf("unexpected %s key %X (%s)", types.ModuleName, key, key))
			}
		},
		Unmarshal: func(key []byte) func(value []byte) interface{} {
			switch {
			case bytes.Equal(key[:1], types.ValidatorSigningInfoKeyPrefix):
				return keeper2.GetValidatorSigningInfoUnmarshalFunc(cdc)
			case bytes.Equal(key[:1], types.ValidatorMissedBlockBitArrayKeyPrefix):
				return keeper2.GetBoolValueUnmarshalFunc(cdc)
			case bytes.Equal(key[:1], types.AddrPubkeyRelationKeyPrefix):
				return GetStringValueValueUnmarshalFunc(cdc)

			default:
				panic(fmt.Sprintf("unexpected %s key %X (%s)", types.ModuleName, key, key))
			}
		},
		LogPair: func(kvA, kvB types4.KOPair) string {
			switch {
			case bytes.Equal(kvA.Key[:1], types.ValidatorSigningInfoKeyPrefix):
				infoA := *kvA.Value.(*types.ValidatorSigningInfo)
				infoB := *kvB.Value.(*types.ValidatorSigningInfo)
				return fmt.Sprintf("%v\n%v", infoA, infoB)

			case bytes.Equal(kvA.Key[:1], types.ValidatorMissedBlockBitArrayKeyPrefix):
				missedA := *kvA.Value.(*gogotypes.BoolValue)
				missedB := *kvB.Value.(*gogotypes.BoolValue)
				return fmt.Sprintf("missedA: %v\nmissedB: %v", missedA.Value, missedB.Value)

			case bytes.Equal(kvA.Key[:1], types.AddrPubkeyRelationKeyPrefix):
				pubKeyA := *kvA.Value.(*gogotypes.StringValue)
				pubKeyB := *kvB.Value.(*gogotypes.StringValue)
				return fmt.Sprintf("PubKeyA: %s\nPubKeyB: %s", pubKeyA.Value, pubKeyB.Value)

			default:
				panic(fmt.Sprintf("invalid slashing key prefix %X", kvA.Key[:1]))
			}
		},
	}
}
