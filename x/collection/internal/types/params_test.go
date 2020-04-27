package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewParams(t *testing.T) {
	params := NewParams(10, 20)
	require.Equal(t, uint64(10), params.MaxComposableDepth)
	require.Equal(t, uint64(20), params.MaxComposableWidth)
}

func TestValidate(t *testing.T) {
	require.NoError(t, NewParams(20, 20).Validate())
	require.Error(t, NewParams(0, 20).Validate())
	require.Error(t, NewParams(20, 0).Validate())
}

func TestDefaultParams(t *testing.T) {
	params := DefaultParams()
	require.Equal(t, DefaultMaxComposableDepth, params.MaxComposableDepth)
	require.Equal(t, DefaultMaxComposableWidth, params.MaxComposableWidth)
}
