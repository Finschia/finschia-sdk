package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPermission(t *testing.T) {
	// Given permissions
	mintPerm := NewMintPermission()
	burnPerm := NewBurnPermission()
	modifyPerm := NewModifyPermission()

	require.True(t, mintPerm.Validate())
	require.True(t, burnPerm.Validate())
	require.True(t, modifyPerm.Validate())

	require.True(t, mintPerm.Equal(mintPerm))
	require.False(t, mintPerm.Equal(burnPerm))
	require.False(t, mintPerm.Equal(modifyPerm))

	// When make resource or action empty
	mintPerm = ""
	burnPerm = ""

	// Then they are invalid
	require.False(t, mintPerm.Validate())
	require.False(t, burnPerm.Validate())
}

func TestPermissionString(t *testing.T) {
	mintPerm := NewMintPermission()
	burnPerm := NewBurnPermission()
	modifyPerm := NewModifyPermission()

	require.Equal(t, mintPerm.String(), "mint")
	require.Equal(t, burnPerm.String(), "burn")
	require.Equal(t, modifyPerm.String(), "modify")
}

func TestPermissionsString(t *testing.T) {
	perms := Permissions{
		NewMintPermission(),
		NewBurnPermission(),
		NewModifyPermission(),
	}

	require.Equal(t, `types.Permissions{"mint", "burn", "modify"}`, perms.String())
}
