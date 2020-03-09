package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/token/internal/types"
	"github.com/stretchr/testify/require"
)

func TestKeeper_MintTokens(t *testing.T) {
	ctx := cacheKeeper()
	t.Log("Issue Token")
	{
		token := types.NewToken(defaultContractID, defaultName, defaultSymbol, defaultMeta, defaultImageURI, sdk.NewInt(defaultDecimals), true)
		require.NoError(t, keeper.IssueToken(ctx, token, sdk.NewInt(defaultAmount), addr1, addr1))
	}
	t.Log("TotalSupply supply")
	{
		supply, err := keeper.GetTotalInt(ctx, defaultContractID, types.QuerySupply)
		require.NoError(t, err)
		require.Equal(t, int64(defaultAmount), supply.Int64())
	}
	t.Log("Balance of Account")
	{
		supply := keeper.GetBalance(ctx, defaultContractID, addr1)
		require.Equal(t, int64(defaultAmount), supply.Int64())
	}

	t.Log("Mint Tokens addr1 -> addr1")
	{
		err := keeper.MintToken(ctx, defaultContractID, sdk.NewInt(defaultAmount), addr1, addr1)
		require.NoError(t, err)
	}
	t.Log("TotalSupply supply")
	{
		supply, err := keeper.GetTotalInt(ctx, defaultContractID, types.QuerySupply)
		require.Equal(t, int64(defaultAmount+defaultAmount), supply.Int64())
		require.NoError(t, err)
	}
	t.Log("Balance of Account 1")
	{
		supply := keeper.GetBalance(ctx, defaultContractID, addr1)
		require.Equal(t, int64(defaultAmount+defaultAmount), supply.Int64())
	}
	t.Log("Mint Tokens addr1 -> addr2")
	{
		err := keeper.MintToken(ctx, defaultContractID, sdk.NewInt(defaultAmount), addr1, addr2)
		require.NoError(t, err)
	}
	t.Log("TotalSupply supply")
	{
		supply, err := keeper.GetTotalInt(ctx, defaultContractID, types.QuerySupply)
		require.Equal(t, int64(defaultAmount+defaultAmount+defaultAmount), supply.Int64())
		require.NoError(t, err)
	}
	t.Log("Balance of Account 1")
	{
		supply := keeper.GetBalance(ctx, defaultContractID, addr1)
		require.Equal(t, int64(defaultAmount+defaultAmount), supply.Int64())
	}
	t.Log("Balance of Account 2")
	{
		supply := keeper.GetBalance(ctx, defaultContractID, addr2)
		require.Equal(t, int64(defaultAmount), supply.Int64())
	}
}

func TestKeeper_MintTokensWithoutPermissions(t *testing.T) {
	ctx := cacheKeeper()
	t.Log("Issue Token")
	{
		token := types.NewToken(defaultContractID, defaultName, defaultSymbol, defaultMeta, defaultImageURI, sdk.NewInt(defaultDecimals), true)
		require.NoError(t, keeper.IssueToken(ctx, token, sdk.NewInt(defaultAmount), addr1, addr1))
	}

	t.Log("Mint Tokens by addr2. Expect Fail")
	{
		err := keeper.MintToken(ctx, defaultContractID, sdk.NewInt(defaultAmount), addr2, addr2)
		require.Error(t, err)
		require.EqualError(t, err, types.ErrTokenNoPermission(types.DefaultCodespace, addr2, types.NewMintPermission(defaultContractID)).Error())
	}
}
