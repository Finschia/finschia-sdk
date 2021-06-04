package simulation_test

import (
	"fmt"
	"testing"

	types2 "github.com/line/lbm-sdk/v2/store/types"
	"github.com/stretchr/testify/require"

	"github.com/line/lbm-sdk/v2/simapp"
	"github.com/line/lbm-sdk/v2/x/ibc/applications/transfer/simulation"
	"github.com/line/lbm-sdk/v2/x/ibc/applications/transfer/types"
)

func TestDecodeStore(t *testing.T) {
	cdc, _ := simapp.MakeCodecs()
	dec := simulation.NewDecodeStore(cdc)

	trace := types.DenomTrace{
		BaseDenom: "uatom",
		Path:      "transfer/channelToA",
	}

	kvPairs := []types2.KOPair{
			{
				Key:   types.PortKey,
				Value: []byte(types.PortID),
			},
			{
				Key:   types.DenomTraceKey,
				Value: &trace,
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
		{"PortID", fmt.Sprintf("Port A: %s\nPort B: %s", types.PortID, types.PortID)},
		{"DenomTrace", fmt.Sprintf("DenomTrace A: %s\nDenomTrace B: %s", trace.IBCDenom(), trace.IBCDenom())},
		{"other", ""},
	}

	for i, tt := range tests {
		i, tt := i, tt
		t.Run(tt.name, func(t *testing.T) {
			if i == len(tests)-1 {
				require.Panics(t, func() { dec.LogPair(kvPairs[i], kvPairs[i]) }, tt.name)
			} else {
				require.Equal(t, tt.expectedLog, dec.LogPair(kvPairs[i], kvPairs[i]), tt.name)
			}
		})
	}
}
