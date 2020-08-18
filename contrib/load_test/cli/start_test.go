// +build !integration

package cli

import (
	"encoding/json"
	"testing"

	"github.com/line/link/contrib/load_test/tests"
	"github.com/line/link/contrib/load_test/tests/mock"
	"github.com/line/link/contrib/load_test/types"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestStartCmd(t *testing.T) {
	// Given Mock Server
	server := mock.NewServer()
	defer server.Close()
	// Given Flags
	viper.Set(FlagMsgsPerTxLoadTest, tests.TestMsgsPerTxLoadTest)
	viper.Set(FlagMaxGasLoadTest, tests.TestMaxGasPrepare)
	viper.Set(FlagTPS, tests.TestTPS)
	viper.Set(FlagDuration, tests.TestDuration)
	viper.Set(FlagRampUpTime, tests.TestRampUpTime)
	viper.Set(FlagMaxWorkers, tests.TestMaxWorkers)
	viper.Set(FlagLoadTargetURL, server.URL)
	viper.Set(FlagChainID, tests.TestChainID)
	viper.Set(FlagCoinName, tests.TestCoinName)
	viper.Set(FlagTestnet, tests.TestNet)
	viper.Set(FlagOutputDir, "")
	// Given slave Flag
	slavesMap := make(map[string]types.Slave)
	slavesMap["slave1"] = types.NewSlave(server.URL, tests.TestMnemonic, types.QueryAccount, []string{})
	slavesMap["slave2"] = types.NewSlave(server.URL, tests.TestMnemonic2, types.TxSend, []string{})
	bytes, err := json.Marshal(slavesMap)
	require.NoError(t, err)
	viper.Set(FlagSlaves, string(bytes))

	// When
	require.NoError(t, start(StartCmd(), nil))

	// Then
	require.Equal(t, len(slavesMap), mock.GetCallCounter(server.URL).TargetLoadCallCount)
	require.Equal(t, len(slavesMap), mock.GetCallCounter(server.URL).TargetFireCallCount)
}
