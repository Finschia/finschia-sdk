package scenario

import (
	"testing"

	"github.com/line/link/contrib/load_test/tests"
	"github.com/line/link/contrib/load_test/types"
	"github.com/stretchr/testify/require"
)

func TestNewScenarios(t *testing.T) {
	// Given Config
	config := types.Config{
		MsgsPerTxPrepare: tests.TestMsgsPerTxPrepare,
		TPS:              tests.TestTPS,
		Duration:         tests.TestDuration,
		TargetURL:        tests.TestTargetURL,
		ChainID:          tests.TestChainID,
		CoinName:         tests.TestCoinName,
	}
	scenarios := NewScenarios(config, nil, nil)

	require.IsType(t, &QueryAccountScenario{}, scenarios[types.QueryAccount])
	require.IsType(t, &QueryBlockScenario{}, scenarios[types.QueryBlock])
	require.IsType(t, &QueryCoinScenario{}, scenarios[types.QueryCoin])
	require.IsType(t, &TxSendScenario{}, scenarios[types.TxSend])
	require.IsType(t, &TxEmptyScenario{}, scenarios[types.TxEmpty])
	require.IsType(t, &TxMintNFTScenario{}, scenarios[types.TxMintNFT])
	require.IsType(t, &TxTransferFTScenario{}, scenarios[types.TxTransferFT])
	require.IsType(t, &TxTransferNFTScenario{}, scenarios[types.TxTransferNFT])
	require.IsType(t, &TxTokenScenario{}, scenarios[types.TxToken])
	require.IsType(t, &TxCollectionScenario{}, scenarios[types.TxCollection])
	require.IsType(t, &TxAndQueryAllScenario{}, scenarios[types.TxAndQueryAll])
}
