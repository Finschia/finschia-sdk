package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPermission(t *testing.T) {
	issuePerm := NewIssuePermission()
	mintPerm := NewMintPermission()
	burnPerm := NewBurnPermission()
	modifyPerm := NewModifyPermission()

	require.True(t, issuePerm.Validate())
	require.True(t, mintPerm.Validate())
	require.True(t, burnPerm.Validate())
	require.True(t, modifyPerm.Validate())

	require.True(t, mintPerm.Equal(mintPerm))
	require.False(t, mintPerm.Equal(burnPerm))
	require.False(t, mintPerm.Equal(modifyPerm))
	require.False(t, mintPerm.Equal(issuePerm))
}
