package keeper

import (
	"testing"

	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/token/internal/types"
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
		supply, err := keeper.GetTotalInt(ctx, types.QuerySupply)
		require.NoError(t, err)
		require.Equal(t, int64(defaultAmount), supply.Int64())
	}
	t.Log("Balance of Account")
	{
		supply := keeper.GetBalance(ctx, addr1)
		require.Equal(t, int64(defaultAmount), supply.Int64())
	}

	t.Log("Mint Tokens addr1 -> addr1")
	{
		err := keeper.MintToken(ctx, sdk.NewInt(defaultAmount), addr1, addr1)
		require.NoError(t, err)
	}
	t.Log("Mint 0 Token")
	{
		err := keeper.MintToken(ctx, sdk.NewInt(0), addr1, addr1)
		require.Error(t, err)
	}
	t.Log("TotalSupply supply")
	{
		supply, err := keeper.GetTotalInt(ctx, types.QuerySupply)
		require.Equal(t, int64(defaultAmount+defaultAmount), supply.Int64())
		require.NoError(t, err)
	}
	t.Log("Balance of Account 1")
	{
		supply := keeper.GetBalance(ctx, addr1)
		require.Equal(t, int64(defaultAmount+defaultAmount), supply.Int64())
	}
	t.Log("Mint Tokens addr1 -> addr2")
	{
		err := keeper.MintToken(ctx, sdk.NewInt(defaultAmount), addr1, addr2)
		require.NoError(t, err)
	}
	t.Log("TotalSupply supply")
	{
		supply, err := keeper.GetTotalInt(ctx, types.QuerySupply)
		require.Equal(t, int64(defaultAmount+defaultAmount+defaultAmount), supply.Int64())
		require.NoError(t, err)
	}
	t.Log("Balance of Account 1")
	{
		supply := keeper.GetBalance(ctx, addr1)
		require.Equal(t, int64(defaultAmount+defaultAmount), supply.Int64())
	}
	t.Log("Balance of Account 2")
	{
		supply := keeper.GetBalance(ctx, addr2)
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
		err := keeper.MintToken(ctx, sdk.NewInt(defaultAmount), addr2, addr2)
		require.Error(t, err)
		require.EqualError(t, err, sdkerrors.Wrapf(types.ErrTokenNoPermission, "Account: %s, Permission: %s", addr2.String(), types.NewMintPermission().String()).Error())
	}
}
