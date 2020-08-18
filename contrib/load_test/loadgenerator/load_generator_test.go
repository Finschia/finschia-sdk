// +build !integration

package loadgenerator

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	"github.com/line/link/app"
	"github.com/line/link/contrib/load_test/scenario"
	"github.com/line/link/contrib/load_test/tests"
	"github.com/line/link/contrib/load_test/tests/mock"
	"github.com/line/link/contrib/load_test/types"
	"github.com/line/link/contrib/load_test/wallet"
	"github.com/stretchr/testify/require"
	vegeta "github.com/tsenart/vegeta/v12/lib"
)

func TestLoadGenerator_GenerateQueryTargets(t *testing.T) {
	// Given
	config := types.Config{
		MsgsPerTxLoadTest: tests.TestMsgsPerTxLoadTest,
		TPS:               tests.TestTPS,
		Duration:          tests.TestDuration,
		RampUpTime:        tests.TestRampUpTime,
		MaxWorkers:        tests.TestMaxWorkers,
		TargetURL:         tests.TestTargetURL,
		ChainID:           tests.TestChainID,
		CoinName:          tests.TestCoinName,
		Testnet:           tests.TestNet,
		Mnemonic:          tests.TestMnemonic,
	}
	scenarios := scenario.NewScenarios(config, nil, nil)
	loadGenerator := NewLoadGenerator()
	err := loadGenerator.ApplyConfig(config)
	require.NoError(t, err)

	// When
	require.NoError(t, loadGenerator.RunWithGoroutines(scenarios[types.QueryAccount].GenerateTarget))

	// Then
	require.Regexp(t, tests.TestTargetURL+"/auth/accounts/link1[a-z0-9]{38}$", loadGenerator.targets[0].URL)
	require.NotContains(t, loadGenerator.targets, vegeta.Target{})
}

func TestLoadGenerator_GenerateTxTargets(t *testing.T) {
	cdc := app.MakeCodec()
	// Given Mock Server
	server := mock.NewServer()
	defer server.Close()
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
		Testnet:           tests.TestNet,
		Mnemonic:          tests.TestMnemonic,
	}
	t.Log("Test success")
	{
		// Given Scenario
		scenarios := scenario.NewScenarios(config, nil, nil)
		// Given LoadGenerator
		loadGenerator := NewLoadGenerator()
		err := loadGenerator.ApplyConfig(config)
		require.NoError(t, err)

		// When
		require.NoError(t, loadGenerator.RunWithGoroutines(scenarios[types.TxSend].GenerateTarget))

		// Then
		require.Len(t, loadGenerator.targets, config.TPS*config.Duration)
		require.Equal(t, 1, loadGenerator.numTargetsPerUser)
		require.Equal(t, server.URL+"/txs", loadGenerator.targets[0].URL)
		require.NotContains(t, loadGenerator.targets, vegeta.Target{})
		require.Equal(t, tests.ExpectedNumTargets, mock.GetCallCounter(server.URL).QueryAccountCallCount)

		var req rest.BroadcastReq
		err = cdc.UnmarshalJSON(loadGenerator.targets[0].Body, &req)
		require.NoError(t, err)
		require.Len(t, req.Tx.Msgs, tests.TestMsgsPerTxLoadTest)
	}
	t.Log("Test with GenerateTarget that returns a numTargets greater than maxTargetsPerUser")
	{
		// Given invalid GenerateTarget
		invalidGenerateTarget := func(*wallet.KeyWallet, int) (*[]*vegeta.Target, int, error) {
			return nil, maxTargetsPerUser + 1, nil
		}
		// Given LoadGenerator
		loadGenerator := NewLoadGenerator()
		err := loadGenerator.ApplyConfig(config)
		require.NoError(t, err)

		// When
		err = loadGenerator.RunWithGoroutines(invalidGenerateTarget)

		// Then
		require.EqualError(t, err, types.ExceedMaxNumTargetsError{NumTargets: maxTargetsPerUser + 1,
			MaxTargetsPerUser: maxTargetsPerUser}.Error())
	}
	t.Log("Test with empty chain id")
	{
		// Given invalid config
		config.ChainID = ""
		// Given Scenario
		scenarios := scenario.NewScenarios(config, nil, nil)
		// Given LoadGenerator
		loadGenerator := NewLoadGenerator()
		err := loadGenerator.ApplyConfig(config)
		require.NoError(t, err)

		// When
		err = loadGenerator.RunWithGoroutines(scenarios[types.TxSend].GenerateTarget)
		require.EqualError(t, err, "chain ID required but not specified")

		// Then
		require.Len(t, loadGenerator.targets, 0)
	}
}

func TestLoadGenerator_Fire(t *testing.T) {
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
				Testnet:           tests.TestNet,
				Mnemonic:          tests.TestMnemonic,
			}
			scenario := scenario.NewScenarios(config, nil, nil)[types.TxSend]
			// Given LoadGenerator
			loadGenerator := NewLoadGenerator()
			err := loadGenerator.ApplyConfig(config)
			require.NoError(t, err)
			// Given Targets
			require.NoError(t, loadGenerator.RunWithGoroutines(scenario.GenerateTarget))

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
