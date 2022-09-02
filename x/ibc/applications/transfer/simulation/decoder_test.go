package simulation_test

import (
	"fmt"
	"testing"

	"github.com/line/lbm-sdk/types/kv"
	"github.com/stretchr/testify/require"

	"github.com/line/lbm-sdk/simapp"
	"github.com/line/lbm-sdk/x/ibc/applications/transfer/simulation"
	"github.com/line/lbm-sdk/x/ibc/applications/transfer/types"
)

func TestDecodeStore(t *testing.T) {
	app := simapp.Setup(false)
	dec := simulation.NewDecodeStore(app.TransferKeeper)

	trace := types.DenomTrace{
		BaseDenom: "uatom",
		Path:      "transfer/channelToA",
	}

	kvPairs := kv.Pairs{
		Pairs: []kv.Pair{
			{
				Key:   types.PortKey,
				Value: []byte(types.PortID),
			},
			{
				Key:   types.DenomTraceKey,
				Value: app.TransferKeeper.MustMarshalDenomTrace(trace),
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
		{"PortID", fmt.Sprintf("Port A: %s\nPort B: %s", types.PortID, types.PortID)},
		{"DenomTrace", fmt.Sprintf("DenomTrace A: %s\nDenomTrace B: %s", trace.IBCDenom(), trace.IBCDenom())},
		{"other", ""},
	}

	for i, tt := range tests {
		i, tt := i, tt
		t.Run(tt.name, func(t *testing.T) {
			if i == len(tests)-1 {
				require.Panics(t, func() { dec(kvPairs.Pairs[i], kvPairs.Pairs[i]) }, tt.name)
			} else {
				require.Equal(t, tt.expectedLog, dec(kvPairs.Pairs[i], kvPairs.Pairs[i]), tt.name)
			}
		})
	}
}
