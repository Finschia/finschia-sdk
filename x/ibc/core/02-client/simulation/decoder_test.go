package simulation_test

import (
	"fmt"
	"testing"
	"time"

	types2 "github.com/line/lbm-sdk/v2/store/types"
	"github.com/stretchr/testify/require"

	"github.com/line/lbm-sdk/v2/x/ibc/core/02-client/simulation"
	"github.com/line/lbm-sdk/v2/x/ibc/core/02-client/types"
	host "github.com/line/lbm-sdk/v2/x/ibc/core/24-host"
	ibctmtypes "github.com/line/lbm-sdk/v2/x/ibc/light-clients/07-tendermint/types"
)

func TestDecodeStore(t *testing.T) {
	clientID := "clientidone"

	height := types.NewHeight(0, 10)

	clientState := &ibctmtypes.ClientState{
		FrozenHeight: height,
	}

	consState := &ibctmtypes.ConsensusState{
		Timestamp: time.Now().UTC(),
	}

	kvPairs := []types2.KOPair{
			{
				Key:   host.FullClientStateKey(clientID),
				Value: clientState,
			},
			{
				Key:   host.FullConsensusStateKey(clientID, height),
				Value: consState,
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
		{"ClientState", fmt.Sprintf("ClientState A: %v\nClientState B: %v", clientState, clientState)},
		{"ConsensusState", fmt.Sprintf("ConsensusState A: %v\nConsensusState B: %v", consState, consState)},
		{"other", ""},
	}

	for i, tt := range tests {
		i, tt := i, tt
		t.Run(tt.name, func(t *testing.T) {
			res, found := simulation.NewDecodeStore(kvPairs[i], kvPairs[i])
			if i == len(tests)-1 {
				require.False(t, found, string(kvPairs[i].Key))
				require.Empty(t, res, string(kvPairs[i].Key))
			} else {
				require.True(t, found, string(kvPairs[i].Key))
				require.Equal(t, tt.expectedLog, res, string(kvPairs[i].Key))
			}
		})
	}
}
