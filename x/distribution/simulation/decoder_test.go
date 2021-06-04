package simulation_test

import (
	"fmt"
	"testing"

	types2 "github.com/line/lbm-sdk/v2/store/types"
	"github.com/stretchr/testify/require"

	"github.com/line/lbm-sdk/v2/crypto/keys/ed25519"
	"github.com/line/lbm-sdk/v2/simapp"
	sdk "github.com/line/lbm-sdk/v2/types"
	"github.com/line/lbm-sdk/v2/x/distribution/simulation"
	"github.com/line/lbm-sdk/v2/x/distribution/types"
)

var (
	delPk1    = ed25519.GenPrivKey().PubKey()
	delAddr1  = sdk.AccAddress(delPk1.Address())
	valAddr1  = sdk.ValAddress(delPk1.Address())
	consAddr1 = sdk.ConsAddress(delPk1.Address().Bytes())
)

func TestDecodeDistributionStore(t *testing.T) {
	cdc, _ := simapp.MakeCodecs()
	dec := simulation.NewDecodeStore(cdc)

	decCoins := sdk.DecCoins{sdk.NewDecCoinFromDec(sdk.DefaultBondDenom, sdk.OneDec())}
	feePool := types.InitialFeePool()
	feePool.CommunityPool = decCoins
	info := types.NewDelegatorStartingInfo(2, sdk.OneDec(), 200)
	outstanding := types.ValidatorOutstandingRewards{Rewards: decCoins}
	commission := types.ValidatorAccumulatedCommission{Commission: decCoins}
	historicalRewards := types.NewValidatorHistoricalRewards(decCoins, 100)
	currentRewards := types.NewValidatorCurrentRewards(decCoins, 5)
	slashEvent := types.NewValidatorSlashEvent(10, sdk.OneDec())

	kvPairs := []types2.KOPair{
			{Key: types.FeePoolKey, Value: &feePool},
			{Key: types.ProposerKey, Value: consAddr1.Bytes()},
			{Key: types.GetValidatorOutstandingRewardsKey(valAddr1), Value: &outstanding},
			{Key: types.GetDelegatorWithdrawAddrKey(delAddr1), Value: delAddr1.Bytes()},
			{Key: types.GetDelegatorStartingInfoKey(valAddr1, delAddr1), Value: &info},
			{Key: types.GetValidatorHistoricalRewardsKey(valAddr1, 100), Value: &historicalRewards},
			{Key: types.GetValidatorCurrentRewardsKey(valAddr1), Value: &currentRewards},
			{Key: types.GetValidatorAccumulatedCommissionKey(valAddr1), Value: &commission},
			{Key: types.GetValidatorSlashEventKeyPrefix(valAddr1, 13), Value: &slashEvent},
			{Key: []byte{0x99}, Value: []byte{0x99}},
	}

	tests := []struct {
		name        string
		expectedLog string
	}{
		{"FeePool", fmt.Sprintf("%v\n%v", feePool, feePool)},
		{"Proposer", fmt.Sprintf("%v\n%v", consAddr1, consAddr1)},
		{"ValidatorOutstandingRewards", fmt.Sprintf("%v\n%v", outstanding, outstanding)},
		{"DelegatorWithdrawAddr", fmt.Sprintf("%v\n%v", delAddr1, delAddr1)},
		{"DelegatorStartingInfo", fmt.Sprintf("%v\n%v", info, info)},
		{"ValidatorHistoricalRewards", fmt.Sprintf("%v\n%v", historicalRewards, historicalRewards)},
		{"ValidatorCurrentRewards", fmt.Sprintf("%v\n%v", currentRewards, currentRewards)},
		{"ValidatorAccumulatedCommission", fmt.Sprintf("%v\n%v", commission, commission)},
		{"ValidatorSlashEvent", fmt.Sprintf("%v\n%v", slashEvent, slashEvent)},
		{"other", ""},
	}
	for i, tt := range tests {
		i, tt := i, tt
		t.Run(tt.name, func(t *testing.T) {
			switch i {
			case len(tests) - 1:
				require.Panics(t, func() { dec.LogPair(kvPairs[i], kvPairs[i]) }, tt.name)
			default:
				require.Equal(t, tt.expectedLog, dec.LogPair(kvPairs[i], kvPairs[i]), tt.name)
			}
		})
	}
}
