// +build !integration

package cli

import (
	"encoding/json"
	"testing"

	"github.com/line/link/contrib/load_test/tests"
	"github.com/line/link/contrib/load_test/types"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestGetSlaves(t *testing.T) {
	// Given
	slavesMap := make(map[string]types.Slave)
	slavesMap["slave1"] = types.NewSlave(tests.TestTargetURL, tests.TestMnemonic, types.QueryAccount)
	slavesMap["slave2"] = types.NewSlave(tests.TestTargetURL, tests.TestMnemonic2, types.TxSend)
	bytes, err := json.Marshal(slavesMap)
	require.NoError(t, err)
	viper.Set(FlagSlaves, string(bytes))

	// When
	slaves := getSlaves()

	// Then
	require.ElementsMatch(t, []types.Slave{slavesMap["slave1"], slavesMap["slave2"]}, slaves)
}
