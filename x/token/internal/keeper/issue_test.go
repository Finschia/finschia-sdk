package keeper

import (
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/token/internal/types"
	"github.com/stretchr/testify/require"
)

func TestKeeper_IssueToken(t *testing.T) {
	ctx := cacheKeeper()
	t.Log("Issue Token")
	expected := types.NewToken(defaultContractID, defaultName, defaultSymbol, defaultImageURI, sdk.NewInt(defaultDecimals), true)
	{
		require.NoError(t, keeper.IssueToken(ctx, expected, sdk.NewInt(defaultAmount), addr1, addr1))
	}
	t.Log("Get Token")
	{
		actual, err := keeper.GetToken(ctx, defaultContractID)
		require.NoError(t, err)
		verifyTokenFunc(t, expected, actual)
	}
	t.Log("Permission")
	{
		require.True(t, keeper.HasPermission(ctx, addr1, types.NewModifyPermission(defaultContractID)))
		require.True(t, keeper.HasPermission(ctx, addr1, types.NewMintPermission(defaultContractID)))
		require.True(t, keeper.HasPermission(ctx, addr1, types.NewBurnPermission(defaultContractID)))
	}
	t.Log("Permission only addr1 has the permissions")
	{
		require.False(t, keeper.HasPermission(ctx, addr2, types.NewModifyPermission(defaultContractID)))
		require.False(t, keeper.HasPermission(ctx, addr2, types.NewMintPermission(defaultContractID)))
		require.False(t, keeper.HasPermission(ctx, addr2, types.NewBurnPermission(defaultContractID)))
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
}

func TestKeeper_IssueTokenNotMintable(t *testing.T) {
	ctx := cacheKeeper()
	t.Log("Issue a Token Not Mintable")
	expected := types.NewToken(defaultContractID, defaultName, defaultSymbol, defaultImageURI, sdk.NewInt(defaultDecimals), false)
	{
		require.NoError(t, keeper.IssueToken(ctx, expected, sdk.NewInt(defaultAmount), addr1, addr1))
	}
	{
		actual, err := keeper.GetToken(ctx, defaultContractID)
		require.NoError(t, err)
		verifyTokenFunc(t, expected, actual)
	}
	t.Log("Permission only addr1 has no mint/burn permissions")
	{
		require.True(t, keeper.HasPermission(ctx, addr1, types.NewModifyPermission(defaultContractID)))
		require.False(t, keeper.HasPermission(ctx, addr1, types.NewMintPermission(defaultContractID)))
		require.False(t, keeper.HasPermission(ctx, addr1, types.NewBurnPermission(defaultContractID)))
	}
}

func TestKeeper_IssueTokenTooLongTokenURI(t *testing.T) {
	ctx := cacheKeeper()

	length1001String := strings.Repeat("Eng글자日本語はスゲ", 91) // 11 * 91 = 1001

	t.Log("issue a token with too long token uri")
	{
		token := types.NewToken(defaultContractID, defaultName, defaultSymbol, length1001String, sdk.NewInt(defaultDecimals), true)
		require.EqualError(t, keeper.IssueToken(ctx, token, sdk.NewInt(defaultAmount), addr1, addr1), types.ErrInvalidImageURILength(types.DefaultCodespace, length1001String).Error())
	}
}
