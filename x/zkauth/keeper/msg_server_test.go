package keeper_test

import (
	"testing"

	"github.com/Finschia/finschia-sdk/simapp"
	sdk "github.com/Finschia/finschia-sdk/types"
	banktypes "github.com/Finschia/finschia-sdk/x/bank/types"
	"github.com/Finschia/finschia-sdk/x/zkauth/keeper"
	datest "github.com/Finschia/finschia-sdk/x/zkauth/testutil"
	"github.com/Finschia/finschia-sdk/x/zkauth/types"
	"github.com/stretchr/testify/require"
)

func setupMsgServer(t testing.TB) (types.MsgServer, *simapp.SimApp, sdk.Context) {
	testApp := datest.ZkAuthKeeper(t)
	k, ctx := testApp.Keeper, testApp.Ctx
	return keeper.NewMsgServerImpl(k), testApp.Simapp, ctx
}

func TestExecution(t *testing.T) {
	msgServer, app, ctx := setupMsgServer(t)
	newCoins := sdk.NewCoins(sdk.NewInt64Coin("stake", 5))
	addrs := simapp.AddTestAddrs(app, ctx, 2, sdk.NewInt(100))
	fromAddr := addrs[0]
	toAddr := addrs[1]

	bankMsg := banktypes.MsgSend{
		Amount:      newCoins,
		FromAddress: fromAddr.String(),
		ToAddress:   toAddr.String(),
	}

	zkAuthSig := types.ZKAuthSignature{}

	msgs := types.NewMsgExecution([]sdk.Msg{&bankMsg}, zkAuthSig)

	resp, err := msgServer.Execution(sdk.WrapSDKContext(ctx), msgs)
	require.NoError(t, err)
	require.NotNil(t, resp)
}
