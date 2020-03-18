package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/line/link/x/token/internal/types"
	"github.com/stretchr/testify/require"
)

func TestKeeper_BurnTokens(t *testing.T) {
	ctx := cacheKeeper()
	t.Log("Issue Token")
	{
		token := types.NewToken(defaultContractID, defaultName, defaultSymbol, defaultMeta, defaultImageURI, sdk.NewInt(defaultDecimals), true)
		require.NoError(t, keeper.IssueToken(ctx, token, sdk.NewInt(defaultAmount+defaultAmount), addr1, addr1))
	}
	t.Log("TotalSupply supply")
	{
		supply, err := keeper.GetTotalInt(ctx, defaultContractID, types.QuerySupply)
		require.NoError(t, err)
		require.Equal(t, int64(defaultAmount+defaultAmount), supply.Int64())
	}
	t.Log("Balance of Account")
	{
		supply := keeper.GetBalance(ctx, defaultContractID, addr1)
		require.Equal(t, int64(defaultAmount+defaultAmount), supply.Int64())
	}

	t.Log("Burn Tokens by addr1")
	{
		err := keeper.BurnToken(ctx, defaultContractID, sdk.NewInt(defaultAmount), addr1)
		require.NoError(t, err)
	}
	t.Log("Burn 0 Token by addr1")
	{
		err := keeper.BurnToken(ctx, defaultContractID, sdk.NewInt(0), addr1)
		require.Error(t, err)
	}
	t.Log("TotalSupply supply")
	{
		supply, err := keeper.GetTotalInt(ctx, defaultContractID, types.QuerySupply)
		require.Equal(t, int64(defaultAmount), supply.Int64())
		require.NoError(t, err)
	}
	t.Log("Balance of Account 1")
	{
		supply := keeper.GetBalance(ctx, defaultContractID, addr1)
		require.Equal(t, int64(defaultAmount), supply.Int64())
	}
}

func TestKeeper_BurnTokensWithoutPermissions(t *testing.T) {
	ctx := cacheKeeper()
	t.Log("Issue Token")
	{
		token := types.NewToken(defaultContractID, defaultName, defaultSymbol, defaultMeta, defaultImageURI, sdk.NewInt(defaultDecimals), true)
		require.NoError(t, keeper.IssueToken(ctx, token, sdk.NewInt(defaultAmount), addr1, addr1))
	}

	t.Log("Transfer Enough Token")
	{
		err := keeper.Transfer(ctx, addr1, addr2, defaultContractID, sdk.NewInt(defaultAmount))
		require.NoError(t, err)
	}

	t.Log("Burn Tokens by addr2. Expect Fail")
	{
		err := keeper.BurnToken(ctx, defaultContractID, sdk.NewInt(defaultAmount), addr2)
		require.Error(t, err)
		require.EqualError(t, err, sdkerrors.Wrapf(types.ErrTokenNoPermission, "Account: %s, Permission: %s", addr2.String(), types.NewBurnPermission(defaultContractID).String()).Error())
	}
}
