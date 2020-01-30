package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/proxy/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"testing"
)

func TestProxy(t *testing.T) {
	input := SetupTestInput(t)
	_, ctx, keeper, ak := input.Cdc, input.Ctx, input.Keeper, input.Ak

	var err sdk.Error

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

	// create receiver
	receiver := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	{
		acc := ak.NewAccountWithAddress(ctx, receiver)
		ak.SetAccount(ctx, acc)
	}
	require.NotNil(t, ak.GetAccount(ctx, receiver))

	// onBehalfOf to have 10 links
	initialBalance := sdk.NewInt(10)
	{
		_, err = input.Bk.AddCoins(ctx, onBehalfOf, sdk.NewCoins(sdk.NewCoin("link", initialBalance)))
		require.NoError(t, err)
	}

	// `proxy` tries to send coins to `receiver` on behalf of `onBehalfOf`
	{
		err = keeper.SendCoinsFrom(ctx, types.NewMsgProxySendCoinsFrom(proxy, onBehalfOf, receiver, "link", sdk.NewInt(1)))

		// should fail as it's not approved
		require.EqualError(t, err, types.ErrProxyNotExist(types.DefaultCodespace, proxy.String(), onBehalfOf.String()).Error())
	}

	// `onBehalfOf` approves 5 link for `proxy`
	approvedAmount := sdk.NewInt(5)
	{
		err = keeper.ApproveCoins(ctx, types.NewMsgProxyApproveCoins(proxy, onBehalfOf, "link", approvedAmount))
		require.NoError(t, err)
	}

	// 'proxy' tries to send 6 link to `receiver` on behalf of `onBehalfOf`
	{
		err = keeper.SendCoinsFrom(ctx, types.NewMsgProxySendCoinsFrom(proxy, onBehalfOf, receiver, "link", sdk.NewInt(6)))

		// should fail as it's more than approved
		require.EqualError(t, err, types.ErrProxyNotEnoughApprovedCoins(types.DefaultCodespace, approvedAmount, sdk.NewInt(6)).Error())
	}

	// `proxy` sends 2 link to `receiver` on behalf of `onBehalfOf`
	sentAmount1 := sdk.NewInt(2)
	{
		err = keeper.SendCoinsFrom(ctx, types.NewMsgProxySendCoinsFrom(proxy, onBehalfOf, receiver, "link", sentAmount1))

		// should succeed
		require.NoError(t, err)
	}

	// check balance of `onBehalfOf` and `receiver`
	{
		onBehalfOfBalance := ak.GetAccount(ctx, onBehalfOf).GetCoins().AmountOf("link")
		require.Equal(t, initialBalance.Sub(sentAmount1), onBehalfOfBalance)
		receiverBalance := ak.GetAccount(ctx, receiver).GetCoins().AmountOf("link")
		require.Equal(t, sentAmount1, receiverBalance)
	}

	// `onBehalfOf` tries to disapprove 4 link from `proxy`
	{
		err = keeper.DisapproveCoins(ctx, types.NewMsgProxyDisapproveCoins(proxy, onBehalfOf, "link", sdk.NewInt(4)))

		// should fail as only 3 approved coins are left
		require.EqualError(t, err, types.ErrProxyNotEnoughApprovedCoins(types.DefaultCodespace, approvedAmount.Sub(sentAmount1), sdk.NewInt(4)).Error())
	}

	// `onBehalfOf` disapprove 1 link from `proxy`
	{
		err = keeper.DisapproveCoins(ctx, types.NewMsgProxyDisapproveCoins(proxy, onBehalfOf, "link", sdk.NewInt(1)))

		// should succeed
		require.NoError(t, err)
	}

	// `proxy` sends 2 link to `receiver` on behalf of `onBehalfOf`
	sentAmount2 := sdk.NewInt(2)
	{
		err = keeper.SendCoinsFrom(ctx, types.NewMsgProxySendCoinsFrom(proxy, onBehalfOf, receiver, "link", sentAmount2))

		// should succeed
		require.NoError(t, err)
	}

	// check balance of `onBehalfOf` and `receiver`
	{
		onBehalfOfBalance := ak.GetAccount(ctx, onBehalfOf).GetCoins().AmountOf("link")
		require.Equal(t, initialBalance.Sub(sentAmount1).Sub(sentAmount2), onBehalfOfBalance)
		receiverBalance := ak.GetAccount(ctx, receiver).GetCoins().AmountOf("link")
		require.Equal(t, sentAmount1.Add(sentAmount2), receiverBalance)
	}

	// 'proxy' tries to send 1 link to `receiver` on behalf of `onBehalfOf`
	{
		err = keeper.SendCoinsFrom(ctx, types.NewMsgProxySendCoinsFrom(proxy, onBehalfOf, receiver, "link", sdk.NewInt(1)))

		// should fail as there is no coin approved (all sent!)
		require.EqualError(t, err, types.ErrProxyNotExist(types.DefaultCodespace, proxy.String(), onBehalfOf.String()).Error())
	}

	// 'onBehalfOf' tries to disapprove 1 link from `proxy`
	{
		err = keeper.DisapproveCoins(ctx, types.NewMsgProxyDisapproveCoins(proxy, onBehalfOf, "link", sdk.NewInt(1)))

		// should fail as there is no proxy anymore (all sent!)
		require.EqualError(t, err, types.ErrProxyNotExist(types.DefaultCodespace, proxy.String(), onBehalfOf.String()).Error())
	}
}
