// +build !integration

package loadgenerator

import (
	"net/http"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	"github.com/line/link/app"
	"github.com/line/link/contrib/load_test/tests"
	"github.com/line/link/contrib/load_test/tests/mock"
	"github.com/line/link/contrib/load_test/types"
	linktypes "github.com/line/link/types"
	"github.com/stretchr/testify/require"
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

func TestLoadGenerator_GenerateQueryTargets(t *testing.T) {
	// Given
	sdk.GetConfig().SetBech32PrefixForAccount(linktypes.Bech32PrefixAcc(tests.TestNet), linktypes.Bech32PrefixAccPub(tests.TestNet))
	config := types.Config{
		MsgsPerTxLoadTest: tests.TestMsgsPerTxLoadTest,
		TPS:               tests.TestTPS,
		Duration:          tests.TestDuration,
		MaxWorkers:        tests.TestMaxWorkers,
		PacerType:         types.ConstantPacer,
		TargetURL:         tests.TestTargetURL,
		ChainID:           tests.TestChainID,
		CoinName:          tests.TestCoinName,
		Mnemonic:          tests.TestMnemonic,
	}
	loadGenerator := NewLoadGenerator()
	err := loadGenerator.ApplyConfig(config)
	require.NoError(t, err)

	// When
	require.NoError(t, loadGenerator.RunWithGoroutines(loadGenerator.GenerateAccountQueryTarget))

	// Then
	require.Regexp(t, tests.TestTargetURL+"/auth/accounts/link1[a-z0-9]{38}$", loadGenerator.targets[0].URL)
	require.NotContains(t, loadGenerator.targets, vegeta.Target{})
}

func TestLoadGenerator_GenerateTxTargets(t *testing.T) {
	cdc := app.MakeCodec()
	// Given Mock Server
	server := mock.NewServer()
	defer server.Close()
	sdk.GetConfig().SetBech32PrefixForAccount(linktypes.Bech32PrefixAcc(tests.TestNet), linktypes.Bech32PrefixAccPub(tests.TestNet))
	// Given Config
	config := types.Config{
		MsgsPerTxLoadTest: tests.TestMsgsPerTxLoadTest,
		TPS:               tests.TestTPS,
		Duration:          tests.TestDuration,
		MaxWorkers:        tests.TestMaxWorkers,
		PacerType:         types.ConstantPacer,
		TargetURL:         server.URL,
		ChainID:           tests.TestChainID,
		CoinName:          tests.TestCoinName,
		Mnemonic:          tests.TestMnemonic,
	}

	t.Log("Test success")
	{
		// Given LoadGenerator
		loadGenerator := NewLoadGenerator()
		err := loadGenerator.ApplyConfig(config)
		require.NoError(t, err)

		// When
		require.NoError(t, loadGenerator.RunWithGoroutines(loadGenerator.GenerateMsgSendTxTarget))

		// Then
		require.Equal(t, server.URL+"/txs", loadGenerator.targets[0].URL)
		require.NotContains(t, loadGenerator.targets, vegeta.Target{})
		require.Equal(t, tests.ExpectedNumTargetsConstant, mock.GetCallCounter(server.URL).QueryAccountCallCount)

		var req rest.BroadcastReq
		err = cdc.UnmarshalJSON(loadGenerator.targets[0].Body, &req)
		require.NoError(t, err)
		require.Len(t, req.Tx.Msgs, tests.TestMsgsPerTxLoadTest)
	}
	t.Log("Test with empty chain id")
	{
		// Given invalid config
		config.ChainID = ""
		// Given LoadGenerator
		loadGenerator := NewLoadGenerator()
		err := loadGenerator.ApplyConfig(config)
		require.NoError(t, err)

		// When
		err = loadGenerator.RunWithGoroutines(loadGenerator.GenerateMsgSendTxTarget)
		require.EqualError(t, err, "chain ID required but not specified")

		// Then
		require.Equal(t, vegeta.Target{}, loadGenerator.targets[0])
	}
}

func TestLoadGenerator_Fire(t *testing.T) {
	sdk.GetConfig().SetBech32PrefixForAccount(linktypes.Bech32PrefixAcc(tests.TestNet), linktypes.Bech32PrefixAccPub(tests.TestNet))
	// Given Mock Server
	server := mock.NewServer()
	defer server.Close()
	// Given config
	config := types.Config{
		MsgsPerTxLoadTest: tests.TestMsgsPerTxLoadTest,
		TPS:               tests.TestTPS,
		Duration:          tests.TestDuration,
		MaxWorkers:        tests.TestMaxWorkers,
		PacerType:         types.ConstantPacer,
		TargetURL:         server.URL,
		ChainID:           tests.TestChainID,
		CoinName:          tests.TestCoinName,
		Mnemonic:          tests.TestMnemonic,
	}

	t.Log("Test with constant pacer")
	{
		// Given LoadGenerator
		loadGenerator := NewLoadGenerator()
		err := loadGenerator.ApplyConfig(config)
		require.NoError(t, err)
		// Given Targets
		require.NoError(t, loadGenerator.RunWithGoroutines(loadGenerator.GenerateMsgSendTxTarget))

		// When
		for res := range loadGenerator.Fire(tests.TestLoadGeneratorURL) {
			// Then
			require.Equal(t, "LINK v2 load test: "+tests.TestLoadGeneratorURL, res.Attack)
			require.Equal(t, uint16(http.StatusOK), res.Code)
			require.NotEmpty(t, res.Body)
		}
		require.Equal(t, tests.ExpectedNumTargetsConstant, mock.GetCallCounter(server.URL).BroadcastTxCallCount)
	}
	t.Log("Test with linear pacer")
	{
		// Given config
		config.PacerType = types.LinearPacer
		// Given LoadGenerator
		loadGenerator := NewLoadGenerator()
		err := loadGenerator.ApplyConfig(config)
		require.NoError(t, err)
		// Given Targets
		require.NoError(t, loadGenerator.RunWithGoroutines(loadGenerator.GenerateMsgSendTxTarget))
		// Clear Call Counter
		mock.ClearCallCounter(server.URL)

		// When
		for res := range loadGenerator.Fire(tests.TestLoadGeneratorURL) {
			// Then
			require.Equal(t, "LINK v2 load test: "+tests.TestLoadGeneratorURL, res.Attack)
			require.Equal(t, uint16(http.StatusOK), res.Code)
			require.NotEmpty(t, res.Body)
		}
		require.Equal(t, tests.ExpectedNumTargetsLinear, mock.GetCallCounter(server.URL).BroadcastTxCallCount)
	}
}
