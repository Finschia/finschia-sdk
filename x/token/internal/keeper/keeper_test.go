package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/token/internal/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"testing"
)

func TestGetAllTokens(t *testing.T) {
	input := SetupTestInput(t)
	_, ctx, keeper, ak := input.Cdc, input.Ctx, input.Keeper, input.Ak

	addr1 := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	{
		acc := ak.NewAccountWithAddress(ctx, addr1)
		ak.SetAccount(ctx, acc)
	}
	require.NotNil(t, ak.GetAccount(ctx, addr1))
	require.Equal(t, uint64(0), ak.GetAccount(ctx, addr1).GetAccountNumber())

	{
		token := types.NewFT("name", "link", sdk.NewInt(8), true)
		require.NoError(t, keeper.IssueFT(ctx, token, sdk.NewInt(10), addr1))
		tokens := keeper.GetAllTokens(ctx)
		require.Equal(t, "name", tokens[0].Name)
		require.Equal(t, "link", tokens[0].Symbol)
		require.Equal(t, int64(8), tokens[0].Decimals.Int64())
		require.Equal(t, true, tokens[0].Mintable)
		require.Equal(t, int64(10), keeper.accountKeeper.GetAccount(ctx, addr1).GetCoins().AmountOf("link").Int64())
	}
	{
		require.NoError(t, keeper.IssueFT(ctx, types.NewFT("name", "link1", sdk.NewInt(8), true), sdk.NewInt(100), addr1))
		require.NoError(t, keeper.IssueFT(ctx, types.NewFT("name", "link2", sdk.NewInt(8), true), sdk.NewInt(200), addr1))
		require.NoError(t, keeper.IssueFT(ctx, types.NewFT("name", "link3", sdk.NewInt(8), true), sdk.NewInt(300), addr1))
		require.NoError(t, keeper.IssueFT(ctx, types.NewFT("name", "link4", sdk.NewInt(8), true), sdk.NewInt(400), addr1))
		tokens := keeper.GetAllTokens(ctx)
		{
			require.Equal(t, "name", tokens[0].Name)
			require.Equal(t, "link", tokens[0].Symbol)
			require.Equal(t, true, tokens[0].Mintable)
		}
		{
			require.Equal(t, "link1", tokens[1].Symbol)
			require.Equal(t, "link2", tokens[2].Symbol)
			require.Equal(t, "link3", tokens[3].Symbol)
			require.Equal(t, "link4", tokens[4].Symbol)
		}

		{
			require.Equal(t, int64(100), keeper.accountKeeper.GetAccount(ctx, addr1).GetCoins().AmountOf("link1").Int64())
			require.Equal(t, int64(200), keeper.accountKeeper.GetAccount(ctx, addr1).GetCoins().AmountOf("link2").Int64())
			require.Equal(t, int64(300), keeper.accountKeeper.GetAccount(ctx, addr1).GetCoins().AmountOf("link3").Int64())
			require.Equal(t, int64(400), keeper.accountKeeper.GetAccount(ctx, addr1).GetCoins().AmountOf("link4").Int64())
		}
	}

	keeper.Logger(ctx).Info("test", "test", "success")
}

