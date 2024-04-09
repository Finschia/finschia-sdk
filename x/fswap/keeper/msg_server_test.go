package keeper_test

import (
	"context"
	"testing"

	keepertest "fswap/testutil/keeper"
	"fswap/x/fswap/keeper"
	"fswap/x/fswap/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.FswapKeeper(t)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}
