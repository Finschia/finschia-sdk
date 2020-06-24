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

func TestTxAndQueryAllScenario_GenerateStateSettingMsgs(t *testing.T) {
	// Given Mock Server
	server := mock.NewServer()
	defer server.Close()
	// Given Test Environments
	scenario, hdWallet, masterWallet := GivenTestEnvironments(t, server.URL, types.TxAndQueryAll, nil)
	txAndQueryAllScenario, ok := scenario.(*TxAndQueryAllScenario)
	require.True(t, ok)

	msgs, params, err := txAndQueryAllScenario.GenerateStateSettingMsgs(masterWallet, hdWallet)
	require.NoError(t, err)

	require.Len(t, msgs, tests.TestTPS*tests.TestDuration*21)
	require.Equal(t, "send", msgs[tests.TestTPS*tests.TestDuration-1].Type())
	require.Equal(t, "grant_perm", msgs[tests.TestTPS*tests.TestDuration].Type())
	require.Equal(t, "grant_perm", msgs[8*tests.TestTPS*tests.TestDuration-1].Type())
	require.Equal(t, "mint_nft", msgs[8*tests.TestTPS*tests.TestDuration].Type())
	require.Equal(t, "9be17165", params["token_contract_id"])
	require.Equal(t, "678c146a", params["collection_contract_id"])
	require.Equal(t, "0000000100000000", params["ft_token_id"])
	require.Equal(t, "10000001", params["nft_token_type"])
	require.Equal(t, "16EFE7CF722157A57E03E947C6171B24A7FC3731E1A24FAE0D9168F80845407F", params["tx_hash"])
}

func TestTxAndQueryAllScenario_GenerateTarget(t *testing.T) {
	// Given Mock Server
	server := mock.NewServer()
	defer server.Close()
	// Given Test Environments
	msgTypes := []string{"send", "send", "issue_token", "mint", "mint", "transfer_ft", "transfer_ft", "modify_token",
		"modify_token", "burn", "burn", "create_collection", "approve_collection", "issue_ft", "mint_ft", "transfer_ft",
		"transfer_ft", "burn_ft", "modify_token", "issue_nft", "mint_nft", "mint_nft", "attach", "detach",
		"transfer_nft", "transfer_nft", "transfer_nft", "transfer_nft", "burn_nft"}
	scenario, _, keyWallet := GivenTestEnvironments(t, server.URL, types.TxAndQueryAll,
		map[string]string{
			"token_contract_id":      "9be17165",
			"nft_token_type":         "10000001",
			"ft_token_id":            "0000000100000000",
			"tx_hash":                "16EFE7CF722157A57E03E947C6171B24A7FC3731E1A24FAE0D9168F80845407F",
			"collection_contract_id": "678c146a",
		})
	txAndQueryAllScenario, ok := scenario.(*TxAndQueryAllScenario)
	require.True(t, ok)

	// When
	targets, numTargets, err := txAndQueryAllScenario.GenerateTarget(keyWallet, 0)
	require.NoError(t, err)

	// Then
	require.Equal(t, 98, numTargets)
	require.Equal(t, "POST", (*targets)[0].Method)
	for i := 1; i < numTargets; i++ {
		require.Equal(t, "GET", (*targets)[i].Method)
	}
	require.Equal(t, fmt.Sprintf("%s%s", server.URL, TxURL), (*targets)[0].URL)
	// Then tx target is valid
	var req rest.BroadcastReq
	require.NoError(t, app.MakeCodec().UnmarshalJSON((*targets)[0].Body, &req))
	require.Equal(t, service.BroadcastSync, req.Mode)

	require.Len(t, req.Tx.Msgs, len(msgTypes))
	for i, msg := range req.Tx.Msgs {
		require.Equal(t, msgTypes[i%len(msgTypes)], msg.Type())
		require.NoError(t, msg.ValidateBasic())
	}
}
