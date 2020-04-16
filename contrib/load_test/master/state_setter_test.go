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

func TestStateSetter_PrepareTestState(t *testing.T) {
	sdk.GetConfig().SetBech32PrefixForAccount(linktypes.Bech32PrefixAcc(tests.TestNet), linktypes.Bech32PrefixAccPub(tests.TestNet))
	// Given Mock Server
	server := mock.NewServer()
	defer server.Close()
	// Given Config
	config := types.Config{
		MsgsPerTxPrepare: tests.TestMsgsPerTxPrepare,
		TPS:              tests.TestTPS,
		Duration:         tests.TestDuration,
		TargetURL:        server.URL,
		ChainID:          tests.TestChainID,
		CoinName:         tests.TestCoinName,
	}
	// Given StateSetter
	ss, err := NewStateSetter(tests.TestMasterMnemonic, config)
	require.NoError(t, err)
	// Given Slaves
	slaves := []types.Slave{
		types.NewSlave(server.URL, tests.TestMnemonic, types.QueryAccount),
		types.NewSlave(server.URL, tests.TestMnemonic2, types.TxSend),
	}

	require.NoError(t, ss.PrepareTestState(slaves))
	require.Equal(t, tests.TestNumPrepareRequest*len(slaves), mock.GetCallCounter(server.URL).BroadcastTxCallCount)
}

func TestStateSetter_RegisterAccounts(t *testing.T) {
	sdk.GetConfig().SetBech32PrefixForAccount(linktypes.Bech32PrefixAcc(tests.TestNet), linktypes.Bech32PrefixAccPub(tests.TestNet))
	// Given Mock Server
	server := mock.NewServer()
	defer server.Close()
	// Given Config
	config := types.Config{
		MsgsPerTxPrepare: tests.TestMsgsPerTxPrepare,
		TPS:              tests.TestTPS,
		Duration:         tests.TestDuration,
		TargetURL:        server.URL,
		ChainID:          tests.TestChainID,
		CoinName:         tests.TestCoinName,
	}
	t.Log("Test success")
	{
		ss, err := NewStateSetter(tests.TestMasterMnemonic, config)
		require.NoError(t, err)

		require.NoError(t, ss.RegisterAccounts(tests.TestMnemonic))
		require.Equal(t, tests.TestNumPrepareRequest, mock.GetCallCounter(server.URL).BroadcastTxCallCount)
	}
	t.Log("Test with invalid mnemonic")
	{
		ss, err := NewStateSetter(tests.TestMasterMnemonic, config)
		require.NoError(t, err)

		err = ss.RegisterAccounts(tests.InvalidMnemonic)

		require.EqualError(t, err, "Invalid mnemonic: invalid mnemonic")
	}
	t.Log("Test with empty chain id")
	{
		config.ChainID = ""
		ss, err := NewStateSetter(tests.TestMasterMnemonic, config)
		require.NoError(t, err)

		err = ss.RegisterAccounts(tests.TestMnemonic)

		require.EqualError(t, err, "chain ID required but not specified")
	}
}
