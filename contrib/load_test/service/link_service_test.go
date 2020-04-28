// +build !integration

package service

import (
	"net/http"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/app"
	"github.com/line/link/contrib/load_test/tests/mock"
	"github.com/line/link/contrib/load_test/transaction"
	linktypes "github.com/line/link/types"
	"github.com/line/link/x/bank"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

const (
	TestAddress  = "link1muu5cza33kttadr5wylsqhfgxnwlcdrxls0wwn"
	TestNet      = false
	TestChainID  = "chain-id"
	TestCoinName = "link"
	TestHeight   = 3
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
	msgs := []sdk.Msg{bank.NewMsgSend(from, to, coins)}
	// And StdTx
	txBuilder := transaction.NewTxBuilder().WithChainID(TestChainID)
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
		require.Equal(t, "D20985E8B70B54B7C79D37B8E214EE815EB8D9818CF793A20304678FFA2A4A92", res.TxHash)
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
		require.Equal(t, "D20985E8B70B54B7C79D37B8E214EE815EB8D9818CF793A20304678FFA2A4A92", res.TxHash)
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
		require.Equal(t, "D20985E8B70B54B7C79D37B8E214EE815EB8D9818CF793A20304678FFA2A4A92", res.TxHash)
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
