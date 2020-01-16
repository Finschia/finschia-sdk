package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/token/internal/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"testing"
)

const (
	defaultName     = "description"
	defaultSymbol   = "linktkn"
	defaultTokenURI = "token-uri"
	defaultDecimals = 6
	defaultAmount   = 1000
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

	t.Log("issue a token")
	{
		token := types.NewFT(defaultName, defaultSymbol, sdk.NewInt(defaultDecimals), true)
		require.NoError(t, keeper.IssueFT(ctx, token, sdk.NewInt(defaultAmount), addr1))
		tokens := keeper.GetAllTokens(ctx)
		require.Equal(t, defaultName, tokens[0].GetName())
		require.Equal(t, defaultSymbol, tokens[0].GetSymbol())
		require.Equal(t, int64(defaultDecimals), tokens[0].(types.FT).GetDecimals().Int64())
		require.Equal(t, true, tokens[0].(types.FT).GetMintable())
		require.Equal(t, int64(defaultAmount), keeper.accountKeeper.GetAccount(ctx, addr1).GetCoins().AmountOf(defaultSymbol).Int64())
	}
	t.Log("issue tokens and get tokens")
	{
		require.NoError(t, keeper.IssueFT(ctx, types.NewFT(defaultName, defaultSymbol+"1", sdk.NewInt(defaultDecimals), true), sdk.NewInt(100), addr1))
		require.NoError(t, keeper.IssueFT(ctx, types.NewFT(defaultName, defaultSymbol+"2", sdk.NewInt(defaultDecimals), true), sdk.NewInt(200), addr1))
		require.NoError(t, keeper.IssueFT(ctx, types.NewFT(defaultName, defaultSymbol+"3", sdk.NewInt(defaultDecimals), true), sdk.NewInt(300), addr1))
		require.NoError(t, keeper.IssueFT(ctx, types.NewFT(defaultName, defaultSymbol+"4", sdk.NewInt(defaultDecimals), true), sdk.NewInt(400), addr1))
		tokens := keeper.GetAllTokens(ctx)
		{
			require.Equal(t, defaultName, tokens[0].GetName())
			require.Equal(t, defaultSymbol, tokens[0].GetSymbol())
			require.Equal(t, true, tokens[0].(types.FT).GetMintable())
		}
		{
			require.Equal(t, defaultSymbol+"1", tokens[1].GetSymbol())
			require.Equal(t, defaultSymbol+"2", tokens[2].GetSymbol())
			require.Equal(t, defaultSymbol+"3", tokens[3].GetSymbol())
			require.Equal(t, defaultSymbol+"4", tokens[4].GetSymbol())
		}

		{
			require.Equal(t, int64(100), keeper.accountKeeper.GetAccount(ctx, addr1).GetCoins().AmountOf(defaultSymbol+"1").Int64())
			require.Equal(t, int64(200), keeper.accountKeeper.GetAccount(ctx, addr1).GetCoins().AmountOf(defaultSymbol+"2").Int64())
			require.Equal(t, int64(300), keeper.accountKeeper.GetAccount(ctx, addr1).GetCoins().AmountOf(defaultSymbol+"3").Int64())
			require.Equal(t, int64(400), keeper.accountKeeper.GetAccount(ctx, addr1).GetCoins().AmountOf(defaultSymbol+"4").Int64())
		}
	}
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

	t.Log("Issue a token")
	{
		setuptoken := types.NewFT(defaultName, defaultSymbol, sdk.NewInt(defaultDecimals), true)
		require.NoError(t, keeper.IssueFT(ctx, setuptoken, sdk.NewInt(900), addr1))

		token, err := keeper.GetToken(ctx, defaultSymbol)
		require.NoError(t, err)
		require.Equal(t, defaultName, token.GetName())
		require.Equal(t, defaultSymbol, token.GetSymbol())
		require.Equal(t, int64(defaultDecimals), token.(types.FT).GetDecimals().Int64())
		require.Equal(t, true, token.(types.FT).GetMintable())
		require.Equal(t, int64(900), keeper.accountKeeper.GetAccount(ctx, addr1).GetCoins().AmountOf(defaultSymbol).Int64())

		require.NoError(t, keeper.MintTokens(ctx, sdk.NewCoins(sdk.NewCoin(defaultSymbol, sdk.NewInt(99))), addr1))

		require.Equal(t, int64(999), keeper.supplyKeeper.GetSupply(ctx).GetTotal().AmountOf(defaultSymbol).Int64())
		require.Equal(t, int64(999), bk.GetCoins(ctx, addr1).AmountOf(defaultSymbol).Int64())
	}
	t.Log("Issue a token again FAIL")
	{
		token := types.NewFT(defaultName, defaultSymbol, sdk.NewInt(defaultDecimals), true)
		require.Error(t, keeper.IssueFT(ctx, token, sdk.NewInt(900), addr1))
	}

	t.Log("Transfer Token")
	{
		require.NoError(t, bk.SendCoins(ctx, addr1, addr2, sdk.NewCoins(sdk.NewCoin(defaultSymbol, sdk.NewInt(100)))))
		require.Equal(t, int64(899), bk.GetCoins(ctx, addr1).AmountOf(defaultSymbol).Int64())
		require.Equal(t, int64(100), bk.GetCoins(ctx, addr2).AmountOf(defaultSymbol).Int64())
		require.Equal(t, int64(999), keeper.supplyKeeper.GetSupply(ctx).GetTotal().AmountOf(defaultSymbol).Int64())
	}

	t.Log("Transfer Token again")
	{
		require.NoError(t, bk.SendCoins(ctx, addr1, addr2, sdk.NewCoins(sdk.NewCoin(defaultSymbol, sdk.NewInt(100)))))
		require.Equal(t, int64(799), bk.GetCoins(ctx, addr1).AmountOf(defaultSymbol).Int64())
		require.Equal(t, int64(200), bk.GetCoins(ctx, addr2).AmountOf(defaultSymbol).Int64())
		require.Equal(t, int64(999), keeper.supplyKeeper.GetSupply(ctx).GetTotal().AmountOf(defaultSymbol).Int64())
	}

	t.Log("Mint Token")
	{
		require.NoError(t, keeper.MintTokens(ctx, sdk.NewCoins(sdk.NewCoin(defaultSymbol, sdk.NewInt(100))), addr1))
		require.Equal(t, int64(899), bk.GetCoins(ctx, addr1).AmountOf(defaultSymbol).Int64())
		require.Equal(t, int64(200), bk.GetCoins(ctx, addr2).AmountOf(defaultSymbol).Int64())
		require.Equal(t, int64(1099), keeper.supplyKeeper.GetSupply(ctx).GetTotal().AmountOf(defaultSymbol).Int64())
	}

	t.Log("Burn Token")
	{
		require.NoError(t, keeper.BurnTokens(ctx, sdk.NewCoins(sdk.NewCoin(defaultSymbol, sdk.NewInt(100))), addr1))
		require.Equal(t, int64(799), bk.GetCoins(ctx, addr1).AmountOf(defaultSymbol).Int64())
		require.Equal(t, int64(200), bk.GetCoins(ctx, addr2).AmountOf(defaultSymbol).Int64())
		require.Equal(t, int64(999), keeper.supplyKeeper.GetSupply(ctx).GetTotal().AmountOf(defaultSymbol).Int64())
	}

	t.Log("Burn Token again amount > has --> fail")
	{
		require.Error(t, keeper.BurnTokens(ctx, sdk.NewCoins(sdk.NewCoin(defaultSymbol, sdk.NewInt(800))), addr1))
		require.Equal(t, int64(799), bk.GetCoins(ctx, addr1).AmountOf(defaultSymbol).Int64())
		require.Equal(t, int64(200), bk.GetCoins(ctx, addr2).AmountOf(defaultSymbol).Int64())
		require.Equal(t, int64(999), keeper.supplyKeeper.GetSupply(ctx).GetTotal().AmountOf(defaultSymbol).Int64())
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

	t.Log("Issue a nft")
	{
		token := types.NewNFT(defaultName, defaultSymbol, defaultTokenURI, addr1)
		require.NoError(t, keeper.IssueNFT(ctx, token, addr1))
	}
	t.Log("Issue a nft again FAIL")
	{
		token := types.NewNFT(defaultName, defaultSymbol, defaultTokenURI, addr1)
		require.Error(t, keeper.IssueNFT(ctx, token, addr1))
	}

	t.Log("Get the token and check")
	{
		token, err := keeper.GetToken(ctx, defaultSymbol)
		require.NoError(t, err)
		require.Equal(t, defaultName, token.GetName())
		require.Equal(t, defaultSymbol, token.GetSymbol())
		require.Equal(t, defaultTokenURI, token.(types.NFT).GetTokenURI())
		require.Equal(t, int64(1), bk.GetCoins(ctx, addr1).AmountOf(defaultSymbol).Int64())
	}
	t.Log("Mint token -> fail. it is nft")
	{
		require.Error(t, keeper.MintTokens(ctx, sdk.NewCoins(sdk.NewCoin(defaultSymbol, sdk.NewInt(99))), addr1))
		require.Equal(t, int64(1), keeper.supplyKeeper.GetSupply(ctx).GetTotal().AmountOf(defaultSymbol).Int64())
		require.Equal(t, int64(1), bk.GetCoins(ctx, addr1).AmountOf(defaultSymbol).Int64())
	}

	t.Log("Send insufficient")
	{
		require.Error(t, bk.SendCoins(ctx, addr1, addr2, sdk.NewCoins(sdk.NewCoin(defaultSymbol, sdk.NewInt(10)))))
		require.Equal(t, int64(1), bk.GetCoins(ctx, addr1).AmountOf(defaultSymbol).Int64())
		require.Equal(t, int64(0), bk.GetCoins(ctx, addr2).AmountOf(defaultSymbol).Int64())
		require.Equal(t, int64(1), keeper.supplyKeeper.GetSupply(ctx).GetTotal().AmountOf(defaultSymbol).Int64())
	}

	t.Log("Send from addr 1 to addr 2")
	{
		require.NoError(t, bk.SendCoins(ctx, addr1, addr2, sdk.NewCoins(sdk.NewCoin(defaultSymbol, sdk.NewInt(1)))))
		require.Equal(t, int64(0), bk.GetCoins(ctx, addr1).AmountOf(defaultSymbol).Int64())
		require.Equal(t, int64(1), bk.GetCoins(ctx, addr2).AmountOf(defaultSymbol).Int64())
		require.Equal(t, int64(1), keeper.supplyKeeper.GetSupply(ctx).GetTotal().AmountOf(defaultSymbol).Int64())
	}

	t.Log("Burn from account 1 -> fail.")
	{
		require.Error(t, keeper.BurnTokens(ctx, sdk.NewCoins(sdk.NewCoin(defaultSymbol, sdk.NewInt(1))), addr1))
		require.Equal(t, int64(0), bk.GetCoins(ctx, addr1).AmountOf(defaultSymbol).Int64())
		require.Equal(t, int64(1), bk.GetCoins(ctx, addr2).AmountOf(defaultSymbol).Int64())
		require.Equal(t, int64(1), keeper.supplyKeeper.GetSupply(ctx).GetTotal().AmountOf(defaultSymbol).Int64())
	}

	t.Log("Burn from account 2 ( the owner)")
	{
		require.NoError(t, keeper.BurnTokens(ctx, sdk.NewCoins(sdk.NewCoin(defaultSymbol, sdk.NewInt(1))), addr2))
		require.Equal(t, int64(0), bk.GetCoins(ctx, addr1).AmountOf(defaultSymbol).Int64())
		require.Equal(t, int64(0), bk.GetCoins(ctx, addr2).AmountOf(defaultSymbol).Int64())
		require.Equal(t, int64(0), keeper.supplyKeeper.GetSupply(ctx).GetTotal().AmountOf(defaultSymbol).Int64())
	}
}

func TestCollectionAndPermission(t *testing.T) {
	input := SetupTestInput(t)
	_, ctx, keeper, ak, _ := input.Cdc, input.Ctx, input.Keeper, input.Ak, input.Bk

	const (
		resource01 = "resource01"
		resource02 = "resource02"
	)

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
	issuePerm := types.NewIssuePermission(resource01)
	{
		require.NoError(t, keeper.OccupySymbol(ctx, resource01, addr1))
		require.True(t, keeper.HasPermission(ctx, addr1, issuePerm))
		require.Error(t, keeper.OccupySymbol(ctx, resource01, addr1))
		collection, err := keeper.GetCollection(ctx, resource01)
		require.NoError(t, err)
		require.Equal(t, resource01, collection.Symbol)
	}
	{
		require.NoError(t, keeper.GrantPermission(ctx, addr1, addr2, issuePerm))
		require.True(t, keeper.HasPermission(ctx, addr1, issuePerm))
		require.True(t, keeper.HasPermission(ctx, addr2, issuePerm))
	}

	issuePerm2 := types.NewIssuePermission(resource02)
	{
		require.NoError(t, keeper.OccupySymbol(ctx, resource02, addr1))
		require.True(t, keeper.HasPermission(ctx, addr1, issuePerm2))
		require.Error(t, keeper.OccupySymbol(ctx, resource02, addr1))
		collection, err := keeper.GetCollection(ctx, resource02)
		require.NoError(t, err)
		require.Equal(t, resource02, collection.Symbol)
	}
	{
		collections := keeper.GetAllCollections(ctx)
		require.Equal(t, 2, len(collections))
		require.Equal(t, resource01, collections[0].Symbol)
		require.Equal(t, resource02, collections[1].Symbol)
	}
}

func TestGetPrefixedTokens(t *testing.T) {
	input := SetupTestInput(t)
	_, ctx, keeper, ak := input.Cdc, input.Ctx, input.Keeper, input.Ak

	const (
		symbolPrefixLink = "link"
		symbolPrefixCony = "cony"
		symbolPrefixLine = "line"
		symbolPrefixLi   = "li"
	)

	addr1 := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	{
		acc := ak.NewAccountWithAddress(ctx, addr1)
		ak.SetAccount(ctx, acc)
	}
	require.NotNil(t, ak.GetAccount(ctx, addr1))
	require.Equal(t, uint64(0), ak.GetAccount(ctx, addr1).GetAccountNumber())

	{
		require.NoError(t, keeper.IssueFT(ctx, types.NewFT(defaultName, symbolPrefixLink+"1", sdk.NewInt(defaultDecimals), true), sdk.NewInt(defaultAmount), addr1))
		require.NoError(t, keeper.IssueFT(ctx, types.NewFT(defaultName, symbolPrefixLink+"2", sdk.NewInt(defaultDecimals), true), sdk.NewInt(defaultAmount), addr1))
		require.NoError(t, keeper.IssueFT(ctx, types.NewFT(defaultName, symbolPrefixLink+"3", sdk.NewInt(defaultDecimals), true), sdk.NewInt(defaultAmount), addr1))
		require.NoError(t, keeper.IssueFT(ctx, types.NewFT(defaultName, symbolPrefixCony+"1", sdk.NewInt(defaultDecimals), true), sdk.NewInt(defaultAmount), addr1))
		require.NoError(t, keeper.IssueFT(ctx, types.NewFT(defaultName, symbolPrefixCony+"2", sdk.NewInt(defaultDecimals), true), sdk.NewInt(defaultAmount), addr1))
		require.NoError(t, keeper.IssueFT(ctx, types.NewFT(defaultName, symbolPrefixLine+"1", sdk.NewInt(defaultDecimals), true), sdk.NewInt(defaultAmount), addr1))
		require.NoError(t, keeper.IssueFT(ctx, types.NewFT(defaultName, symbolPrefixLine+"2", sdk.NewInt(defaultDecimals), true), sdk.NewInt(defaultAmount), addr1))
	}
	{
		tokens := keeper.GetAllTokens(ctx)
		require.Equal(t, 7, len(tokens))
	}
	{
		tokens := keeper.GetPrefixedTokens(ctx, symbolPrefixLink)
		require.Equal(t, 3, len(tokens))
	}
	{
		tokens := keeper.GetPrefixedTokens(ctx, symbolPrefixCony)
		require.Equal(t, 2, len(tokens))
	}
	{
		tokens := keeper.GetPrefixedTokens(ctx, symbolPrefixLine)
		require.Equal(t, 2, len(tokens))
	}
	{
		tokens := keeper.GetPrefixedTokens(ctx, symbolPrefixLi)
		require.Equal(t, 5, len(tokens))
	}
}
