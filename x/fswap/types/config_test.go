package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	require.Equal(t, false, config.UpdateAllowed)
	require.Equal(t, 1, config.MaxSwaps)
}
