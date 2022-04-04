package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	ocproto "github.com/line/ostracon/proto/ostracon/types"

	"github.com/line/lbm-sdk/simapp"
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/bank/types"
)

var (
	coins1000 = sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(1000)))
	coins500  = sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(500)))
	fromAddr  = sdk.BytesToAccAddress([]byte("_____from _____"))
	toAddr    = sdk.BytesToAccAddress([]byte("_______to________"))
)

func TestSendAuthorization(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, ocproto.Header{})
	authorization := types.NewSendAuthorization(coins1000)

	t.Log("verify authorization returns valid method name")
	require.Equal(t, authorization.MsgTypeURL(), "/lbm.bank.v1.MsgSend")
	require.NoError(t, authorization.ValidateBasic())
	send := types.NewMsgSend(fromAddr, toAddr, coins1000)

	require.NoError(t, authorization.ValidateBasic())

	t.Log("verify updated authorization returns nil")
	resp, err := authorization.Accept(ctx, send)
	require.NoError(t, err)
	require.True(t, resp.Delete)
	require.Nil(t, resp.Updated)

	authorization = types.NewSendAuthorization(coins1000)
	require.Equal(t, authorization.MsgTypeURL(), "/lbm.bank.v1.MsgSend")
	require.NoError(t, authorization.ValidateBasic())
	send = types.NewMsgSend(fromAddr, toAddr, coins500)
	require.NoError(t, authorization.ValidateBasic())
	resp, err = authorization.Accept(ctx, send)

	t.Log("verify updated authorization returns remaining spent limit")
	require.NoError(t, err)
	require.False(t, resp.Delete)
	require.NotNil(t, resp.Updated)
	sendAuth := types.NewSendAuthorization(coins500)
	require.Equal(t, sendAuth.String(), resp.Updated.String())

	t.Log("expect updated authorization nil after spending remaining amount")
	resp, err = resp.Updated.Accept(ctx, send)
	require.NoError(t, err)
	require.True(t, resp.Delete)
	require.Nil(t, resp.Updated)
}
