package simulation_test

import (
	"fmt"
	"testing"

	"github.com/line/lbm-sdk/v2/store/types"
	"github.com/stretchr/testify/require"

	"github.com/line/lbm-sdk/v2/simapp"
	clienttypes "github.com/line/lbm-sdk/v2/x/ibc/core/02-client/types"
	connectiontypes "github.com/line/lbm-sdk/v2/x/ibc/core/03-connection/types"
	channeltypes "github.com/line/lbm-sdk/v2/x/ibc/core/04-channel/types"
	host "github.com/line/lbm-sdk/v2/x/ibc/core/24-host"
	"github.com/line/lbm-sdk/v2/x/ibc/core/simulation"
	ibctmtypes "github.com/line/lbm-sdk/v2/x/ibc/light-clients/07-tendermint/types"
)

func TestDecodeStore(t *testing.T) {
	cdc, _ := simapp.MakeCodecs()
	dec := simulation.NewDecodeStore(cdc)

	clientID := "clientidone"
	connectionID := "connectionidone"
	channelID := "channelidone"
	portID := "portidone"

	clientState := &ibctmtypes.ClientState{
		FrozenHeight: clienttypes.NewHeight(0, 10),
	}
	connection := connectiontypes.ConnectionEnd{
		ClientId: "clientidone",
		Versions: []*connectiontypes.Version{connectiontypes.NewVersion("1", nil)},
	}
	channel := channeltypes.Channel{
		State:   channeltypes.OPEN,
		Version: "1.0",
	}

	kvPairs := []types.KOPair{
			{
				Key:   host.FullClientStateKey(clientID),
				Value: clientState,
			},
			{
				Key:   host.ConnectionKey(connectionID),
				Value: &connection,
			},
			{
				Key:   host.ChannelKey(portID, channelID),
				Value: &channel,
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
		{"ConnectionEnd", fmt.Sprintf("ConnectionEnd A: %v\nConnectionEnd B: %v", connection, connection)},
		{"Channel", fmt.Sprintf("Channel A: %v\nChannel B: %v", channel, channel)},
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
