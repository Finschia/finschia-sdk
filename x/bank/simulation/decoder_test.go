package simulation_test

import (
	"fmt"
	"testing"

	types2 "github.com/line/lbm-sdk/v2/store/types"
	"github.com/stretchr/testify/require"

	"github.com/line/lbm-sdk/v2/simapp"
	sdk "github.com/line/lbm-sdk/v2/types"
	"github.com/line/lbm-sdk/v2/x/bank/simulation"
	"github.com/line/lbm-sdk/v2/x/bank/types"
)

func TestDecodeStore(t *testing.T) {
	cdc, _ := simapp.MakeCodecs()
	dec := simulation.NewDecodeStore(cdc)
	totalSupply := types.NewSupply(sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 1000)))

	kvPairs := []types2.KOPair{
		{
			Key: types.SupplyKey, Value: totalSupply,
		},
		{
			Key: []byte{0x99}, Value: []byte{0x99},
		},
	}

	tests := []struct {
		name        string
		expectedLog string
	}{
		{"Supply", fmt.Sprintf("%v\n%v", totalSupply, totalSupply)},
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
