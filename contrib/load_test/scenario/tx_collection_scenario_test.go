package scenario

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	"github.com/line/link/app"
	"github.com/line/link/contrib/load_test/tests"
	"github.com/line/link/contrib/load_test/tests/mock"
	"github.com/line/link/contrib/load_test/types"
	"github.com/line/link/x/collection"
	"github.com/stretchr/testify/require"
)

func TestTxCollectionScenario_GenerateStateSettingMsgs(t *testing.T) {
	// Given Mock Server
	server := mock.NewServer()
	defer server.Close()
	// Given Test Environments
	scenario, hdWallet, masterWallet := GivenTestEnvironments(t, server.URL, types.TxCollection, nil, nil)
	collectionScenario, ok := scenario.(*TxCollectionScenario)
	require.True(t, ok)

	msgs, params, err := collectionScenario.GenerateStateSettingMsgs(masterWallet, hdWallet, []string{})
	require.NoError(t, err)

	require.Len(t, msgs, tests.TestTPS*tests.TestDuration*(4+2*tests.TestMsgsPerTxLoadTest))
	require.Equal(t, "send", msgs[tests.TestTPS*tests.TestDuration-1].Type())
	require.Equal(t, "grant_perm", msgs[tests.TestTPS*tests.TestDuration].Type())
	require.Equal(t, "grant_perm", msgs[4*tests.TestTPS*tests.TestDuration-1].Type())
	require.Equal(t, "mint_nft", msgs[4*tests.TestTPS*tests.TestDuration].Type())
	require.Equal(t, "678c146a", params["collection_contract_id"])
	require.Equal(t, "0000000100000000", params["ft_token_id"])
	require.Equal(t, "10000001", params["nft_token_type"])
}

func TestTxCollectionScenario_GenerateTarget(t *testing.T) {
	// Given Mock Server
	server := mock.NewServer()
	defer server.Close()
	// Given Test Environments
	msgTypes := []string{"mint_ft", "transfer_ft", "modify_token", "burn_ft", "attach", "detach", "transfer_nft",
		"burn_nft"}
	scenario, _, keyWallet := GivenTestEnvironments(t, server.URL, types.TxCollection,
		map[string]string{"collection_contract_id": "678c146a", "ft_token_id": "0000000100000000",
			"nft_token_type": "10000001", "num_nft_per_user": fmt.Sprintf("%d", 2*tests.TestMsgsPerTxLoadTest)}, nil)
	collectionScenario, ok := scenario.(*TxCollectionScenario)
	require.True(t, ok)

	// When
	targets, numTargets, err := collectionScenario.GenerateTarget(keyWallet, 0)
	require.NoError(t, err)

	// Then
	require.Equal(t, 1, numTargets)
	// Then request is valid
	var req rest.BroadcastReq
	require.NoError(t, app.MakeCodec().UnmarshalJSON((*targets)[0].Body, &req))

	require.Len(t, req.Tx.Msgs, len(msgTypes)*tests.TestMsgsPerTxLoadTest)
	for i, msg := range req.Tx.Msgs {
		require.Equal(t, msgTypes[i%8], msg.Type())
		require.NoError(t, msg.ValidateBasic())
	}
	msgAttach, ok := req.Tx.Msgs[4].(collection.MsgAttach)
	require.True(t, ok)
	require.Equal(t, "1000000100000001", msgAttach.ToTokenID)
	require.Equal(t, "1000000100000002", msgAttach.TokenID)
}
