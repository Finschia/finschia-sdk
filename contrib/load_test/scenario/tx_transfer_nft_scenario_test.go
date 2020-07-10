package scenario

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	"github.com/line/link/app"
	"github.com/line/link/contrib/load_test/service"
	"github.com/line/link/contrib/load_test/tests"
	"github.com/line/link/contrib/load_test/tests/mock"
	"github.com/line/link/contrib/load_test/types"
	"github.com/stretchr/testify/require"
)

func TestTxTransferNFTScenario_GenerateStateSettingMsgs(t *testing.T) {
	// Given Mock Server
	server := mock.NewServer()
	defer server.Close()
	// Given Test Environments
	scenario, hdWallet, masterWallet := GivenTestEnvironments(t, server.URL, types.TxTransferNFT, nil)
	transferNFTScenario, ok := scenario.(*TxTransferNFTScenario)
	require.True(t, ok)

	msgs, params, err := transferNFTScenario.GenerateStateSettingMsgs(masterWallet, hdWallet)
	require.NoError(t, err)

	require.Len(t, msgs, tests.TestTPS*tests.TestDuration*(1+tests.TestMsgsPerTxLoadTest))
	require.Equal(t, "send", msgs[tests.TestTPS*tests.TestDuration-1].Type())
	require.Equal(t, "mint_nft", msgs[tests.TestTPS*tests.TestDuration].Type())
	require.Equal(t, "678c146a", params["contract_id"])
	require.Equal(t, "10000001", params["nft_token_type"])
}

//nolint: dupl
func TestTxTransferNFTScenario_GenerateTarget(t *testing.T) {
	// Given Mock Server
	server := mock.NewServer()
	defer server.Close()
	// Given Test Environments
	scenario, _, keyWallet := GivenTestEnvironments(t, server.URL, types.TxTransferNFT,
		map[string]string{"contract_id": "678c146a", "nft_token_type": "10000001"})
	transferNFTScenario, ok := scenario.(*TxTransferNFTScenario)
	require.True(t, ok)

	// When
	targets, numTargets, err := transferNFTScenario.GenerateTarget(keyWallet, 0)
	require.NoError(t, err)

	// Then
	require.Equal(t, 1, numTargets)
	require.Equal(t, "POST", (*targets)[0].Method)
	require.Equal(t, fmt.Sprintf("%s%s", server.URL, TxURL), (*targets)[0].URL)
	// Then request is valid
	var req rest.BroadcastReq
	require.NoError(t, app.MakeCodec().UnmarshalJSON((*targets)[0].Body, &req))
	require.Equal(t, service.BroadcastSync, req.Mode)

	for _, msg := range req.Tx.Msgs {
		require.Equal(t, "transfer_nft", msg.Type())
		require.NoError(t, msg.ValidateBasic())
	}
}