func TestIssueTokenAndSendTokens(t *testing.T) {
	input := SetupTestInput(t)
	_, ctx, keeper, ak, bk := input.Cdc, input.Ctx, input.Keeper, input.Ak, input.Bk

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

	// Issue a token
	{
		token := types.NewFT("name", "link", sdk.NewInt(8), true)
		require.NoError(t, keeper.IssueFT(ctx, token, sdk.NewInt(900), addr1))

		token, err := keeper.GetToken(ctx, "link")
		require.NoError(t, err)
		require.Equal(t, "name", token.Name)
		require.Equal(t, "link", token.Symbol)
		require.Equal(t, int64(8), token.Decimals.Int64())
		require.Equal(t, true, token.Mintable)
		require.Equal(t, int64(900), keeper.accountKeeper.GetAccount(ctx, addr1).GetCoins().AmountOf("link").Int64())

		require.NoError(t, keeper.MintTokens(ctx, sdk.NewCoins(sdk.NewCoin("link", sdk.NewInt(99))), addr1))

		require.Equal(t, int64(999), keeper.supplyKeeper.GetSupply(ctx).GetTotal().AmountOf("link").Int64())
		require.Equal(t, int64(999), bk.GetCoins(ctx, addr1).AmountOf("link").Int64())
	}
	// Issue a token again FAIL
	{
		token := types.NewFT("name", "link", sdk.NewInt(8), true)
		require.Error(t, keeper.IssueFT(ctx, token, sdk.NewInt(900), addr1))
	}

	// Transfer Token
	{
		require.NoError(t, bk.SendCoins(ctx, addr1, addr2, sdk.NewCoins(sdk.NewCoin("link", sdk.NewInt(100)))))
		require.Equal(t, int64(899), bk.GetCoins(ctx, addr1).AmountOf("link").Int64())
		require.Equal(t, int64(100), bk.GetCoins(ctx, addr2).AmountOf("link").Int64())
		require.Equal(t, int64(999), keeper.supplyKeeper.GetSupply(ctx).GetTotal().AmountOf("link").Int64())
	}

	// Transfer Token again
	{
		require.NoError(t, bk.SendCoins(ctx, addr1, addr2, sdk.NewCoins(sdk.NewCoin("link", sdk.NewInt(100)))))
		require.Equal(t, int64(799), bk.GetCoins(ctx, addr1).AmountOf("link").Int64())
		require.Equal(t, int64(200), bk.GetCoins(ctx, addr2).AmountOf("link").Int64())
		require.Equal(t, int64(999), keeper.supplyKeeper.GetSupply(ctx).GetTotal().AmountOf("link").Int64())
	}

	// Mint Token
	{
		require.NoError(t, keeper.MintTokens(ctx, sdk.NewCoins(sdk.NewCoin("link", sdk.NewInt(100))), addr1))
		require.Equal(t, int64(899), bk.GetCoins(ctx, addr1).AmountOf("link").Int64())
		require.Equal(t, int64(200), bk.GetCoins(ctx, addr2).AmountOf("link").Int64())
		require.Equal(t, int64(1099), keeper.supplyKeeper.GetSupply(ctx).GetTotal().AmountOf("link").Int64())
	}

	// Burn Token
	{
		require.NoError(t, keeper.BurnTokens(ctx, sdk.NewCoins(sdk.NewCoin("link", sdk.NewInt(100))), addr1))
		require.Equal(t, int64(799), bk.GetCoins(ctx, addr1).AmountOf("link").Int64())
		require.Equal(t, int64(200), bk.GetCoins(ctx, addr2).AmountOf("link").Int64())
		require.Equal(t, int64(999), keeper.supplyKeeper.GetSupply(ctx).GetTotal().AmountOf("link").Int64())
	}

	// Burn Token again amount > has --> fail
	{
		require.Error(t, keeper.BurnTokens(ctx, sdk.NewCoins(sdk.NewCoin("link", sdk.NewInt(800))), addr1))
		require.Equal(t, int64(799), bk.GetCoins(ctx, addr1).AmountOf("link").Int64())
		require.Equal(t, int64(200), bk.GetCoins(ctx, addr2).AmountOf("link").Int64())
		require.Equal(t, int64(999), keeper.supplyKeeper.GetSupply(ctx).GetTotal().AmountOf("link").Int64())
	}

}

func TestIssueNFTAndSendTokens(t *testing.T) {
	input := SetupTestInput(t)
	_, ctx, keeper, ak, bk := input.Cdc, input.Ctx, input.Keeper, input.Ak, input.Bk

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

	// Issue a nft
	{
		token := types.NewNFT("name", "linknft", "metadata")
		require.NoError(t, keeper.IssueNFT(ctx, token, addr1))
	}
	// Issue a nft again FAIL
	{
		token := types.NewNFT("name", "linknft", "metadata")
		require.Error(t, keeper.IssueNFT(ctx, token, addr1))
	}

	// Test
	{
		token, err := keeper.GetToken(ctx, "linknft")
		require.NoError(t, err)
		require.Equal(t, "name", token.Name)
		require.Equal(t, "linknft", token.Symbol)
		require.Equal(t, int64(0), token.Decimals.Int64())
		require.Equal(t, false, token.Mintable)
		require.Equal(t, int64(1), bk.GetCoins(ctx, addr1).AmountOf("linknft").Int64())

		require.Error(t, keeper.MintTokens(ctx, sdk.NewCoins(sdk.NewCoin("linknft", sdk.NewInt(99))), addr1))

		require.Equal(t, int64(1), keeper.supplyKeeper.GetSupply(ctx).GetTotal().AmountOf("linknft").Int64())
		require.Equal(t, int64(1), bk.GetCoins(ctx, addr1).AmountOf("linknft").Int64())
	}

	// Send insufficient

	{
		require.Error(t, bk.SendCoins(ctx, addr1, addr2, sdk.NewCoins(sdk.NewCoin("linknft", sdk.NewInt(10)))))
		require.Equal(t, int64(1), bk.GetCoins(ctx, addr1).AmountOf("linknft").Int64())
		require.Equal(t, int64(0), bk.GetCoins(ctx, addr2).AmountOf("linknft").Int64())
		require.Equal(t, int64(1), keeper.supplyKeeper.GetSupply(ctx).GetTotal().AmountOf("linknft").Int64())
	}

	// Send
	{
		require.NoError(t, bk.SendCoins(ctx, addr1, addr2, sdk.NewCoins(sdk.NewCoin("linknft", sdk.NewInt(1)))))
		require.Equal(t, int64(0), bk.GetCoins(ctx, addr1).AmountOf("linknft").Int64())
		require.Equal(t, int64(1), bk.GetCoins(ctx, addr2).AmountOf("linknft").Int64())
		require.Equal(t, int64(1), keeper.supplyKeeper.GetSupply(ctx).GetTotal().AmountOf("linknft").Int64())
	}

	// Burn from account 1
	{
		require.Error(t, keeper.BurnTokens(ctx, sdk.NewCoins(sdk.NewCoin("linknft", sdk.NewInt(1))), addr1))
		require.Equal(t, int64(0), bk.GetCoins(ctx, addr1).AmountOf("linknft").Int64())
		require.Equal(t, int64(1), bk.GetCoins(ctx, addr2).AmountOf("linknft").Int64())
		require.Equal(t, int64(1), keeper.supplyKeeper.GetSupply(ctx).GetTotal().AmountOf("linknft").Int64())
	}

	// Burn from account 2 ( the owner)
	{
		require.NoError(t, keeper.BurnTokens(ctx, sdk.NewCoins(sdk.NewCoin("linknft", sdk.NewInt(1))), addr2))
		require.Equal(t, int64(0), bk.GetCoins(ctx, addr1).AmountOf("linknft").Int64())
		require.Equal(t, int64(0), bk.GetCoins(ctx, addr2).AmountOf("linknft").Int64())
		require.Equal(t, int64(0), keeper.supplyKeeper.GetSupply(ctx).GetTotal().AmountOf("linknft").Int64())
	}
}

