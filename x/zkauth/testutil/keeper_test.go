package testutil

import (
	"testing"

	"github.com/Finschia/finschia-sdk/x/zkauth/types"
	"github.com/stretchr/testify/require"
)

func TestKeeper(t *testing.T) {
	k, ctx := ZkAuthKeeper(t)

	params, err := k.GetParams(ctx)
	require.Equal(t, types.DefaultParams(), params)
	require.NoError(t, err)
}
