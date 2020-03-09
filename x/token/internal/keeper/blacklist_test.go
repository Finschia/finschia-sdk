package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestKeeper_SetBlackList(t *testing.T) {
	ctx := cacheKeeper()
	keeper.SetBlackList(ctx, addr1, "every")
	require.True(t, keeper.IsBlacklisted(ctx, addr1, "every"))
}

func TestKeeper_IsBlacklisted(t *testing.T) {
	ctx := cacheKeeper()
	keeper.SetBlackList(ctx, addr1, "every")
	require.True(t, keeper.IsBlacklisted(ctx, addr1, "every"))
	require.False(t, keeper.IsBlacklisted(ctx, addr2, "every"))
	require.False(t, keeper.IsBlacklisted(ctx, addr1, "every2"))
}
