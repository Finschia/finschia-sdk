package keeper_test

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Finschia/finschia-sdk/simapp"
	sdk "github.com/Finschia/finschia-sdk/types"
	banktypes "github.com/Finschia/finschia-sdk/x/bank/types"
	"github.com/Finschia/finschia-sdk/x/zkauth/keeper"
	datest "github.com/Finschia/finschia-sdk/x/zkauth/testutil"
	"github.com/Finschia/finschia-sdk/x/zkauth/types"
)

func setupMsgServer(t testing.TB) (types.MsgServer, *simapp.SimApp, sdk.Context) {
	testApp := datest.ZkAuthKeeper(t)
	k, ctx := testApp.ZKAuthKeeper, testApp.Ctx
	return keeper.NewMsgServerImpl(*k), testApp.Simapp, ctx
}

func TestExecution(t *testing.T) {
	msgServer, app, ctx := setupMsgServer(t)
	newCoins := sdk.NewCoins(sdk.NewInt64Coin("stake", 5))
	addrs := simapp.AddTestAddrs(app, ctx, 2, sdk.NewInt(100))
	fromAddr := sdk.MustAccAddressFromBech32("link1f7j46mgxr3d5tgn3qsnn9ep8j77yr8w4ks0f3lpzz9rhnwja9vcqej33ld")
	toAddr := addrs[0]
	err := simapp.FundAccount(app, ctx, fromAddr, sdk.NewCoins(sdk.NewCoin(app.StakingKeeper.BondDenom(ctx), sdk.NewInt(100))))

	bankMsg := banktypes.MsgSend{
		Amount:      newCoins,
		FromAddress: fromAddr.String(),
		ToAddress:   toAddr.String(),
	}

	testProofB64 := "eyJwaV9hIjpbIjM0MjM1MjQ3MjQzMDgyMDE5Mzc5NTQyMDUwMjA3MTMyNDg0MTg5OTQ5NjEzODk1OTcxODgxNjk1ODg4Njk5Mzg0MDgyMjI1NjE2MzAiLCI4OTQ1NzMxODIxOTg1NDU0MTcwNDA5MjQ2MjA4MDUxMzYwMzcxNzQ0NjUxOTYwNjg2MjQ5Njg1NjExMDI0NzQ4MjQ5MDk2NjUyNDAwIiwiMSJdLCJwaV9iIjpbWyIxNjE3NzgwNjczNzQ1MTg2MDY3NDk3NTAxNjI1MjExODM3NDU0NDU5NjE5NzcxMDAyNjYzMjAxMTg2MTg4OTM4NDkxMDc3MTU5ODUwIiwiMTAzODU2MDMyMjQzMDc2MDQxMjM5NjI2NzAyMzA3NDcxNDM2NDc2ODAxOTAwNTQzNzU3OTE5NjAxNjYwMTY2NzMyNTcwNjE5ODExNjgiXSxbIjYyMTE1NjA5NTUzMTE1NzgxNDcwMzEzMjAzMTE1MzAyODM5MzgwMDY5MTIyMjAwMTUwODM1MTIxODg2MDI5MzU0MjA4Mjk5MjY1NCIsIjEzODUwMDYzMTA5ODg2MTE3MDcwNzgzMzIxNDk1ODI4MzQ2MTg3ODQ3OTM3OTIyMDMyNTI2NzU3MTU2MDg0MjUzMTU4NTc0NjgyMTYyIl0sWyIxIiwiMCJdXSwicGlfYyI6WyIxMDE0MTM1NTQwNTE1MDg5MzQ1MjY3NTUwMTE4Mjk4MTY1MzY5NDI1MzM0MzMyMjAyNzM1OTgwNTU0NzYxODgyODE4NjU3NDU5MTA4IiwiNDkzMjY5MjI2OTcyOTIyNDQ4NjI1MDI0MDAyMjI4NTQ4NzE1MjYzNTkwNjAzMzA1ODYwNjgwMjI2Nzg4Nzg5NjE2MTU2NTM4MzYiLCIxIl19"
	testProof, err := base64.StdEncoding.DecodeString(testProofB64)
	require.NoError(t, err)

	zkAuthSig := types.ZKAuthSignature{
		ZkAuthInputs: &types.ZKAuthInputs{
			ProofPoints:  testProof,
			IssBase64:    "aHR0cHM6Ly9hY2NvdW50cy5nb29nbGUuY29t",
			HeaderBase64: "eyJhbGciOiJSUzI1NiIsImtpZCI6ImFkZjVlNzEwZWRmZWJlY2JlZmE5YTYxNDk1NjU0ZDAzYzBiOGVkZjgiLCJ0eXAiOiJKV1QifQ",
			AddressSeed:  "16929294337693897740349056748881744503581363933815798702166939594477679063350",
		},
		MaxBlockHeight: 104,
	}

	msgs := types.NewMsgExecution([]sdk.Msg{&bankMsg}, zkAuthSig)

	resp, err := msgServer.Execution(sdk.WrapSDKContext(ctx), msgs)
	require.NoError(t, err)
	require.NotNil(t, resp)
}
