package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	simappparams "github.com/Finschia/finschia-sdk/simapp/params"
	datest "github.com/Finschia/finschia-sdk/x/or/da/testutil"
	"github.com/Finschia/finschia-sdk/x/or/da/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
		CCList: dummyCCList(),
	}

	k, ctx, _ := datest.DaKeeper(t, simappparams.MakeTestEncodingConfig())
	k.InitGenesis(ctx, genesisState)
	got := k.ExportGenesis(ctx)
	require.NotNil(t, got)
	require.Equal(t, genesisState.Params, got.Params)

	expected := genesisState.CCList[0]
	var actual types.CC
	if got.CCList[0].RollupName == expected.RollupName {
		actual = got.CCList[0]
	} else {
		actual = got.CCList[1]
	}
	require.Equal(t, expected.RollupName, actual.RollupName)
	require.Equal(t, expected.CCState, actual.CCState)
	require.Equal(t, expected.QueueTxState, actual.QueueTxState)
	require.ElementsMatch(t, expected.History, actual.History)
	require.ElementsMatch(t, expected.QueueList, actual.QueueList)
	require.ElementsMatch(t, expected.L2BatchMap, actual.L2BatchMap)
}

func dummyCCList() []types.CC {
	t := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	return []types.CC{
		{
			RollupName: "dummy",
			CCState: types.CCState{
				Base:             5,
				Height:           10,
				L1Height:         1,
				Timestamp:        t,
				ProcessedL2Block: 7,
			},
			History: []types.CCRef{
				{
					BatchRoot:   []byte("dummy batch root"),
					TotalFrames: 10,
					BatchSize:   100,
					TxHash:      []byte("dummy tx hash"),
					MsgIndex:    3,
				},
				{
					BatchRoot:   []byte("dummy batch root2"),
					TotalFrames: 11,
					BatchSize:   23,
					TxHash:      []byte("dummy tx hash2"),
					MsgIndex:    1,
				},
			},
			QueueTxState: types.QueueTxState{
				ProcessedQueueIndex: 20,
				NextQueueIndex:      40,
			},
			QueueList: []types.L1ToL2Queue{
				{
					Timestamp: t,
					L1Height:  10,
					Status:    types.QUEUE_TX_PENDING,
					Txraw:     []byte("dummy txraw"),
				},
				{
					Timestamp: t,
					L1Height:  13,
					Status:    types.QUEUE_TX_EXPIRED,
					Txraw:     []byte("dummy txraw2"),
				},
			},
			L2BatchMap: []types.L2BatchMap{
				{
					L2Height: 20,
					BatchIdx: 5,
				},
				{
					L2Height: 44,
					BatchIdx: 10,
				},
			},
		},
		{
			RollupName: "dummy2",
			CCState: types.CCState{
				Base:             5,
				Height:           10,
				L1Height:         1,
				Timestamp:        t,
				ProcessedL2Block: 7,
			},
			History: []types.CCRef{
				{
					BatchRoot:   []byte("dummy batch root"),
					TotalFrames: 10,
					BatchSize:   100,
					TxHash:      []byte("dummy tx hash"),
					MsgIndex:    3,
				},
				{
					BatchRoot:   []byte("dummy batch root2"),
					TotalFrames: 11,
					BatchSize:   23,
					TxHash:      []byte("dummy tx hash2"),
					MsgIndex:    1,
				},
			},
			QueueTxState: types.QueueTxState{
				ProcessedQueueIndex: 20,
				NextQueueIndex:      40,
			},
			QueueList: []types.L1ToL2Queue{
				{
					Timestamp: t,
					L1Height:  10,
					Status:    types.QUEUE_TX_PENDING,
					Txraw:     []byte("dummy txraw"),
				},
			},
			L2BatchMap: []types.L2BatchMap{
				{
					L2Height: 20,
					BatchIdx: 5,
				},
			},
		},
	}
}
