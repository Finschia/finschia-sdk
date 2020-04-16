// +build !integration

package loadgenerator

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	"github.com/line/link/app"
	"github.com/line/link/contrib/load_test/tests"
	"github.com/line/link/contrib/load_test/transaction"
	"github.com/line/link/x/bank"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

func TestTargetBuilder_MakeQueryTarget(t *testing.T) {
	cdc := app.MakeCodec()
	testURL := "/test/url"
	targetBuilder := NewTargetBuilder(cdc, tests.TestTargetURL)

	target := targetBuilder.MakeQueryTarget(testURL)

	require.Equal(t, "GET", target.Method)
	require.Equal(t, tests.TestTargetURL+testURL, target.URL)
}

func TestTargetBuilder_MakeTxQuery(t *testing.T) {
	// Given TargetBuilder
	cdc := app.MakeCodec()
	targetBuilder := NewTargetBuilder(cdc, tests.TestTargetURL)
	// And MsgSend
	fromPrivateKey := secp256k1.GenPrivKey()
	from := fromPrivateKey.PubKey().Address().Bytes()
	to := secp256k1.GenPrivKey().PubKey().Address().Bytes()
	coins := sdk.NewCoins(sdk.NewCoin(tests.TestCoinName, sdk.NewInt(10)))
	msgs := []sdk.Msg{bank.NewMsgSend(from, to, coins)}
	// And TxBuilder
	txBuilder := transaction.NewTxBuilder().WithChainID(tests.TestChainID)
	stdTx, err := txBuilder.BuildAndSign(fromPrivateKey, msgs)
	require.NoError(t, err)

	// When
	target, err := targetBuilder.MakeTxTarget(stdTx, BroadcastMode)
	require.NoError(t, err)
	// And
	var broadcastReq rest.BroadcastReq
	err = cdc.UnmarshalJSON(target.Body, &broadcastReq)
	require.NoError(t, err)

	// Then
	require.Equal(t, "POST", target.Method)
	require.Equal(t, tests.TestTargetURL+TxURL, target.URL)
	require.Equal(t, BroadcastMode, broadcastReq.Mode)
	require.Equal(t, stdTx, broadcastReq.Tx)
}
