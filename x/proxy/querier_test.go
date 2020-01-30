package proxy

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	testCommon "github.com/line/link/x/proxy/keeper"
	"github.com/line/link/x/proxy/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"testing"
)

func TestProxyQuerierAllowance(t *testing.T) {
	input := testCommon.SetupTestInput(t)
	codec, ctx, keeper, ak := input.Cdc, input.Ctx, input.Keeper, input.Ak
	h := NewHandler(input.Keeper)
	denom := "link"
	amount := sdk.NewInt(1)

	// create proxy
	proxy := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	{
		acc := ak.NewAccountWithAddress(ctx, proxy)
		ak.SetAccount(ctx, acc)
	}
	require.NotNil(t, ak.GetAccount(ctx, proxy))

	// create on behalf of
	onBehalfOf := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	{
		acc := ak.NewAccountWithAddress(ctx, onBehalfOf)
		ak.SetAccount(ctx, acc)
	}
	require.NotNil(t, ak.GetAccount(ctx, onBehalfOf))

	// approve 1 link
	msgProxyApproveCoins := types.NewMsgProxyApproveCoins(proxy, onBehalfOf, denom, amount)
	r := h(ctx, msgProxyApproveCoins)
	require.True(t, r.IsOK())

	// query allowance
	params := types.QueryProxyAllowance{ProxyDenom: types.NewProxyDenom(proxy, onBehalfOf, denom)}
	req := abci.RequestQuery{
		Path: fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryAllowance),
		Data: []byte(params.String()),
	}
	path := []string{types.QueryAllowance}
	querier := NewQuerier(keeper)
	res, err := querier(ctx, path, req)
	require.NoError(t, err)

	// unmarshal the response
	var pxa types.ProxyAllowance
	err2 := codec.UnmarshalJSON(res, &pxa)
	require.NoError(t, err2)

	// verify 1 link allowance
	require.Equal(t, proxy, pxa.Proxy)
	require.Equal(t, onBehalfOf, pxa.OnBehalfOf)
	require.Equal(t, denom, pxa.Denom)
	require.Equal(t, amount, pxa.Amount)

	// onBehalfOf to have 5 links
	initialBalance := sdk.NewInt(5)
	{
		_, err = input.Bk.AddCoins(ctx, onBehalfOf, sdk.NewCoins(sdk.NewCoin(denom, initialBalance)))
		require.NoError(t, err)
	}

	// create receiver
	receiver := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	{
		acc := ak.NewAccountWithAddress(ctx, receiver)
		ak.SetAccount(ctx, acc)
	}
	require.NotNil(t, ak.GetAccount(ctx, receiver))

	// the proxy sends 1 link to the receiver on behalf of ...
	msgProxySendCoinsFrom := types.NewMsgProxySendCoinsFrom(proxy, onBehalfOf, receiver, denom, amount)
	r = h(ctx, msgProxySendCoinsFrom)
	require.True(t, r.IsOK())

	// query allowance
	params = types.QueryProxyAllowance{ProxyDenom: types.NewProxyDenom(proxy, onBehalfOf, denom)}
	req = abci.RequestQuery{
		Path: fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryAllowance),
		Data: []byte(params.String()),
	}
	path = []string{types.QueryAllowance}
	querier = NewQuerier(keeper)
	_, err = querier(ctx, path, req)

	// no proxy expected as all the allowance is used
	require.EqualError(t, types.ErrProxyNotExist(DefaultCodespace, proxy.String(), onBehalfOf.String()), err.Error())
}
