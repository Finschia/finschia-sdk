package simulation_test

import (
	"fmt"
	"testing"

	gogotypes "github.com/gogo/protobuf/types"
	types2 "github.com/line/lbm-sdk/v2/store/types"
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
	acc := types.NewBaseAccountWithAddress(delAddr1)
	dec := simulation.NewDecodeStore(app.AccountKeeper)

	globalAccNumber := gogotypes.UInt64Value{Value: 10}

	kvPairs := []types2.KOPair{
			{
				Key:   types.AddressStoreKey(delAddr1),
				Value: acc,
			},
			{
				Key:   types.GlobalAccountNumberKey,
				Value: &globalAccNumber,
			},
			{
				Key:   []byte{0x99},
				Value: []byte{0x99},
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
			switch i {
			case len(tests) - 1:
				require.Panics(t, func() { dec.LogPair(kvPairs[i], kvPairs[i]) }, tt.name)
			default:
				require.Equal(t, tt.expectedLog, dec.LogPair(kvPairs[i], kvPairs[i]), tt.name)
			}
		})
	}
}
