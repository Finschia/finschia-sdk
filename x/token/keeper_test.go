package token

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"testing"
)

func TestGetAllTokens(t *testing.T) {
	input := setupTestInput(t)
	_, ctx, keeper, ak := input.cdc, input.ctx, input.keeper, input.ak

	addr1 := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	{
		acc := ak.NewAccountWithAddress(ctx, addr1)
		ak.SetAccount(ctx, acc)
	}
	require.NotNil(t, ak.GetAccount(ctx, addr1))
	require.Equal(t, uint64(0), ak.GetAccount(ctx, addr1).GetAccountNumber())

	{
		token := Token{Name: "name", Symbol: "link", Owner: addr1, Mintable: true}
		require.NoError(t, keeper.SetToken(ctx, token))
		tokens := keeper.GetAllTokens(ctx)
		require.Equal(t, "name", tokens[0].Name)
		require.Equal(t, "link", tokens[0].Symbol)
		require.Equal(t, true, tokens[0].Mintable)
		require.Equal(t, addr1, tokens[0].Owner)
	}
	{
		require.NoError(t, keeper.SetToken(ctx, Token{Name: "name", Symbol: "link1", Owner: addr1, Mintable: true}))
		require.NoError(t, keeper.SetToken(ctx, Token{Name: "name", Symbol: "link2", Owner: addr1, Mintable: true}))
		require.NoError(t, keeper.SetToken(ctx, Token{Name: "name", Symbol: "link3", Owner: addr1, Mintable: true}))
		require.NoError(t, keeper.SetToken(ctx, Token{Name: "name", Symbol: "link4", Owner: addr1, Mintable: true}))
		tokens := keeper.GetAllTokens(ctx)
		{
			require.Equal(t, "name", tokens[0].Name)
			require.Equal(t, "link", tokens[0].Symbol)
			require.Equal(t, true, tokens[0].Mintable)
			require.Equal(t, addr1, tokens[0].Owner)
		}
		{
			require.Equal(t, "link1", tokens[1].Symbol)
			require.Equal(t, "link2", tokens[2].Symbol)
			require.Equal(t, "link3", tokens[3].Symbol)
			require.Equal(t, "link4", tokens[4].Symbol)
		}
	}
}

func TestPublishTokenAndSendTokens(t *testing.T) {
	input := setupTestInput(t)
	_, ctx, keeper, ak := input.cdc, input.ctx, input.keeper, input.ak

	// Register account 1
	addr1 := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	{
		acc := ak.NewAccountWithAddress(ctx, addr1)
		ak.SetAccount(ctx, acc)
		require.NotNil(t, ak.GetAccount(ctx, addr1))
		require.Equal(t, uint64(0), ak.GetAccount(ctx, addr1).GetAccountNumber())
	}

	// Register account 2
	addr2 := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	{
		acc := ak.NewAccountWithAddress(ctx, addr2)
		ak.SetAccount(ctx, acc)
		require.NotNil(t, ak.GetAccount(ctx, addr2))
		require.Equal(t, uint64(1), ak.GetAccount(ctx, addr2).GetAccountNumber())
	}

	// Publish a token
	{
		token := Token{Name: "name", Symbol: "link", Owner: addr1, Mintable: true}
		require.NoError(t, keeper.SetToken(ctx, token))
		require.NoError(t, keeper.MintToken(ctx, sdk.NewCoin("link", sdk.NewInt(999)), addr1))

		token, err := keeper.GetToken(ctx, "link")
		require.NoError(t, err)
		require.Equal(t, "name", token.Name)
		require.Equal(t, "link", token.Symbol)
		require.Equal(t, true, token.Mintable)
		require.Equal(t, addr1, token.Owner)

		require.Equal(t, int64(999), keeper.supplyKeeper.GetSupply(ctx).GetTotal().AmountOf("link").Int64())
		require.Equal(t, int64(999), keeper.bankKeeper.GetCoins(ctx, addr1).AmountOf("link").Int64())
	}

	// Transfer Token
	{
		require.NoError(t, keeper.bankKeeper.SendCoins(ctx, addr1, addr2, sdk.NewCoins(sdk.NewCoin("link", sdk.NewInt(100)))))
		require.Equal(t, int64(899), keeper.bankKeeper.GetCoins(ctx, addr1).AmountOf("link").Int64())
		require.Equal(t, int64(100), keeper.bankKeeper.GetCoins(ctx, addr2).AmountOf("link").Int64())
		require.Equal(t, int64(999), keeper.supplyKeeper.GetSupply(ctx).GetTotal().AmountOf("link").Int64())
	}

	// Transfer Token again
	{
		require.NoError(t, keeper.bankKeeper.SendCoins(ctx, addr1, addr2, sdk.NewCoins(sdk.NewCoin("link", sdk.NewInt(100)))))
		require.Equal(t, int64(799), keeper.bankKeeper.GetCoins(ctx, addr1).AmountOf("link").Int64())
		require.Equal(t, int64(200), keeper.bankKeeper.GetCoins(ctx, addr2).AmountOf("link").Int64())
		require.Equal(t, int64(999), keeper.supplyKeeper.GetSupply(ctx).GetTotal().AmountOf("link").Int64())
	}

	// Mint Token
	{
		require.NoError(t, keeper.MintToken(ctx, sdk.NewCoin("link", sdk.NewInt(100)), addr1))
		require.Equal(t, int64(899), keeper.bankKeeper.GetCoins(ctx, addr1).AmountOf("link").Int64())
		require.Equal(t, int64(200), keeper.bankKeeper.GetCoins(ctx, addr2).AmountOf("link").Int64())
		require.Equal(t, int64(1099), keeper.supplyKeeper.GetSupply(ctx).GetTotal().AmountOf("link").Int64())
	}

	// Burn Token
	{
		require.NoError(t, keeper.BurnToken(ctx, sdk.NewCoin("link", sdk.NewInt(100)), addr1))
		require.Equal(t, int64(799), keeper.bankKeeper.GetCoins(ctx, addr1).AmountOf("link").Int64())
		require.Equal(t, int64(200), keeper.bankKeeper.GetCoins(ctx, addr2).AmountOf("link").Int64())
		require.Equal(t, int64(999), keeper.supplyKeeper.GetSupply(ctx).GetTotal().AmountOf("link").Int64())
	}

	// Burn Token again amount > has --> fail
	{
		require.Error(t, keeper.BurnToken(ctx, sdk.NewCoin("link", sdk.NewInt(800)), addr1))
		require.Equal(t, int64(799), keeper.bankKeeper.GetCoins(ctx, addr1).AmountOf("link").Int64())
		require.Equal(t, int64(200), keeper.bankKeeper.GetCoins(ctx, addr2).AmountOf("link").Int64())
		require.Equal(t, int64(999), keeper.supplyKeeper.GetSupply(ctx).GetTotal().AmountOf("link").Int64())
	}

}
