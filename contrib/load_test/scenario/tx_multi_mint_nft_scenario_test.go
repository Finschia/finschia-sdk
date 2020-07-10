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
	"github.com/line/link/x/collection"
	"github.com/stretchr/testify/require"
)

func TestTxMultiMintNFTScenario_GenerateTarget(t *testing.T) {
	// Given Mock Server
	server := mock.NewServer()
	defer server.Close()
	// Given Test Environments
	scenario, _, keyWallet := GivenTestEnvironments(t, server.URL, types.TxMultiMintNFT,
		map[string]string{"contract_id": "678c146a", "nft_token_type": "10000001"})
	mintNFTScenario, ok := scenario.(*TxMultiMintNFTScenario)
	require.True(t, ok)

	// When
	targets, numTargets, err := mintNFTScenario.GenerateTarget(keyWallet, 0)
	require.NoError(t, err)

	// Then
	require.Equal(t, 1, numTargets)
	require.Equal(t, "POST", (*targets)[0].Method)
	require.Equal(t, fmt.Sprintf("%s%s", server.URL, TxURL), (*targets)[0].URL)
	// Then request is valid
	var req rest.BroadcastReq
	require.NoError(t, app.MakeCodec().UnmarshalJSON((*targets)[0].Body, &req))
	require.Equal(t, service.BroadcastSync, req.Mode)

	require.Len(t, req.Tx.Msgs, tests.TestMsgsPerTxLoadTest)
	for _, msg := range req.Tx.Msgs {
		require.Equal(t, "mint_nft", msg.Type())
		require.NoError(t, msg.ValidateBasic())
		msgMintNFT, ok := msg.(collection.MsgMintNFT)
		require.True(t, ok)
		require.Len(t, msgMintNFT.MintNFTParams, 5)
	}
}
