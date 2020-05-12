// +build !integration

package loadgenerator

import (
	"fmt"
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
		RampUpTime:        tests.TestRampUpTime,
		MaxWorkers:        tests.TestMaxWorkers,
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

func TestLoadGenerator_GenerateCustomQueryTarget(t *testing.T) {
	// Given
	sdk.GetConfig().SetBech32PrefixForAccount(linktypes.Bech32PrefixAcc(tests.TestNet), linktypes.Bech32PrefixAccPub(tests.TestNet))
	config := types.Config{
		MsgsPerTxLoadTest: tests.TestMsgsPerTxLoadTest,
		TPS:               tests.TestTPS,
		Duration:          tests.TestDuration,
		RampUpTime:        tests.TestRampUpTime,
		MaxWorkers:        tests.TestMaxWorkers,
		TargetURL:         tests.TestTargetURL,
		ChainID:           tests.TestChainID,
		CoinName:          tests.TestCoinName,
		Mnemonic:          tests.TestMnemonic,
	}
	loadGenerator := NewLoadGenerator()
	err := loadGenerator.ApplyConfig(config)
	require.NoError(t, err)
	// Given custom url
	loadGenerator.customURL = tests.TestCustomURL

	// When
	require.NoError(t, loadGenerator.RunWithGoroutines(loadGenerator.GenerateCustomQueryTarget))

	// Then
	require.Equal(t, tests.TestTargetURL+tests.TestCustomURL, loadGenerator.targets[0].URL)
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
		RampUpTime:        tests.TestRampUpTime,
		MaxWorkers:        tests.TestMaxWorkers,
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
		require.Equal(t, tests.ExpectedNumTargets, mock.GetCallCounter(server.URL).QueryAccountCallCount)

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

	var testCase = []struct {
		duration   int
		rampUpTime int
	}{
		{3, 0},
		{3, 2},

		{4, 0},
		{4, 2},
		{4, 4},
	}

	for i, tt := range testCase {
		tt := tt
		t.Run(fmt.Sprintf("Test #%d", i), func(t *testing.T) {
			t.Parallel()
			// Given Mock Server
			server := mock.NewServer()
			defer server.Close()
			// Given config
			config := types.Config{
				MsgsPerTxLoadTest: tests.TestMsgsPerTxLoadTest,
				TPS:               tests.TestTPS,
				Duration:          tt.duration,
				RampUpTime:        tt.rampUpTime,
				MaxWorkers:        tests.TestMaxWorkers,
				TargetURL:         server.URL,
				ChainID:           tests.TestChainID,
				CoinName:          tests.TestCoinName,
				Mnemonic:          tests.TestMnemonic,
			}
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

			expectedAttackCount := (tt.duration-tt.rampUpTime/2)*tests.TestTPS + tt.duration
			require.InDelta(t, expectedAttackCount, mock.GetCallCounter(server.URL).BroadcastTxCallCount,
				tests.TestMaxAttackDifference)
		})
	}
}
