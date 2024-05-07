package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Finschia/finschia-sdk/x/fbridge/testutil"
	"github.com/Finschia/finschia-sdk/x/fbridge/types"
)

func TestAssignRole(t *testing.T) {
	key, memKey, ctx, encCfg, authKeeper, bankKeeper, addrs := testutil.PrepareFbridgeTest(t, 3)
	auth := types.DefaultAuthority()
	k := NewKeeper(encCfg.Codec, key, memKey, authKeeper, bankKeeper, "stake", auth.String())
	err := k.InitGenesis(ctx, types.DefaultGenesisState())
	require.NoError(t, err)

	// 1. Bridge authority assigns an address to a guardian role
	p, err := k.RegisterRoleProposal(ctx, addrs[0], addrs[1], types.RoleGuardian)
	require.Error(t, err, "role proposal must not be passed without authority")
	require.Equal(t, types.RoleProposal{}, p)
	p, err = k.RegisterRoleProposal(ctx, auth, addrs[0], types.RoleGuardian)
	require.NoError(t, err)
	require.EqualValues(t, 1, p.Id)
	err = k.updateRole(ctx, types.RoleGuardian, addrs[0])
	require.NoError(t, err)
	require.Equal(t, types.RoleGuardian, k.GetRole(ctx, addrs[0]))
	require.Equal(t, types.RoleMetadata{Guardian: 1, Operator: 0, Judge: 0}, k.GetRoleMetadata(ctx))

	// 2. Guardian assigns an address to a guardian role
	_, err = k.RegisterRoleProposal(ctx, auth, addrs[1], types.RoleGuardian)
	require.Error(t, err, "role proposal must be passed with guardian role after guardian group is formed")
	p, err = k.RegisterRoleProposal(ctx, addrs[0], addrs[1], types.RoleGuardian)
	require.NoError(t, err, "role proposal must be passed with guardian role")
	require.EqualValues(t, 2, p.Id)
	err = k.addVote(ctx, p.Id, addrs[0], types.OptionYes)
	require.NoError(t, err)
	opt, err := k.GetVote(ctx, p.Id, addrs[0])
	require.NoError(t, err)
	require.Equal(t, types.OptionYes, opt)
	err = k.updateRole(ctx, types.RoleGuardian, addrs[1])
	require.NoError(t, err)
	require.Equal(t, types.RoleMetadata{Guardian: 2, Operator: 0, Judge: 0}, k.GetRoleMetadata(ctx))
	sws := k.GetBridgeSwitches(ctx)
	require.Len(t, sws, 2)
	for _, sw := range sws {
		require.Equal(t, types.StatusActive, sw.Status)
	}

	// 3. Guardian assigns an address to an operator role
	err = k.updateRole(ctx, types.RoleOperator, addrs[1])
	require.NoError(t, err)
	require.Equal(t, types.RoleMetadata{Guardian: 1, Operator: 1, Judge: 0}, k.GetRoleMetadata(ctx))

	// 4. Guardian assigns an address to a same role
	err = k.updateRole(ctx, types.RoleOperator, addrs[1])
	require.Error(t, err, "role must not be updated to the same role")
}

func TestBridgeHaltAndResume(t *testing.T) {
	key, memKey, ctx, encCfg, authKeeper, bankKeeper, addrs := testutil.PrepareFbridgeTest(t, 3)
	auth := types.DefaultAuthority()
	k := NewKeeper(encCfg.Codec, key, memKey, authKeeper, bankKeeper, "stake", auth.String())
	err := k.InitGenesis(ctx, types.DefaultGenesisState())
	require.NoError(t, err)
	for _, addr := range addrs {
		err = k.updateRole(ctx, types.RoleGuardian, addr)
		require.NoError(t, err)
	}

	require.Equal(t, types.StatusActive, k.GetBridgeStatus(ctx), "bridge status must be active (3/3)")
	require.Equal(t, types.BridgeStatusMetadata{Active: 3, Inactive: 0}, k.GetBridgeStatusMetadata(ctx))

	err = k.updateBridgeSwitch(ctx, addrs[0], types.StatusInactive)
	require.NoError(t, err)
	require.Equal(t, types.StatusActive, k.GetBridgeStatus(ctx), "bridge status must be active (2/3)")
	require.Equal(t, types.BridgeStatusMetadata{Active: 2, Inactive: 1}, k.GetBridgeStatusMetadata(ctx))

	err = k.updateBridgeSwitch(ctx, addrs[1], types.StatusInactive)
	require.NoError(t, err)
	require.Equal(t, types.StatusInactive, k.GetBridgeStatus(ctx), "bridge status must be inactive (1/3)")
	require.Equal(t, types.BridgeStatusMetadata{Active: 1, Inactive: 2}, k.GetBridgeStatusMetadata(ctx))

	err = k.updateBridgeSwitch(ctx, addrs[0], types.StatusActive)
	require.NoError(t, err)
	require.Equal(t, types.StatusActive, k.GetBridgeStatus(ctx), "bridge status must be active (2/3)")
	require.Equal(t, types.BridgeStatusMetadata{Active: 2, Inactive: 1}, k.GetBridgeStatusMetadata(ctx))
}
