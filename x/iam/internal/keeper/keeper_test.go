package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/iam/internal/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"testing"
)

type permission struct {
	Resource string
	Action   string
}

func newPermission(res, act string) permission {
	return permission{Resource: res, Action: act}
}

func (p permission) Equal(res, act string) bool {
	if p.Resource == res && p.Action == act {
		return true
	}
	return false
}

func (p permission) GetResource() string {
	return p.Resource
}

func (p permission) GetAction() string {
	return p.Action
}

func TestKeeper(t *testing.T) {

	testInput := setupTestInput(t)
	_, ctx, keeper := testInput.cdc, testInput.ctx, testInput.keeper

	addr := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())

	perm := newPermission("resource", "mint")
	accPerm := types.NewAccountPermission(addr)
	accPerm.AddPermission(perm)

	keeper.SetAccountPermission(ctx, accPerm)
	accPerm2 := keeper.GetAccountPermission(ctx, addr)

	require.Equal(t, accPerm.GetAddress(), accPerm2.GetAddress())
	require.True(t, accPerm.HasPermission(perm))
	require.Equal(t, accPerm.HasPermission(perm), accPerm2.HasPermission(perm))
}

func TestPermission(t *testing.T) {
	pms := types.NewPermissions(
		newPermission("1", "action"),
		newPermission("2", "action"),
	)
	require.Equal(t, 2, len(pms))
	require.Equal(t, "1", pms[0].GetResource())

	pms.AddPermission(newPermission("3", "action"))
	pms.AddPermission(newPermission("4", "action"))

	require.Equal(t, 4, len(pms))
	require.Equal(t, "1", pms[0].GetResource())
	require.Equal(t, "4", pms[3].GetResource())

	require.True(t, pms.HasPermission(newPermission("3", "action")))
	pms.RemovePermission(newPermission("2", "action"))

	require.Equal(t, 3, len(pms))

	require.False(t, pms.HasPermission(newPermission("2", "action")))
}

func TestAccountPermissionInheritance(t *testing.T) {
	testInput := setupTestInput(t)
	_, ctx, keeper := testInput.cdc, testInput.ctx, testInput.keeper
	addr := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	childAddr := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	perm1 := newPermission("resource1", "mint")
	perm2 := newPermission("resource2", "burn")

	{
		accPerm := keeper.GetAccountPermission(ctx, addr)
		accPerm.AddPermission(perm1)
		keeper.SetAccountPermission(ctx, accPerm)
		accPermChild := keeper.GetAccountPermission(ctx, childAddr)
		accPermChild = accPermChild.InheritAccountPermission(accPerm)
		accPermChild.AddPermission(perm2)
		keeper.SetAccountPermission(ctx, accPermChild)
	}

	{
		accPermChild := keeper.GetAccountPermission(ctx, childAddr)
		require.True(t, accPermChild.HasPermission(perm1))
		require.True(t, accPermChild.HasPermission(perm2))
	}
	{
		accPerm := keeper.GetAccountPermission(ctx, addr)
		accPerm.RemovePermission(perm1)
		keeper.SetAccountPermission(ctx, accPerm)

		accPermChild := keeper.GetAccountPermission(ctx, childAddr)

		require.False(t, accPermChild.HasPermission(perm1))
		require.True(t, accPermChild.HasPermission(perm2))
	}
}

func TestAccountPermissionInheritanceGenerations(t *testing.T) {
	testInput := setupTestInput(t)
	_, ctx, keeper := testInput.cdc, testInput.ctx, testInput.keeper
	addr1 := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	addr2 := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	addr3 := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())

	perm1 := newPermission("resource1", "mint")
	perm2 := newPermission("resource2", "burn")
	perm3 := newPermission("resource3", "burn")

	{
		accPerm1 := keeper.GetAccountPermission(ctx, addr1)
		accPerm2 := keeper.GetAccountPermission(ctx, addr2)
		accPerm3 := keeper.GetAccountPermission(ctx, addr3)
		accPerm1.AddPermission(perm1)
		accPerm2.AddPermission(perm2)
		accPerm3.AddPermission(perm3)
		keeper.SetAccountPermission(ctx, accPerm1)
		keeper.SetAccountPermission(ctx, accPerm2)
		keeper.SetAccountPermission(ctx, accPerm3)
	}

	{
		accPerm1 := keeper.GetAccountPermission(ctx, addr1)
		accPerm2 := keeper.GetAccountPermission(ctx, addr2)
		accPerm3 := keeper.GetAccountPermission(ctx, addr3)

		accPerm2 = accPerm2.InheritAccountPermission(accPerm1)
		accPerm3 = accPerm3.InheritAccountPermission(accPerm2)

		keeper.SetAccountPermission(ctx, accPerm1)
		keeper.SetAccountPermission(ctx, accPerm2)
		keeper.SetAccountPermission(ctx, accPerm3)

	}

	{
		accPerm1 := keeper.GetAccountPermission(ctx, addr1)
		accPerm2 := keeper.GetAccountPermission(ctx, addr2)
		accPerm3 := keeper.GetAccountPermission(ctx, addr3)

		require.True(t, accPerm1.HasPermission(perm1))

		require.True(t, accPerm2.HasPermission(perm1))
		require.True(t, accPerm2.HasPermission(perm2))

		require.True(t, accPerm3.HasPermission(perm1))
		require.True(t, accPerm3.HasPermission(perm2))
		require.True(t, accPerm3.HasPermission(perm3))

	}

	{
		accPerm1 := keeper.GetAccountPermission(ctx, addr1)
		accPerm2 := keeper.GetAccountPermission(ctx, addr2)
		accPerm3 := keeper.GetAccountPermission(ctx, addr3)

		accPerm1.RemovePermission(perm1)

		require.False(t, accPerm1.HasPermission(perm1))

		//Even perm1 is removed from acc1, acc2 already loaded from store it is not changed
		require.True(t, accPerm2.HasPermission(perm1))
		require.True(t, accPerm2.HasPermission(perm2))

		require.True(t, accPerm3.HasPermission(perm1))
		require.True(t, accPerm3.HasPermission(perm2))
		require.True(t, accPerm3.HasPermission(perm3))

		keeper.SetAccountPermission(ctx, accPerm1)
	}
	{
		accPerm1 := keeper.GetAccountPermission(ctx, addr1)
		accPerm2 := keeper.GetAccountPermission(ctx, addr2)
		accPerm3 := keeper.GetAccountPermission(ctx, addr3)

		require.False(t, accPerm1.HasPermission(perm1))

		require.False(t, accPerm2.HasPermission(perm1))
		require.True(t, accPerm2.HasPermission(perm2))

		require.False(t, accPerm3.HasPermission(perm1))
		require.True(t, accPerm3.HasPermission(perm2))
		require.True(t, accPerm3.HasPermission(perm3))
	}

}
