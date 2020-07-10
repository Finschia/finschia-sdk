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

func TestTxTokenScenario_GenerateStateSettingMsgs(t *testing.T) {
	// Given Mock Server
	server := mock.NewServer()
	defer server.Close()
	// Given Test Environments
	scenario, hdWallet, masterWallet := GivenTestEnvironments(t, server.URL, types.TxToken, nil)
	tokenScenario, ok := scenario.(*TxTokenScenario)
	require.True(t, ok)

	msgs, params, err := tokenScenario.GenerateStateSettingMsgs(masterWallet, hdWallet)
	require.NoError(t, err)

	require.Len(t, msgs, tests.TestTPS*tests.TestDuration*4)
	require.Equal(t, "send", msgs[tests.TestTPS*tests.TestDuration-1].Type())
	require.Equal(t, "grant_perm", msgs[tests.TestTPS*tests.TestDuration].Type())
	require.Equal(t, "9be17165", params["contract_id"])
}

func TestTxTokenScenario_GenerateTarget(t *testing.T) {
	// Given Mock Server
	server := mock.NewServer()
	defer server.Close()
	// Given Test Environments
	msgTypes := []string{"mint", "transfer_ft", "grant_perm", "revoke_perm", "modify_token", "burn"}
	scenario, _, keyWallet := GivenTestEnvironments(t, server.URL, types.TxToken,
		map[string]string{"contract_id": "9be17165"})
	tokenSenario, ok := scenario.(*TxTokenScenario)
	require.True(t, ok)

	// When
	targets, numTargets, err := tokenSenario.GenerateTarget(keyWallet, 0)
	require.NoError(t, err)

	// Then
	require.Equal(t, 1, numTargets)
	require.Equal(t, "POST", (*targets)[0].Method)
	require.Equal(t, fmt.Sprintf("%s%s", server.URL, TxURL), (*targets)[0].URL)
	// Then request is valid
	var req rest.BroadcastReq
	require.NoError(t, app.MakeCodec().UnmarshalJSON((*targets)[0].Body, &req))
	require.Equal(t, service.BroadcastSync, req.Mode)

	require.Len(t, req.Tx.Msgs, 6)
	for i, msg := range req.Tx.Msgs {
		require.Equal(t, msgTypes[i%6], msg.Type())
		require.NoError(t, msg.ValidateBasic())
	}
}
