// +build !integration

package master

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/contrib/load_test/tests"
	"github.com/line/link/contrib/load_test/tests/mock"
	"github.com/line/link/contrib/load_test/types"
	linktypes "github.com/line/link/types"
	"github.com/stretchr/testify/require"
)

func TestStateSetter_PrepareTestState(t *testing.T) {
	defer tests.RemoveFile(ParamsFileName)
	sdk.GetConfig().SetBech32PrefixForAccount(linktypes.Bech32PrefixAcc(tests.TestNet), linktypes.Bech32PrefixAccPub(tests.TestNet))
	// Given Mock Server
	server := mock.NewServer()
	defer server.Close()
	// Given Config
	config := types.Config{
		MsgsPerTxPrepare:  tests.TestMsgsPerTxPrepare,
		MsgsPerTxLoadTest: tests.TestMsgsPerTxLoadTest,
		TPS:               tests.TestTPS,
		Duration:          tests.TestDuration,
		TargetURL:         server.URL,
		ChainID:           tests.TestChainID,
		CoinName:          tests.TestCoinName,
	}

	t.Log("Test success")
	{
		testCases := []struct {
			scenario     string
			numPrepareTx int
			fileCheck    func(require.TestingT, string, ...interface{})
			params       string
		}{
			{types.QueryAccount,
				tests.TestNumPrepareRequest,
				require.NoFileExists,
				"",
			},
			{types.QueryBlock,
				tests.TestNumPrepareRequest,
				require.FileExists,
				fmt.Sprintf(`{"%s":{"height":"3"}}`, server.URL),
			},
			{types.TxSend,
				tests.TestNumPrepareRequest,
				require.NoFileExists,
				"",
			},
			{types.TxEmpty,
				tests.TestNumPrepareRequest,
				require.NoFileExists,
				"",
			},
			{types.TxToken,
				tests.GetNumPrepareTx(tests.ExpectedNumTargets*4, tests.TestMsgsPerTxPrepare) + 1,
				require.FileExists,
				fmt.Sprintf(`{"%s":{"token_contract_id":"9be17165"}}`, server.URL),
			},
			{types.TxCollection,
				tests.GetNumPrepareTx(tests.ExpectedNumTargets*(4+2*tests.TestMsgsPerTxLoadTest), tests.TestMsgsPerTxPrepare) + 3,
				require.FileExists,
				fmt.Sprintf(`{"%s":{"collection_contract_id":"678c146a","ft_token_id":"0000000100000000","nft_token_type":"10000001","num_nft_per_user":"6"}}`, server.URL),
			},
		}
		for _, tt := range testCases {
			t.Log(tt.scenario)
			// Given StateSetter
			ss, err := NewStateSetter(tests.TestMasterMnemonic, config)
			require.NoError(t, err)
			// Given Slaves
			slaves := []types.Slave{
				types.NewSlave(server.URL, tests.TestMnemonic, tt.scenario, []string{}),
			}

			require.NoError(t, ss.PrepareTestState(slaves, "."))
			require.Equal(t, tt.numPrepareTx*len(slaves), mock.GetCallCounter(server.URL).BroadcastTxCallCount)
			tt.fileCheck(t, ParamsFileName)
			if _, err := os.Stat(ParamsFileName); err == nil {
				data, err := ioutil.ReadFile(ParamsFileName)
				require.NoError(t, err)
				require.Equal(t, tt.params, string(data))
			}
			// Clear Call Counter
			mock.ClearCallCounter(server.URL)
			tests.RemoveFile(ParamsFileName)
		}
	}
	t.Log("Test with invalid master mnemonic")
	{
		_, err := NewStateSetter(tests.InvalidMnemonic, config)
		require.EqualError(t, err, "Invalid master mnemonic: invalid mnemonic")
		require.NoFileExists(t, ParamsFileName)
	}
	t.Log("Test with invalid mnemonic")
	{
		// Given StateSetter
		ss, err := NewStateSetter(tests.TestMasterMnemonic, config)
		require.NoError(t, err)
		// Given Slaves
		slaves := []types.Slave{
			types.NewSlave(server.URL, tests.InvalidMnemonic, types.QueryAccount, []string{}),
			types.NewSlave(server.URL, tests.InvalidMnemonic, types.TxSend, []string{}),
		}

		require.EqualError(t, ss.PrepareTestState(slaves, "."), "Invalid mnemonic: invalid mnemonic")
		require.NoFileExists(t, ParamsFileName)
	}
	t.Log("Test with empty chain id")
	{
		// Given StateSetter
		config.ChainID = ""
		ss, err := NewStateSetter(tests.TestMasterMnemonic, config)
		require.NoError(t, err)
		// Given Slaves
		slaves := []types.Slave{
			types.NewSlave(server.URL, tests.TestMnemonic, types.QueryAccount, []string{}),
			types.NewSlave(server.URL, tests.TestMnemonic2, types.TxSend, []string{}),
		}

		require.EqualError(t, ss.PrepareTestState(slaves, "."), "chain ID required but not specified")
		require.NoFileExists(t, ParamsFileName)
	}
}
