package simulation

import (
	"bytes"
	"fmt"

	"github.com/line/lbm-sdk/v2/codec"
	types2 "github.com/line/lbm-sdk/v2/store/types"
	sdk "github.com/line/lbm-sdk/v2/types"
	"github.com/line/lbm-sdk/v2/x/staking/keeper"
	"github.com/line/lbm-sdk/v2/x/staking/types"
)

// NewDecodeStore returns a decoder function closure that unmarshals the KVPair's
// Value to the corresponding staking type.
func NewDecodeStore(cdc codec.Marshaler) sdk.StoreDecoder {
	return sdk.StoreDecoder{
		Marshal: func(key []byte) func(obj interface{}) []byte {
			switch {
			case bytes.Equal(key[:1], types.LastTotalPowerKey):
				return keeper.GetIntProtoMarshalFunc(cdc)
			case bytes.Equal(key[:1], types.ValidatorsKey):
				return keeper.GetValidatorMarshalFunc(cdc)
			case bytes.Equal(key[:1], types.LastValidatorPowerKey),
				 bytes.Equal(key[:1], types.ValidatorsByConsAddrKey),
				 bytes.Equal(key[:1], types.ValidatorsByPowerIndexKey):
			 	return types2.GetBytesMarshalFunc()
			case bytes.Equal(key[:1], types.DelegationKey):
				return keeper.GetDelegationMarshalFunc(cdc)
			case bytes.Equal(key[:1], types.UnbondingDelegationKey),
				 bytes.Equal(key[:1], types.UnbondingDelegationByValIndexKey):
				return keeper.GetUnbondingDelegationMarshalFunc(cdc)
			case bytes.Equal(key[:1], types.RedelegationKey),
				 bytes.Equal(key[:1], types.RedelegationByValSrcIndexKey):
				return keeper.GetRedelegationMarshalFunc(cdc)
			default:
				panic(fmt.Sprintf("unexpected %s key %X (%s)", types.ModuleName, key, key))
			}
		},
		Unmarshal: func(key []byte) func(value []byte) interface{} {
			switch {
			case bytes.Equal(key[:1], types.LastTotalPowerKey):
				return keeper.GetIntProtoUnmarshalFunc(cdc)
			case bytes.Equal(key[:1], types.ValidatorsKey):
				return keeper.GetValidatorUnmarshalFunc(cdc)
			case bytes.Equal(key[:1], types.LastValidatorPowerKey),
				bytes.Equal(key[:1], types.ValidatorsByConsAddrKey),
				bytes.Equal(key[:1], types.ValidatorsByPowerIndexKey):
				return types2.GetBytesUnmarshalFunc()
			case bytes.Equal(key[:1], types.DelegationKey):
				return keeper.GetDelegationUnmarshalFunc(cdc)
			case bytes.Equal(key[:1], types.UnbondingDelegationKey),
				bytes.Equal(key[:1], types.UnbondingDelegationByValIndexKey):
				return keeper.GetUnbondingDelegationUnmarshalFunc(cdc)
			case bytes.Equal(key[:1], types.RedelegationKey),
				bytes.Equal(key[:1], types.RedelegationByValSrcIndexKey):
				return keeper.GetRedelegationUnmarshalFunc(cdc)

			default:
				panic(fmt.Sprintf("unexpected %s key %X (%s)", types.ModuleName, key, key))
			}
		},
		LogPair: func(kvA, kvB types2.KOPair) string {
			switch {
			case bytes.Equal(kvA.Key[:1], types.LastTotalPowerKey):
				powerA := *kvA.Value.(*sdk.IntProto)
				powerB := *kvB.Value.(*sdk.IntProto)

				return fmt.Sprintf("%v\n%v", powerA, powerB)
			case bytes.Equal(kvA.Key[:1], types.ValidatorsKey):
				validatorA := *kvA.Value.(*types.Validator)
				validatorB := *kvB.Value.(*types.Validator)

				return fmt.Sprintf("%v\n%v", validatorA, validatorB)
			case bytes.Equal(kvA.Key[:1], types.LastValidatorPowerKey),
				bytes.Equal(kvA.Key[:1], types.ValidatorsByConsAddrKey),
				bytes.Equal(kvA.Key[:1], types.ValidatorsByPowerIndexKey):
				return fmt.Sprintf("%v\n%v", sdk.ValAddress(kvA.Value.([]byte)),
					sdk.ValAddress(kvB.Value.([]byte)))

			case bytes.Equal(kvA.Key[:1], types.DelegationKey):
				delegationA := *kvA.Value.(*types.Delegation)
				delegationB := *kvB.Value.(*types.Delegation)

				return fmt.Sprintf("%v\n%v", delegationA, delegationB)
			case bytes.Equal(kvA.Key[:1], types.UnbondingDelegationKey),
				bytes.Equal(kvA.Key[:1], types.UnbondingDelegationByValIndexKey):
				ubdA := *kvA.Value.(*types.UnbondingDelegation)
				ubdB := *kvB.Value.(*types.UnbondingDelegation)

				return fmt.Sprintf("%v\n%v", ubdA, ubdB)
			case bytes.Equal(kvA.Key[:1], types.RedelegationKey),
				bytes.Equal(kvA.Key[:1], types.RedelegationByValSrcIndexKey):
				redA := *kvA.Value.(*types.Redelegation)
				redB := *kvB.Value.(*types.Redelegation)

				return fmt.Sprintf("%v\n%v", redA, redB)
			default:
				panic(fmt.Sprintf("invalid staking key prefix %X", kvA.Key[:1]))
			}
		},
	}
}
