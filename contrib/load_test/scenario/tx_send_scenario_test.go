package scenario

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	"github.com/line/link/app"
	"github.com/line/link/contrib/load_test/tests/mock"
	"github.com/line/link/contrib/load_test/types"
	"github.com/line/link/x/coin"
	"github.com/stretchr/testify/require"
)

func TestTxSendScenario_GenerateTarget(t *testing.T) {
	// Given Mock Server
	server := mock.NewServer()
	defer server.Close()
	// Given Test Environments
	scenario, _, keyWallet := GivenTestEnvironments(t, server.URL, types.TxSend, nil, nil)
	sendScenario, ok := scenario.(*TxSendScenario)
	require.True(t, ok)

	// When
	targets, numTargets, err := sendScenario.GenerateTarget(keyWallet, 0)
	require.NoError(t, err)

	// Then
	require.Equal(t, 1, numTargets)
	// Then request is valid
	var req rest.BroadcastReq
	require.NoError(t, app.MakeCodec().UnmarshalJSON((*targets)[0].Body, &req))
	require.Len(t, req.Tx.Msgs, sendScenario.config.MsgsPerTxLoadTest)
	for _, msg := range req.Tx.Msgs {
		require.Equal(t, "send", msg.Type())
		require.NoError(t, msg.ValidateBasic())
		msgSend, ok := msg.(coin.MsgSend)
		require.True(t, ok)
		require.Equal(t, sendScenario.config.CoinName, msgSend.Amount[0].Denom)
	}
}
