// +build !integration

package transaction

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/client/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/line/link/app"
	"github.com/line/link/contrib/load_test/tests"
	"github.com/line/link/x/coin"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

const (
	TestChainID   = "chain-id"
	TestCoinName  = "link"
	TestFees      = "1link"
	TestGasPrices = "1.0link"
	TestMemo      = "memo"
)

func TestNewTxBuilder(t *testing.T) {
	t.Log("Test with Accessible parameters")
	{
		testTxEncoder := utils.GetTxEncoder(app.MakeCodec())
		testSequence := uint64(3)
		testAccountNumber := uint64(2)

		txBuilder := NewTxBuilder(tests.TestMaxGasPrepare).WithTxEncoder(testTxEncoder).WithChainID(TestChainID).
			WithFees(TestFees).WithGasPrices(TestGasPrices).WithSequence(testSequence).WithMemo(TestMemo).
			WithAccountNumber(testAccountNumber)

		innerTxBuilder := txBuilder.txBuilder
		linkCoins := sdk.NewCoins(sdk.NewCoin("link", sdk.NewInt(1)))
		linkDecCoins := sdk.NewDecCoins(sdk.NewDecCoin("link", sdk.NewInt(1)))
		require.Equal(t, TestChainID, innerTxBuilder.ChainID())
		require.Equal(t, tests.TestMaxGasPrepare, innerTxBuilder.Gas())
		require.Equal(t, linkCoins, innerTxBuilder.Fees())
		require.Equal(t, linkDecCoins, innerTxBuilder.GasPrices())
		require.Equal(t, testSequence, innerTxBuilder.Sequence())
		require.Equal(t, TestMemo, innerTxBuilder.Memo())
		require.Equal(t, testAccountNumber, innerTxBuilder.AccountNumber())
	}
	t.Log("Test with Inaccessible parameters")
	{
		testKeybase := keys.NewInMemoryKeyBase()

		err := NewTxBuilder(tests.TestMaxGasPrepare).WithKeybase(testKeybase)

		require.EqualError(t, err, "Inaccessible Field Error: TxBuilderWithoutKeybase can not access keybase")
	}
}

func TestTxBuilderWithoutKeybase_BuildAndSign(t *testing.T) {
	// Given private key
	fromPrivateKey := secp256k1.GenPrivKey()
	// And MsgSend
	from := fromPrivateKey.PubKey().Address().Bytes()
	to := secp256k1.GenPrivKey().PubKey().Address().Bytes()
	coins := sdk.NewCoins(sdk.NewCoin(TestCoinName, sdk.NewInt(10)))
	msgs := []sdk.Msg{coin.NewMsgSend(from, to, coins)}
	// And TxBuilder
	txBuilder := NewTxBuilder(tests.TestMaxGasPrepare).WithChainID(TestChainID).WithGasPrices(TestGasPrices).
		WithMemo(TestMemo)

	// When
	stdTx, err := txBuilder.BuildAndSign(fromPrivateKey, msgs)

	// Then
	fees := sdk.NewCoins(sdk.NewCoin("link", sdk.NewInt(int64(tests.TestMaxGasPrepare))))
	require.NoError(t, err)
	require.Equal(t, msgs, stdTx.Msgs)
	require.Equal(t, auth.NewStdFee(tests.TestMaxGasPrepare, fees), stdTx.Fee)
	require.Equal(t, fromPrivateKey.PubKey(), stdTx.Signatures[0].PubKey)
	require.Len(t, stdTx.Signatures[0].Signature, 64)
	require.Equal(t, TestMemo, stdTx.Memo)
}
