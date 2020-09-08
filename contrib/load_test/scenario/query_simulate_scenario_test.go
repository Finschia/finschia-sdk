package scenario

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	"github.com/line/link/app"
	"github.com/line/link/contrib/load_test/tests"
	"github.com/line/link/contrib/load_test/tests/mock"
	"github.com/line/link/contrib/load_test/types"
	"github.com/stretchr/testify/require"
)

func TestQuerySimulateScenario_GenerateStateSettingMsgs(t *testing.T) {
	// Given Mock Server
	server := mock.NewServer()
	defer server.Close()
	// Given Test Environments
	scenario, hdWallet, masterWallet := GivenTestEnvironments(t, server.URL, types.QuerySimulate, nil, nil)
	querySimulateScenario, ok := scenario.(*QuerySimulateScenario)
	require.True(t, ok)

	var cases = []struct {
		msgType        string
		numMsgsPerUser int
	}{
		{"MsgSend", 1},
		{"MsgMintNFT", 2},
		{"MsgTransferFT", 2},
		{"MsgTransferNFT", 1 + tests.TestMsgsPerTxLoadTest},
	}

	for _, tt := range cases {
		// When
		msgs, _, err := querySimulateScenario.GenerateStateSettingMsgs(masterWallet, hdWallet, []string{tt.msgType})
		require.NoError(t, err)

		// Then
		require.Len(t, msgs, tests.TestTPS*tests.TestDuration*tt.numMsgsPerUser)
	}
}

func TestQuerySimulateScenario_GenerateTarget(t *testing.T) {
	// Given Mock Server
	server := mock.NewServer()
	defer server.Close()
	// Given Test Environments
	scenario, _, keyWallet := GivenTestEnvironments(t, server.URL, types.QuerySimulate,
		map[string]string{"collection_contract_id": "678c146a", "ft_token_id": "0000000100000000",
			"nft_token_type": "10000001", "num_nft_per_user": fmt.Sprintf("%d", tests.TestMsgsPerTxLoadTest)},
		nil)
	querySimulateScenario, ok := scenario.(*QuerySimulateScenario)
	require.True(t, ok)

	var cases = []struct {
		scenarioParams []string
		msgtype        string
	}{
		// {[]string{"MsgSend"}, "send"},
		{[]string{"MsgMintNFT", "1"}, "mint_nft"},
		// {[]string{"MsgTransferFT"}, "transfer_ft"},
		// {[]string{"MsgTransferNFT"}, "transfer_nft"},
	}

	for _, tt := range cases {
		// Given
		querySimulateScenario.Info.scenarioParams = tt.scenarioParams

		// When
		targets, numTargets, err := querySimulateScenario.GenerateTarget(keyWallet, 0)
		require.NoError(t, err)

		// Then
		require.Equal(t, 1, numTargets)
		require.Equal(t, fmt.Sprintf("%s%s", server.URL, QuerySimulateURL), (*targets)[0].URL)
		// Then request is valid
		var req rest.BroadcastReq
		require.NoError(t, app.MakeCodec().UnmarshalJSON((*targets)[0].Body, &req))
		require.Len(t, req.Tx.Msgs, querySimulateScenario.config.MsgsPerTxLoadTest)
		for _, msg := range req.Tx.Msgs {
			require.Equal(t, tt.msgtype, msg.Type())
			require.NoError(t, msg.ValidateBasic())
		}
	}
}
