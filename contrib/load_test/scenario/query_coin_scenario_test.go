package scenario

import (
	"fmt"
	"testing"

	"github.com/line/link/contrib/load_test/tests/mock"
	"github.com/line/link/contrib/load_test/types"
	"github.com/stretchr/testify/require"
)

func TestQueryCoinScenario_GenerateTarget(t *testing.T) {
	// Given Mock Server
	server := mock.NewServer()
	defer server.Close()
	// Given Test Environments
	scenario, _, keyWallet := GivenTestEnvironments(t, server.URL, types.QueryCoin, nil, nil)
	coinScenario, ok := scenario.(*QueryCoinScenario)
	require.True(t, ok)

	targets, numTargets, err := coinScenario.GenerateTarget(keyWallet, 0)
	require.NoError(t, err)

	require.Equal(t, 1, numTargets)
	require.Equal(t, "GET", (*targets)[0].Method)
	require.Equal(t, fmt.Sprintf("%s/coin/cony", server.URL), (*targets)[0].URL)
}
