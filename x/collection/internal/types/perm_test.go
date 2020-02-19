package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPermission(t *testing.T) {
	issuePerm := NewIssuePermission(defaultSymbol)
	mintPerm := NewMintPermission(defaultSymbol)
	burnPerm := NewBurnPermission(defaultSymbol)
	modifyPerm := NewModifyTokenURIPermission(defaultSymbol)

	require.True(t, issuePerm.Validate())
	require.True(t, mintPerm.Validate())
	require.True(t, burnPerm.Validate())
	require.True(t, modifyPerm.Validate())

	require.True(t, mintPerm.Equal(mintPerm.GetResource(), mintPerm.GetAction()))
	require.False(t, mintPerm.Equal(burnPerm.GetResource(), burnPerm.GetAction()))
	require.False(t, mintPerm.Equal(modifyPerm.GetResource(), modifyPerm.GetAction()))
	require.False(t, mintPerm.Equal(issuePerm.GetResource(), issuePerm.GetAction()))
}
