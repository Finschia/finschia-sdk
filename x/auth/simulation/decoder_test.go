package simulation_test

import (
	"fmt"
	"testing"

	gogotypes "github.com/gogo/protobuf/types"
	types2 "github.com/line/lbm-sdk/v2/store/types"
	"github.com/line/lbm-sdk/v2/types/kv"
	"github.com/stretchr/testify/require"

	"github.com/line/lbm-sdk/v2/crypto/keys/ed25519"
	"github.com/line/lbm-sdk/v2/simapp"
	sdk "github.com/line/lbm-sdk/v2/types"
	"github.com/line/lbm-sdk/v2/x/auth/simulation"
	"github.com/line/lbm-sdk/v2/x/auth/types"
)

var (
	delPk1   = ed25519.GenPrivKey().PubKey()
	delAddr1 = sdk.AccAddress(delPk1.Address())
)

func TestDecodeStore(t *testing.T) {
	app := simapp.Setup(false)
	cdc, _ := simapp.MakeCodecs()
	acc := types.NewBaseAccountWithAddress(delAddr1)
	dec := simulation.NewDecodeStore(app.AccountKeeper)

	accBz, err := app.AccountKeeper.MarshalAccount(acc)
	require.NoError(t, err)

	globalAccNumber := gogotypes.UInt64Value{Value: 10}

	kvPairs := kv.Pairs{
		Pairs: []kv.Pair{
			{
				Key:   types.AddressStoreKey(delAddr1),
				Value: accBz,
			},
			{
				Key:   types.GlobalAccountNumberKey,
				Value: cdc.MustMarshalBinaryBare(&globalAccNumber),
			},
			{
				Key:   []byte{0x99},
				Value: []byte{0x99},
			},
		},
	}
	tests := []struct {
		name        string
		expectedLog string
	}{
		{"Account", fmt.Sprintf("%v\n%v", acc, acc)},
		{"GlobalAccNumber", fmt.Sprintf("GlobalAccNumberA: %d\nGlobalAccNumberB: %d", globalAccNumber, globalAccNumber)},
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