func TestCollectionAndPermission(t *testing.T) {
	input := SetupTestInput(t)
	_, ctx, keeper, ak, _ := input.Cdc, input.Ctx, input.Keeper, input.Ak, input.Bk

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
	issuePerm := types.NewIssuePermission("resource01")
	{
		require.NoError(t, keeper.OccupySymbol(ctx, "resource01", addr1))
		require.True(t, keeper.HasPermission(ctx, addr1, issuePerm))
		require.Error(t, keeper.OccupySymbol(ctx, "resource01", addr1))
		collection, err := keeper.GetCollection(ctx, "resource01")
		require.NoError(t, err)
		require.Equal(t, "resource01", collection.Symbol)
	}
	{
		require.NoError(t, keeper.GrantPermission(ctx, addr1, addr2, issuePerm))
		require.True(t, keeper.HasPermission(ctx, addr1, issuePerm))
		require.True(t, keeper.HasPermission(ctx, addr2, issuePerm))
	}

	issuePerm2 := types.NewIssuePermission("resource02")
	{
		require.NoError(t, keeper.OccupySymbol(ctx, "resource02", addr1))
		require.True(t, keeper.HasPermission(ctx, addr1, issuePerm2))
		require.Error(t, keeper.OccupySymbol(ctx, "resource02", addr1))
		collection, err := keeper.GetCollection(ctx, "resource02")
		require.NoError(t, err)
		require.Equal(t, "resource02", collection.Symbol)
	}
	{
		collections := keeper.GetAllCollections(ctx)
		require.Equal(t, 2, len(collections))
		require.Equal(t, "resource01", collections[0].Symbol)
		require.Equal(t, "resource02", collections[1].Symbol)
	}
}

func TestGetPrefixedTokens(t *testing.T) {
	input := SetupTestInput(t)
	_, ctx, keeper, ak := input.Cdc, input.Ctx, input.Keeper, input.Ak

	addr1 := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	{
		acc := ak.NewAccountWithAddress(ctx, addr1)
		ak.SetAccount(ctx, acc)
	}
	require.NotNil(t, ak.GetAccount(ctx, addr1))
	require.Equal(t, uint64(0), ak.GetAccount(ctx, addr1).GetAccountNumber())

	{
		require.NoError(t, keeper.IssueFT(ctx, types.NewFT("name", "link1", sdk.NewInt(8), true), sdk.NewInt(100), addr1))
		require.NoError(t, keeper.IssueFT(ctx, types.NewFT("name", "link2", sdk.NewInt(8), true), sdk.NewInt(200), addr1))
		require.NoError(t, keeper.IssueFT(ctx, types.NewFT("name", "link3", sdk.NewInt(8), true), sdk.NewInt(300), addr1))
		require.NoError(t, keeper.IssueFT(ctx, types.NewFT("name", "cony1", sdk.NewInt(8), true), sdk.NewInt(400), addr1))
		require.NoError(t, keeper.IssueFT(ctx, types.NewFT("name", "cony2", sdk.NewInt(8), true), sdk.NewInt(400), addr1))
		require.NoError(t, keeper.IssueFT(ctx, types.NewFT("name", "line1", sdk.NewInt(8), true), sdk.NewInt(400), addr1))
		require.NoError(t, keeper.IssueFT(ctx, types.NewFT("name", "line2", sdk.NewInt(8), true), sdk.NewInt(400), addr1))
	}
	{
		tokens := keeper.GetAllTokens(ctx)
		require.Equal(t, 7, len(tokens))
	}
	{
		tokens := keeper.GetPrefixedTokens(ctx, "link")
		require.Equal(t, 3, len(tokens))
	}
	{
		tokens := keeper.GetPrefixedTokens(ctx, "cony")
		require.Equal(t, 2, len(tokens))
	}
	{
		tokens := keeper.GetPrefixedTokens(ctx, "li")
		require.Equal(t, 5, len(tokens))
	}
}
