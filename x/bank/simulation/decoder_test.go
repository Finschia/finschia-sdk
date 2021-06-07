package simulation_test

import (
	"fmt"
	"testing"

	types2 "github.com/line/lbm-sdk/v2/store/types"
	"github.com/line/lbm-sdk/v2/types/kv"
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

	supplyBz, err := cdc.MarshalInterface(totalSupply)
	require.NoError(t, err)

	kvPairs := kv.Pairs{
		Pairs: []kv.Pair{
			{Key: types.SupplyKey, Value: supplyBz},
			{Key: []byte{0x99}, Value: []byte{0x99}},
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
			var value interface{}
			if i == len(tests) - 1 {
				require.Panics(t, func () {dec.Unmarshal(kvPairs.Pairs[i].Key)(kvPairs.Pairs[i].Value)}, tt.name)
				value = nil
			} else {
				value = dec.Unmarshal(kvPairs.Pairs[i].Key)(kvPairs.Pairs[i].Value)
			}
			pair := types2.KOPair{
				Key:   kvPairs.Pairs[i].Key,
				Value: value,
			}
			switch i {
			case len(tests) - 1:
				require.Panics(t, func() { dec.LogPair(pair, pair) }, tt.name)
			default:
				require.Equal(t, tt.expectedLog, dec.LogPair(pair, pair), tt.name)
			}
		})
	}
}
