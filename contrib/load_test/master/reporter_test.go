package master

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/contrib/load_test/tests"
	"github.com/line/link/contrib/load_test/tests/mock"
	"github.com/line/link/contrib/load_test/types"
	linktypes "github.com/line/link/types"
	"github.com/stretchr/testify/require"
)

func TestReporter_Report(t *testing.T) {
	sdk.GetConfig().SetBech32PrefixForAccount(linktypes.Bech32PrefixAcc(tests.TestNet), linktypes.Bech32PrefixAccPub(tests.TestNet))
	// Given Mock Server
	server := mock.NewServer()
	defer server.Close()
	// Given config
	config := types.Config{
		MsgsPerTxLoadTest: tests.TestMsgsPerTxLoadTest,
		TPS:               tests.TestTPS,
		Duration:          tests.TestDuration,
		RampUpTime:        tests.TestRampUpTime,
		TargetURL:         server.URL,
		ChainID:           tests.TestChainID,
		CoinName:          tests.TestCoinName,
	}
	// Given Slaves
	slaves := []types.Slave{
		types.NewSlave(server.URL, tests.TestMnemonic, types.QueryAccount),
		types.NewSlave(server.URL, tests.TestMnemonic2, types.TxSend),
	}
	// Given Results
	controller := NewController(slaves, config)
	require.NoError(t, controller.StartLoadTest())
	// Given Reporter
	startHeight := int64(1)
	endHeight := int64(3)
	reporter := NewReporter(controller.Results, config, startHeight, endHeight)

	// When
	require.NoError(t, reporter.Report(""))

	// Then
	require.Len(t, reporter.metrics.blockMetrics.tps, int(endHeight-startHeight+1))
	require.Len(t, reporter.metrics.blockMetrics.height, int(endHeight-startHeight+1))
	require.Len(t, reporter.metrics.blockMetrics.blockInterval, int(endHeight-startHeight+1))
	require.Len(t, reporter.metrics.blockMetrics.numTxs, int(endHeight-startHeight+1))
	require.Equal(t, uint64(len(slaves)), reporter.metrics.Requests)
	require.Len(t, reporter.metrics.latencies, len(slaves))
	require.Len(t, reporter.metrics.timeStamps, len(slaves))
	require.Equal(t, float64(1), reporter.metrics.Success)
}
