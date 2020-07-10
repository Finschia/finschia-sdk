// +build !integration

package service

import (
	"net/http"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/app"
	"github.com/line/link/contrib/load_test/tests"
	"github.com/line/link/contrib/load_test/tests/mock"
	"github.com/line/link/contrib/load_test/transaction"
	linktypes "github.com/line/link/types"
	"github.com/line/link/x/coin"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

const (
	TestAddress  = "link1muu5cza33kttadr5wylsqhfgxnwlcdrxls0wwn"
	TestNet      = false
	TestChainID  = "chain-id"
	TestCoinName = "link"
	TestHeight   = 3
	TestTxHash   = "D20985E8B70B54B7C79D37B8E214EE815EB8D9818CF793A20304678FFA2A4A92"
)

func TestLinkService_GetAccount(t *testing.T) {
	// Given Mock Server
	server := mock.NewServer()
	defer server.Close()
	// And LinkService
	sdk.GetConfig().SetBech32PrefixForAccount(linktypes.Bech32PrefixAcc(TestNet), linktypes.Bech32PrefixAccPub(TestNet))
	linkService := NewLinkService(&http.Client{}, app.MakeCodec(), server.URL)

	// When
	testAccount, err := linkService.GetAccount(TestAddress)
	require.NoError(t, err)

	// Then
	require.Equal(t, 1, mock.GetCallCounter(server.URL).QueryAccountCallCount)
	require.Equal(t, TestAddress, testAccount.Address.String())
	require.Equal(t, uint64(0x1), testAccount.AccountNumber)
	require.Equal(t, uint64(0xb), testAccount.Sequence)
}

func TestLinkService_GetBlock(t *testing.T) {
	// Given Mock Server
	server := mock.NewServer()
	defer server.Close()
	// And LinkService
	sdk.GetConfig().SetBech32PrefixForAccount(linktypes.Bech32PrefixAcc(TestNet), linktypes.Bech32PrefixAccPub(TestNet))
	linkService := NewLinkService(&http.Client{}, app.MakeCodec(), server.URL)

	// When
	testBlock, err := linkService.GetBlock(TestHeight)
	require.NoError(t, err)

	// Then
	require.Equal(t, 1, mock.GetCallCounter(server.URL).QueryBlockCallCount)
	require.Equal(t, int64(3), testBlock.Block.Height)
	require.Equal(t, "link", testBlock.Block.ChainID)
	require.Equal(t, 1, testBlock.BlockID.PartsHeader.Total)
}

func TestLinkService_GetLatestBlock(t *testing.T) {
	// Given Mock Server
	server := mock.NewServer()
	defer server.Close()
	// And LinkService
	sdk.GetConfig().SetBech32PrefixForAccount(linktypes.Bech32PrefixAcc(TestNet), linktypes.Bech32PrefixAccPub(TestNet))
	linkService := NewLinkService(&http.Client{}, app.MakeCodec(), server.URL)

	// When
	testBlock, err := linkService.GetLatestBlock()
	require.NoError(t, err)

	// Then
	require.Equal(t, 1, mock.GetCallCounter(server.URL).QueryBlockCallCount)
	require.Equal(t, int64(3), testBlock.Block.Height)
	require.Equal(t, "link", testBlock.Block.ChainID)
	require.Equal(t, 1, testBlock.BlockID.PartsHeader.Total)
}

func TestLinkService_GetBlocksWithTxResults(t *testing.T) {
	// Given Mock Server
	server := mock.NewServer()
	defer server.Close()
	// And LinkService
	sdk.GetConfig().SetBech32PrefixForAccount(linktypes.Bech32PrefixAcc(TestNet), linktypes.Bech32PrefixAccPub(TestNet))
	linkService := NewLinkService(&http.Client{}, app.MakeCodec(), server.URL)

	// When
	blocksWithTxResults, err := linkService.GetBlocksWithTxResults(3, 2)
	require.NoError(t, err)

	// Then
	require.Equal(t, 1, mock.GetCallCounter(server.URL).QueryBlocksWithTxResultsCallCount)
	require.Equal(t, int64(3), blocksWithTxResults[0].ResultBlock.Block.Height)
	require.Equal(t, int64(4), blocksWithTxResults[1].ResultBlock.Block.Height)
	require.Equal(t, int64(4), blocksWithTxResults[1].TxResponses[0].Height)
}

func TestLinkService_GetTx(t *testing.T) {
	// Given Mock Server
	server := mock.NewServer()
	defer server.Close()
	// And LinkService
	sdk.GetConfig().SetBech32PrefixForAccount(linktypes.Bech32PrefixAcc(TestNet), linktypes.Bech32PrefixAccPub(TestNet))
	linkService := NewLinkService(&http.Client{}, app.MakeCodec(), server.URL)

	// When
	res, err := linkService.GetTx(TestTxHash)
	require.NoError(t, err)

	// Then
	require.Equal(t, 1, mock.GetCallCounter(server.URL).QueryTxCallCount)
	require.Equal(t, int64(517257), res.Height)
	require.Equal(t, TestTxHash, res.TxHash)
	require.NotEmpty(t, res.RawLog)
	require.Equal(t, "send", res.Logs[0].Events[0].Attributes[0].Value)
	require.NotZero(t, res.GasWanted)
	require.NotZero(t, res.GasUsed)
	require.Equal(t, "send", res.Tx.GetMsgs()[0].Type())
	require.Equal(t, "2020-05-16T16:41:59Z", res.Timestamp)
}

func TestLinkService_BroadcastTx(t *testing.T) {
	// Given Mock Server
	server := mock.NewServer()
	defer server.Close()
	// And LinkService
	sdk.GetConfig().SetBech32PrefixForAccount(linktypes.Bech32PrefixAcc(TestNet), linktypes.Bech32PrefixAccPub(TestNet))
	linkService := NewLinkService(&http.Client{}, app.MakeCodec(), server.URL)
	// And MsgSend
	fromPrivateKey := secp256k1.GenPrivKey()
	from := fromPrivateKey.PubKey().Address().Bytes()
	to := secp256k1.GenPrivKey().PubKey().Address().Bytes()
	coins := sdk.NewCoins(sdk.NewCoin(TestCoinName, sdk.NewInt(10)))
	msgs := []sdk.Msg{coin.NewMsgSend(from, to, coins)}
	// And StdTx
	txBuilder := transaction.NewTxBuilder(tests.TestMaxGasPrepare).WithChainID(TestChainID)
	stdTx, err := txBuilder.BuildAndSign(fromPrivateKey, msgs)
	require.NoError(t, err)

	t.Log("Test block mode")
	{
		// When
		res, err := linkService.BroadcastTx(stdTx, "block")
		require.NoError(t, err)

		// Then
		require.Equal(t, 1, mock.GetCallCounter(server.URL).BroadcastTxCallCount)
		require.Equal(t, int64(517257), res.Height)
		require.Equal(t, TestTxHash, res.TxHash)
		require.NotEmpty(t, res.RawLog)
		require.Equal(t, "send", res.Logs[0].Events[0].Attributes[0].Value)
		require.NotZero(t, res.GasWanted)
		require.NotZero(t, res.GasUsed)
		require.Zero(t, res.Index)
		require.Zero(t, res.Code)
		require.Zero(t, res.Data)
		require.Zero(t, res.Info)
		require.Zero(t, res.Codespace)
		require.Zero(t, res.Tx)
		require.Zero(t, res.Timestamp)
	}
	t.Log("Test sync mode")
	{
		// When
		res, err := linkService.BroadcastTx(stdTx, "sync")
		require.NoError(t, err)

		// Then
		require.Equal(t, 2, mock.GetCallCounter(server.URL).BroadcastTxCallCount)
		require.Equal(t, int64(0), res.Height)
		require.Equal(t, TestTxHash, res.TxHash)
		require.Equal(t, "[]", res.RawLog)
		require.Zero(t, res.Logs)
		require.Zero(t, res.GasWanted)
		require.Zero(t, res.GasUsed)
		require.Zero(t, res.Index)
		require.Zero(t, res.Code)
		require.Zero(t, res.Data)
		require.Zero(t, res.Info)
		require.Zero(t, res.Codespace)
		require.Zero(t, res.Tx)
		require.Zero(t, res.Timestamp)
	}
	t.Log("Test async mode")
	{
		// When
		res, err := linkService.BroadcastTx(stdTx, "async")
		require.NoError(t, err)

		// Then
		require.Equal(t, 3, mock.GetCallCounter(server.URL).BroadcastTxCallCount)
		require.Equal(t, int64(0), res.Height)
		require.Equal(t, TestTxHash, res.TxHash)
		require.Zero(t, res.RawLog)
		require.Zero(t, res.Logs)
		require.Zero(t, res.GasWanted)
		require.Zero(t, res.GasUsed)
		require.Zero(t, res.Index)
		require.Zero(t, res.Code)
		require.Zero(t, res.Data)
		require.Zero(t, res.Info)
		require.Zero(t, res.Codespace)
		require.Zero(t, res.Tx)
		require.Zero(t, res.Timestamp)
	}
}
