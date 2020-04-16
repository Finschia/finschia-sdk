// +build !integration

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
		PacerType:         types.ConstantPacer,
		TargetURL:         server.URL,
		ChainID:           tests.TestChainID,
		CoinName:          tests.TestCoinName,
	}
	t.Log("Test success")
	{
		// Given Slaves
		slaves := []types.Slave{
			types.NewSlave(server.URL, tests.TestMnemonic, types.QueryAccount),
			types.NewSlave(server.URL, tests.TestMnemonic2, types.TxSend),
		}
		// And Controller
		controller := NewController(slaves, config)

		// When
		require.NoError(t, controller.StartLoadTest())

		// Then
		for _, res := range controller.Results {
			require.Equal(t, "success", string(res))
		}
		require.Equal(t, len(slaves), mock.GetCallCounter(server.URL).TargetLoadCallCount)
		require.Equal(t, len(slaves), mock.GetCallCounter(server.URL).TargetFireCallCount)
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
		controller := NewController(slaves, config)

		// When
		err := controller.StartLoadTest()

		// Then
		require.Error(t, err)
	}
}
