// +build !integration

package cli

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/line/link/contrib/load_test/master"
	"github.com/line/link/contrib/load_test/tests"
	"github.com/line/link/contrib/load_test/tests/mock"
	"github.com/line/link/contrib/load_test/types"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestPrepareCmd(t *testing.T) {
	defer tests.RemoveFile(master.ParamsFileName)
	// Given Mock Server
	server := mock.NewServer()
	defer server.Close()
	// Given Flags
	viper.Set(FlagMasterMnemonic, tests.TestMasterMnemonic)
	viper.Set(FlagMsgsPerTxPrepare, tests.TestMsgsPerTxPrepare)
	viper.Set(FlagMaxGasPrepare, tests.TestMaxGasPrepare)
	viper.Set(FlagTPS, tests.TestTPS)
	viper.Set(FlagDuration, tests.TestDuration)
	viper.Set(FlagPrepareTargetURL, server.URL)
	viper.Set(FlagChainID, tests.TestChainID)
	viper.Set(FlagCoinName, tests.TestCoinName)
	viper.Set(FlagTestnet, tests.TestNet)
	viper.Set(FlagOutputDir, ".")

	testCases := []struct {
		scenario     string
		numPrepareTx int
		fileCheck    func(require.TestingT, string, ...interface{})
		params       string
	}{
		{types.QueryAccount, tests.TestNumPrepareRequest, require.NoFileExists, ""},
		{types.QueryBlock, tests.TestNumPrepareRequest, require.FileExists,
			fmt.Sprintf(`{"%s":{"height":"3"}}`, server.URL)},
		{types.TxSend, tests.TestNumPrepareRequest, require.NoFileExists, ""},
		{types.TxEmpty, tests.TestNumPrepareRequest, require.NoFileExists, ""},
		{types.TxToken, tests.TestNumPrepareRequest*4 + 1, require.FileExists,
			fmt.Sprintf(`{"%s":{"token_contract_id":"9be17165"}}`, server.URL)},
	}
	for _, tt := range testCases {
		t.Log(tt.scenario)
		// Given slave Flag
		slavesMap := make(map[string]types.Slave)
		slavesMap["slave1"] = types.NewSlave(server.URL, tests.TestMnemonic, tt.scenario, []string{})
		bytes, err := json.Marshal(slavesMap)
		require.NoError(t, err)
		viper.Set(FlagSlaves, string(bytes))

		// When
		require.NoError(t, prepare(PrepareCmd(), nil))

		// Then
		require.Equal(t, tt.numPrepareTx*len(slavesMap), mock.GetCallCounter(server.URL).BroadcastTxCallCount)
		tt.fileCheck(t, master.ParamsFileName)
		if _, err := os.Stat(master.ParamsFileName); err == nil {
			data, err := ioutil.ReadFile(master.ParamsFileName)
			require.NoError(t, err)
			require.Equal(t, tt.params, string(data))
		}
		// Clear Call Counter
		mock.ClearCallCounter(server.URL)
		tests.RemoveFile(master.ParamsFileName)
	}
}
