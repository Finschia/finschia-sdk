package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPermission(t *testing.T) {
	// Given permissions
	mintPerm := NewMintPermission(defaultSymbol)
	burnPerm := NewBurnPermission(defaultSymbol)
	modifyPerm := NewModifyPermission(defaultSymbol)

	require.True(t, mintPerm.Validate())
	require.True(t, burnPerm.Validate())
	require.True(t, modifyPerm.Validate())

	require.True(t, mintPerm.Equal(mintPerm.GetResource(), mintPerm.GetAction()))
	require.False(t, mintPerm.Equal(burnPerm.GetResource(), burnPerm.GetAction()))
	require.False(t, mintPerm.Equal(modifyPerm.GetResource(), modifyPerm.GetAction()))

	// When make resource or action empty
	mintPerm.Resource = ""
	burnPerm.Action = ""

	// Then they are invalid
	require.False(t, mintPerm.Validate())
	require.False(t, burnPerm.Validate())
}

func TestPermissionString(t *testing.T) {
	mintPerm := NewMintPermission(defaultSymbol)
	burnPerm := NewBurnPermission(defaultSymbol)
	modifyPerm := NewModifyPermission(defaultSymbol)

	require.Equal(t, mintPerm.String(), defaultSymbol+"-mint")
	require.Equal(t, burnPerm.String(), defaultSymbol+"-burn")
	require.Equal(t, modifyPerm.String(), defaultSymbol+"-modify")
}

func TestPermissionsString(t *testing.T) {
	perms := Permissions{
		NewMintPermission(defaultSymbol),
		NewBurnPermission(defaultSymbol),
		NewModifyPermission(defaultSymbol),
	}

	require.Equal(t, `types.Permissions{types.Permission{Action:"mint", Resource:"linktkn"}, types.Permission{Action:"burn", Resource:"linktkn"}, types.Permission{Action:"modify", Resource:"linktkn"}}`, perms.String())
}
