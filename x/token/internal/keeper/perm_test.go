package keeper

import (
	"testing"

	"github.com/line/lbm-sdk/x/token/internal/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/line/lbm-sdk/types"
)

func preparePermissions(ctx sdk.Context, t *testing.T) types.Permissions {
	expected := types.Permissions{
		types.NewMintPermission(),
		types.NewBurnPermission(),
	}
	t.Log("Prepare Permissions")
	{
		for _, perm := range expected {
			keeper.AddPermission(ctx, addr1, perm)
		}
	}
	return expected
}

func TestKeeper_GetPermissions(t *testing.T) {
	ctx = cacheKeeper()
	expected := preparePermissions(ctx, t)

	t.Log("Compare Permissions")
	{
		actual := keeper.GetPermissions(ctx, addr1)
		for index := range actual {
			require.Equal(t, expected[index].String(), actual[index].String())
		}
	}
}

func TestKeeper_AddPermission(t *testing.T) {
	ctx = cacheKeeper()
	keeper.AddPermission(ctx, addr1, types.NewMintPermission())
	keeper.AddPermission(ctx, addr1, types.NewBurnPermission())
	keeper.AddPermission(ctx, addr1, types.NewModifyPermission())
	require.Equal(t, 3, len(keeper.GetPermissions(ctx, addr1)))
}

func TestKeeper_GrantPermission(t *testing.T) {
	ctx = cacheKeeper()
	expected := preparePermissions(ctx, t)
	t.Log("Grant Permissions addr1 -> addr2")
	{
		for _, perm := range expected {
			err := keeper.GrantPermission(ctx, addr1, addr2, perm)
			require.NoError(t, err)
		}
	}
	t.Log("Grant Permission. addr1 has not the permission")
	{
		err := keeper.GrantPermission(ctx, addr1, addr2, types.NewModifyPermission())
		require.Error(t, err)
	}
}

func TestKeeper_RevokePermission(t *testing.T) {
	ctx = cacheKeeper()
	expected := preparePermissions(ctx, t)
	t.Log("Revoke Permissions addr1")
	{
		for _, perm := range expected {
			err := keeper.RevokePermission(ctx, addr1, perm)
			require.NoError(t, err)
		}
	}
	t.Log("Revoke Permission. addr1 has not the permission")
	{
		err := keeper.RevokePermission(ctx, addr1, types.NewModifyPermission())
		require.Error(t, err)
	}
}

func TestKeeper_HasPermission(t *testing.T) {
	ctx = cacheKeeper()
	expected := preparePermissions(ctx, t)
	t.Log("Has Permissions addr1")
	{
		for _, perm := range expected {
			require.True(t, keeper.HasPermission(ctx, addr1, perm))
		}
	}
	t.Log("Revoke Permission. addr1 has not the permission")
	{
		require.False(t, keeper.HasPermission(ctx, addr1, types.NewModifyPermission()))
	}
}
