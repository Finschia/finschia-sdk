package simulation

import (
	"bytes"
	"fmt"

	"github.com/line/lbm-sdk/v2/codec"
	types2 "github.com/line/lbm-sdk/v2/store/types"
	sdk "github.com/line/lbm-sdk/v2/types"
	keeper2 "github.com/line/lbm-sdk/v2/x/distribution/keeper"
	"github.com/line/lbm-sdk/v2/x/distribution/types"
)

// NewDecodeStore returns a decoder function closure that unmarshals the KVPair's
// Value to the corresponding distribution type.
func NewDecodeStore(cdc codec.Marshaler) sdk.StoreDecoder {
	return sdk.StoreDecoder{
		Marshal: func(key []byte) func(obj interface{}) []byte {
			switch {
			case bytes.Equal(key[:1], types.FeePoolKey):
				return keeper2.GetFeePoolMarshalFunc(cdc)
			case bytes.Equal(key[:1], types.ProposerKey):
				return keeper2.GetByteValueMarshalFunc(cdc)
			case bytes.Equal(key[:1], types.ValidatorOutstandingRewardsPrefix):
				return keeper2.GetValidatorOutstandingRewardsMarshalFunc(cdc)
			case bytes.Equal(key[:1], types.DelegatorWithdrawAddrPrefix):
				return types2.GetBytesMarshalFunc()
			case bytes.Equal(key[:1], types.DelegatorStartingInfoPrefix):
				return keeper2.GetDelegatorStartingInfoMarshalFunc(cdc)
			case bytes.Equal(key[:1], types.ValidatorHistoricalRewardsPrefix):
				return keeper2.GetValidatorHistoricalRewardsMarshalFunc(cdc)
			case bytes.Equal(key[:1], types.ValidatorCurrentRewardsPrefix):
				return keeper2.GetValidatorCurrentRewardsMarshalFunc(cdc)
			case bytes.Equal(key[:1], types.ValidatorAccumulatedCommissionPrefix):
				return keeper2.GetValidatorAccumulatedCommissionMarshalFunc(cdc)
			case bytes.Equal(key[:1], types.ValidatorSlashEventPrefix):
				return keeper2.GetValidatorSlashEventMarshalFunc(cdc)
			default:
				panic(fmt.Sprintf("unexpected %s key %X (%s)", types.ModuleName, key, key))
			}
		},
		Unmarshal: func(key []byte) func(value []byte) interface{} {
			switch {
			case bytes.Equal(key[:1], types.FeePoolKey):
				return keeper2.GetFeePoolUnmarshalFunc(cdc)
			case bytes.Equal(key[:1], types.ProposerKey):
				return keeper2.GetByteValueUnmarshalFunc(cdc)
			case bytes.Equal(key[:1], types.ValidatorOutstandingRewardsPrefix):
				return keeper2.GetValidatorOutstandingRewardsUnmarshalFunc(cdc)
			case bytes.Equal(key[:1], types.DelegatorWithdrawAddrPrefix):
				return types2.GetBytesUnmarshalFunc()
			case bytes.Equal(key[:1], types.DelegatorStartingInfoPrefix):
				return keeper2.GetDelegatorStartingInfoUnmarshalFunc(cdc)
			case bytes.Equal(key[:1], types.ValidatorHistoricalRewardsPrefix):
				return keeper2.GetValidatorHistoricalRewardsUnmarshalFunc(cdc)
			case bytes.Equal(key[:1], types.ValidatorCurrentRewardsPrefix):
				return keeper2.GetValidatorCurrentRewardsUnmarshalFunc(cdc)
			case bytes.Equal(key[:1], types.ValidatorAccumulatedCommissionPrefix):
				return keeper2.GetValidatorAccumulatedCommissionUnmarshalFunc(cdc)
			case bytes.Equal(key[:1], types.ValidatorSlashEventPrefix):
				return keeper2.GetValidatorSlashEventUnmarshalFunc(cdc)
			default:
				panic(fmt.Sprintf("unexpected %s key %X (%s)", types.ModuleName, key, key))
			}
		},
		LogPair: func(kvA, kvB types2.KOPair) string {
			switch {
			case bytes.Equal(kvA.Key[:1], types.FeePoolKey):
				feePoolA := *kvA.Value.(*types.FeePool)
				feePoolB := *kvB.Value.(*types.FeePool)
				return fmt.Sprintf("%v\n%v", feePoolA, feePoolB)

			case bytes.Equal(kvA.Key[:1], types.ProposerKey):
				return fmt.Sprintf("%v\n%v", sdk.ConsAddress(kvA.Value.([]byte)),
					sdk.ConsAddress(kvB.Value.([]byte)))

			case bytes.Equal(kvA.Key[:1], types.ValidatorOutstandingRewardsPrefix):
				rewardsA := *kvA.Value.(*types.ValidatorOutstandingRewards)
				rewardsB := *kvB.Value.(*types.ValidatorOutstandingRewards)
				return fmt.Sprintf("%v\n%v", rewardsA, rewardsB)

			case bytes.Equal(kvA.Key[:1], types.DelegatorWithdrawAddrPrefix):
				return fmt.Sprintf("%v\n%v", sdk.AccAddress(kvA.Value.([]byte)),
					sdk.AccAddress(kvB.Value.([]byte)))

			case bytes.Equal(kvA.Key[:1], types.DelegatorStartingInfoPrefix):
				infoA := *kvA.Value.(*types.DelegatorStartingInfo)
				infoB := *kvB.Value.(*types.DelegatorStartingInfo)
				return fmt.Sprintf("%v\n%v", infoA, infoB)

			case bytes.Equal(kvA.Key[:1], types.ValidatorHistoricalRewardsPrefix):
				rewardsA := *kvA.Value.(*types.ValidatorHistoricalRewards)
				rewardsB := *kvB.Value.(*types.ValidatorHistoricalRewards)
				return fmt.Sprintf("%v\n%v", rewardsA, rewardsB)

			case bytes.Equal(kvA.Key[:1], types.ValidatorCurrentRewardsPrefix):
				rewardsA := *kvA.Value.(*types.ValidatorCurrentRewards)
				rewardsB := *kvB.Value.(*types.ValidatorCurrentRewards)
				return fmt.Sprintf("%v\n%v", rewardsA, rewardsB)

			case bytes.Equal(kvA.Key[:1], types.ValidatorAccumulatedCommissionPrefix):
				commissionA := *kvA.Value.(*types.ValidatorAccumulatedCommission)
				commissionB := *kvA.Value.(*types.ValidatorAccumulatedCommission)
				return fmt.Sprintf("%v\n%v", commissionA, commissionB)

			case bytes.Equal(kvA.Key[:1], types.ValidatorSlashEventPrefix):
				eventA := *kvA.Value.(*types.ValidatorSlashEvent)
				eventB := *kvB.Value.(*types.ValidatorSlashEvent)
				return fmt.Sprintf("%v\n%v", eventA, eventB)

			default:
				panic(fmt.Sprintf("invalid distribution key prefix %X", kvA.Key[:1]))
			}
		},
	}
}
