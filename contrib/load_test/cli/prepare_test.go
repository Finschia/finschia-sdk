// +build !integration

package cli

import (
	"encoding/json"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/contrib/load_test/tests"
	"github.com/line/link/contrib/load_test/tests/mock"
	"github.com/line/link/contrib/load_test/types"
	linktypes "github.com/line/link/types"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestPrepareCmd(t *testing.T) {
	sdk.GetConfig().SetBech32PrefixForAccount(linktypes.Bech32PrefixAcc(tests.TestNet), linktypes.Bech32PrefixAccPub(tests.TestNet))
	// Given Mock Server
	server := mock.NewServer()
	defer server.Close()
	// Given Flags
	viper.Set(FlagMasterMnemonic, tests.TestMasterMnemonic)
	viper.Set(FlagMsgsPerTxPrepare, tests.TestMsgsPerTxPrepare)
	viper.Set(FlagTPS, tests.TestTPS)
	viper.Set(FlagDuration, tests.TestDuration)
	viper.Set(FlagTargetURL, server.URL)
	viper.Set(FlagChainID, tests.TestChainID)
	viper.Set(FlagCoinName, tests.TestCoinName)
	// Given slave Flag
	slavesMap := make(map[string]types.Slave)
	slavesMap["slave1"] = types.NewSlave(server.URL, tests.TestMnemonic, types.QueryAccount)
	slavesMap["slave2"] = types.NewSlave(server.URL, tests.TestMnemonic2, types.TxSend)
	bytes, err := json.Marshal(slavesMap)
	require.NoError(t, err)
	viper.Set(FlagSlaves, string(bytes))

	// When
	require.NoError(t, prepare(PrepareCmd(), nil))

	// Then
	require.Equal(t, tests.TestNumPrepareRequest*len(slavesMap), mock.GetCallCounter(server.URL).BroadcastTxCallCount)
}
