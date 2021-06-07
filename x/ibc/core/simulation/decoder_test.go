package simulation_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/line/lbm-sdk/v2/store/types"
	types2 "github.com/line/lbm-sdk/v2/types"
	"github.com/line/lbm-sdk/v2/types/kv"
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
	height := clienttypes.NewHeight(0, 10)

	clientState := &ibctmtypes.ClientState{
		FrozenHeight: height,
	}
	connection := connectiontypes.ConnectionEnd{
		ClientId: "clientidone",
		Versions: []*connectiontypes.Version{connectiontypes.NewVersion("1", nil)},
	}
	channel := channeltypes.Channel{
		State:   channeltypes.OPEN,
		Version: "1.0",
	}
	consState := &ibctmtypes.ConsensusState{
		Timestamp: time.Now().UTC(),
	}

	cs, err := cdc.MarshalInterface(clientState)
	require.NoError(t, err)

	cs2, err2 := cdc.MarshalInterface(consState)
	require.NoError(t, err2)

	paths := connectiontypes.ClientPaths{
		Paths: []string{connectionID},
	}

	bz := []byte{0x1, 0x2, 0x3}
	kvPairs := kv.Pairs{
		Pairs: []kv.Pair{
			{
				Key:   host.FullClientStateKey(clientID),
				Value: cs,
			},
			{
				Key:   host.FullConsensusStateKey(clientID, height),
				Value: cs2,
			},
			{
				Key:   host.ConnectionKey(connectionID),
				Value: cdc.MustMarshalBinaryBare(&connection),
			},
			{
				Key:   host.ChannelKey(portID, channelID),
				Value: cdc.MustMarshalBinaryBare(&channel),
			},
			{
				Key:   host.ClientConnectionsKey(connection.ClientId),
				Value: cdc.MustMarshalBinaryBare(&paths),
			},
			{
				Key:   host.NextSequenceSendKey(portID, channelID),
				Value: types2.Uint64ToBigEndian(1),
			},
			{
				Key:   host.NextSequenceRecvKey(portID, channelID),
				Value: types2.Uint64ToBigEndian(1),
			},
			{
				Key:   host.NextSequenceAckKey(portID, channelID),
				Value: types2.Uint64ToBigEndian(1),
			},
			{
				Key:   host.PacketCommitmentKey(portID, channelID, 1),
				Value: bz,
			},
			{
				Key:   host.PacketAcknowledgementKey(portID, channelID, 1),
				Value: bz,
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
		{"ClientState", fmt.Sprintf("ClientState A: %v\nClientState B: %v", clientState, clientState)},
		{"ConsensusState", fmt.Sprintf("ConsensusState A: %v\nConsensusState B: %v", consState, consState)},
		{"ConnectionEnd", fmt.Sprintf("ConnectionEnd A: %v\nConnectionEnd B: %v", connection, connection)},
		{"Channel", fmt.Sprintf("Channel A: %v\nChannel B: %v", channel, channel)},
		{"ClientPaths", fmt.Sprintf("ClientPaths A: %v\nClientPaths B: %v", paths, paths)},
		{"NextSeqSend", "NextSeqSend A: 1\nNextSeqSend B: 1"},
		{"NextSeqRecv", "NextSeqRecv A: 1\nNextSeqRecv B: 1"},
		{"NextSeqAck", "NextSeqAck A: 1\nNextSeqAck B: 1"},
		{"CommitmentHash", fmt.Sprintf("CommitmentHash A: %X\nCommitmentHash B: %X", bz, bz)},
		{"AckHash", fmt.Sprintf("AckHash A: %X\nAckHash B: %X", bz, bz)},
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
			pair := types.KOPair{
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
