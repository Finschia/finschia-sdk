package simulation_test

import (
	"fmt"
	"testing"

	types2 "github.com/line/lbm-sdk/v2/store/types"
	"github.com/stretchr/testify/require"

	"github.com/line/lbm-sdk/v2/x/ibc/core/03-connection/simulation"
	"github.com/line/lbm-sdk/v2/x/ibc/core/03-connection/types"
	host "github.com/line/lbm-sdk/v2/x/ibc/core/24-host"
)

func TestDecodeStore(t *testing.T) {
	connectionID := "connectionidone"

	connection := types.ConnectionEnd{
		ClientId: "clientidone",
		Versions: types.ExportedVersionsToProto(types.GetCompatibleVersions()),
	}

	paths := types.ClientPaths{
		Paths: []string{connectionID},
	}

	kvPairs := []types2.KOPair{
			{
				Key:   host.ClientConnectionsKey(connection.ClientId),
				Value: &paths,
			},
			{
				Key:   host.ConnectionKey(connectionID),
				Value: &connection,
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
		{"ClientPaths", fmt.Sprintf("ClientPaths A: %v\nClientPaths B: %v", paths, paths)},
		{"ConnectionEnd", fmt.Sprintf("ConnectionEnd A: %v\nConnectionEnd B: %v", connection, connection)},
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
