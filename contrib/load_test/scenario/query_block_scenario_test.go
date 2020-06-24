package scenario

import (
	"fmt"
	"testing"

	"github.com/line/link/contrib/load_test/tests"
	"github.com/line/link/contrib/load_test/tests/mock"
	"github.com/line/link/contrib/load_test/types"
	"github.com/stretchr/testify/require"
)

func TestQueryBlockScenario_GenerateStateSettingMsgs(t *testing.T) {
	// Given Mock Server
	server := mock.NewServer()
	defer server.Close()
	// Given Test Environments
	scenario, hdWallet, masterWallet := GivenTestEnvironments(t, server.URL, types.QueryBlock, nil)
	blockScenario, ok := scenario.(*QueryBlockScenario)
	require.True(t, ok)

	msgs, params, err := blockScenario.GenerateStateSettingMsgs(masterWallet, hdWallet)
	require.NoError(t, err)

	require.Len(t, msgs, tests.TestTPS*tests.TestDuration)
	require.Equal(t, "3", params["height"])
}

func TestQueryBlockScenario_GenerateTarget(t *testing.T) {
	// Given Mock Server
	server := mock.NewServer()
	defer server.Close()
	// Given Test Environments
	scenario, _, keyWallet := GivenTestEnvironments(t, server.URL, types.QueryBlock, map[string]string{"height": "3"})
	blockScenario, ok := scenario.(*QueryBlockScenario)
	require.True(t, ok)

	targets, numTargets, err := blockScenario.GenerateTarget(keyWallet, 0)
	require.NoError(t, err)

	require.Equal(t, 1, numTargets)
	require.Equal(t, "GET", (*targets)[0].Method)
	require.Equal(t, fmt.Sprintf("%s/blocks_with_tx_results/%d?fetchsize=%d", server.URL, 3, 3), (*targets)[0].URL)
}
