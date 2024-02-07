package keeper_test

import (
	"context"
	"testing"

	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/zkauth/keeper"
	datest "github.com/Finschia/finschia-sdk/x/zkauth/testutil"
	"github.com/Finschia/finschia-sdk/x/zkauth/types"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := datest.ZkAuthKeeper(t)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}
