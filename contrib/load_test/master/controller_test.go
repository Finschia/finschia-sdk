// +build !integration

package master

import (
	"net/http"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/contrib/load_test/tests"
	"github.com/line/link/contrib/load_test/tests/mock"
	"github.com/line/link/contrib/load_test/types"
	linktypes "github.com/line/link/types"
	"github.com/stretchr/testify/require"
)

func TestController_StartLoadTest(t *testing.T) {
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
	t.Log("Test success")
	{
		testCases := []struct {
			scenario     string
			numPrepareTx int
			params       map[string]map[string]string
		}{
			{types.QueryAccount,
				tests.TestNumPrepareRequest,
				nil,
			},
			{types.QueryBlock,
				tests.TestNumPrepareRequest,
				map[string]map[string]string{server.URL: {"height": "3"}},
			},
			{types.TxSend,
				tests.TestNumPrepareRequest,
				nil,
			},
			{types.TxEmpty,
				tests.TestNumPrepareRequest,
				nil,
			},
			{types.TxToken,
				tests.GetNumPrepareTx(tests.ExpectedNumTargets*4, tests.TestMsgsPerTxPrepare) + 1,
				map[string]map[string]string{server.URL: {"contractID": "9be17165"}},
			},
			{types.TxCollection,
				tests.GetNumPrepareTx(tests.ExpectedNumTargets*6, tests.TestMsgsPerTxPrepare) + 3,
				map[string]map[string]string{server.URL: {"contract_id": "678c146a", "ft_token_id": "0000000100000000",
					"nft_token_type": "10000001"}},
			},
		}
		for _, tt := range testCases {
			t.Log(tt.scenario)
			// Given Slaves
			slaves := []types.Slave{
				types.NewSlave(server.URL, tests.TestMnemonic, tt.scenario),
			}
			// Given Controller

			controller := NewController(slaves, config, tt.params)

			// When
			require.NoError(t, controller.StartLoadTest())

			// Then
			for _, res := range controller.Results {
				require.Equal(t, uint16(http.StatusOK), res[0].Code)
				require.Equal(t, "LINK v2 load test: localhost:8000", res[0].Attack)
			}
			require.Equal(t, len(slaves), mock.GetCallCounter(server.URL).TargetLoadCallCount)
			require.Equal(t, len(slaves), mock.GetCallCounter(server.URL).TargetFireCallCount)

			// Clear Call Counter
			mock.ClearCallCounter(server.URL)
		}
	}
	t.Log("Test with invalid slave url")
	{
		// Given Slaves
		invalidURL := "http://invalid_url.com"
		slaves := []types.Slave{
			types.NewSlave(invalidURL, tests.TestMnemonic, types.QueryAccount),
			types.NewSlave(server.URL, tests.TestMnemonic2, types.TxSend),
		}

		// And Controller
		controller := NewController(slaves, config, nil)

		// When
		err := controller.StartLoadTest()

		// Then
		require.Error(t, err)
	}
}
