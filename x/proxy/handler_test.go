package proxy

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	testCommon "github.com/line/link/x/proxy/keeper"
	"github.com/line/link/x/proxy/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

func TestHandler(t *testing.T) {
	input := testCommon.SetupTestInput(t)
	ctx, keeper, ak := input.Ctx, input.Keeper, input.Ak
	ctx = ctx.WithEventManager(sdk.NewEventManager()) // to track emitted events

	h := NewHandler(keeper)

	res := h(ctx, sdk.NewTestMsg())
	require.False(t, res.IsOK())
	require.Equal(t, res.Codespace, types.DefaultCodespace)
	require.Equal(t, res.Code, types.CodeProxyInvalidMsgType)

	var err sdk.Error

	// proxy
	proxy := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	{
		acc := ak.NewAccountWithAddress(ctx, proxy)
		ak.SetAccount(ctx, acc)
	}
	require.NotNil(t, ak.GetAccount(ctx, proxy))

	// on behalf of
	onBehalfOf := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	{
		acc := ak.NewAccountWithAddress(ctx, onBehalfOf)
		ak.SetAccount(ctx, acc)
	}
	require.NotNil(t, ak.GetAccount(ctx, onBehalfOf))

	// receiver
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
		msgPxSendCoinsFrom := types.NewMsgProxySendCoinsFrom(proxy, onBehalfOf, receiver, "link", sdk.OneInt())

		// should fail as it's not approved
		res = h(ctx, msgPxSendCoinsFrom)
		require.False(t, res.IsOK())
		require.Equal(t, res.Codespace, types.DefaultCodespace)
		require.Equal(t, res.Code, types.CodeProxyNotExist)

		// check emitted events (no event should be emitted)
		testCommon.VerifyEventFunc(t, sdk.Events{}, res.Events)
	}

	// refresh the event manager to verify upcoming events only
	ctx = ctx.WithEventManager(sdk.NewEventManager())

	// `onBehalfOf` approves 5 link for `proxy`
	approvedAmount := sdk.NewInt(5)
	{
		msgPxApproveCoins := types.NewMsgProxyApproveCoins(proxy, onBehalfOf, "link", approvedAmount)
		res = h(ctx, msgPxApproveCoins)
		require.True(t, res.IsOK())

		// check emitted events
		e := sdk.Events{
			sdk.NewEvent(
				EventProxyApproveCoins,
				sdk.NewAttribute(AttributeKeyProxyAddress, msgPxApproveCoins.Proxy.String()),
				sdk.NewAttribute(AttributeKeyProxyOnBehalfOfAddress, msgPxApproveCoins.OnBehalfOf.String()),
				sdk.NewAttribute(AttributeKeyProxyDenom, msgPxApproveCoins.Denom),
				sdk.NewAttribute(AttributeKeyProxyAmount, msgPxApproveCoins.Amount.String()),
			),
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeySender, msgPxApproveCoins.OnBehalfOf.String()),
				sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
			),
		}
		testCommon.VerifyEventFunc(t, e, res.Events)
	}

	// refresh the event manager to verify upcoming events only
	ctx = ctx.WithEventManager(sdk.NewEventManager())

	// 'proxy' tries to send 6 link to `receiver` on behalf of `onBehalfOf`
	{
		msgPxSendCoinsFrom := types.NewMsgProxySendCoinsFrom(proxy, onBehalfOf, receiver, "link", sdk.NewInt(6))

		// should fail as it's more than approved
		res = h(ctx, msgPxSendCoinsFrom)
		require.False(t, res.IsOK())
		require.Equal(t, res.Codespace, types.DefaultCodespace)
		require.Equal(t, res.Code, types.CodeProxyNotEnoughApprovedCoins)

		// check emitted events (no event should be emitted)
		testCommon.VerifyEventFunc(t, sdk.Events{}, res.Events)
	}

	// refresh the event manager to verify upcoming events only
	ctx = ctx.WithEventManager(sdk.NewEventManager())

	// `proxy` sends 2 link to `receiver` on behalf of `onBehalfOf`
	sentAmount1 := sdk.NewInt(2)
	{
		msgPxSendCoinsFrom := types.NewMsgProxySendCoinsFrom(proxy, onBehalfOf, receiver, "link", sentAmount1)

		// should succeed
		res = h(ctx, msgPxSendCoinsFrom)
		require.True(t, res.IsOK())

		// check emitted events
		e := sdk.Events{
			sdk.NewEvent(
				"transfer",
				sdk.NewAttribute("recipient", msgPxSendCoinsFrom.ToAddress.String()),
				sdk.NewAttribute(sdk.AttributeKeyAmount, sentAmount1.String()+"link"),
			),
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeySender, msgPxSendCoinsFrom.OnBehalfOf.String()),
			),
			sdk.NewEvent(
				EventProxySendCoinsFrom,
				sdk.NewAttribute(AttributeKeyProxyAddress, msgPxSendCoinsFrom.Proxy.String()),
				sdk.NewAttribute(AttributeKeyProxyOnBehalfOfAddress, msgPxSendCoinsFrom.OnBehalfOf.String()),
				sdk.NewAttribute(AttributeKeyProxyToAddress, msgPxSendCoinsFrom.ToAddress.String()),
				sdk.NewAttribute(AttributeKeyProxyDenom, msgPxSendCoinsFrom.Denom),
				sdk.NewAttribute(AttributeKeyProxyAmount, msgPxSendCoinsFrom.Amount.String()),
			),
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeySender, msgPxSendCoinsFrom.OnBehalfOf.String()),
				sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
			),
		}
		testCommon.VerifyEventFunc(t, e, res.Events)
	}

	// refresh the event manager to verify upcoming events only
	ctx = ctx.WithEventManager(sdk.NewEventManager())

	// check balance of `onBehalfOf` and `receiver`
	{
		onBehalfOfBalance := ak.GetAccount(ctx, onBehalfOf).GetCoins().AmountOf("link")
		require.Equal(t, initialBalance.Sub(sentAmount1), onBehalfOfBalance)
		receiverBalance := ak.GetAccount(ctx, receiver).GetCoins().AmountOf("link")
		require.Equal(t, sentAmount1, receiverBalance)
	}

	// refresh the event manager to verify upcoming events only
	ctx = ctx.WithEventManager(sdk.NewEventManager())

	// `onBehalfOf` tries to disapprove 4 link from `proxy`
	{
		msgPxDisapproveCoins := types.NewMsgProxyDisapproveCoins(proxy, onBehalfOf, "link", sdk.NewInt(4))

		// should fail as only 3 approved coins are left
		res = h(ctx, msgPxDisapproveCoins)
		require.False(t, res.IsOK())
		require.Equal(t, res.Codespace, types.DefaultCodespace)
		require.Equal(t, res.Code, types.CodeProxyNotEnoughApprovedCoins)

		// check emitted events (no event should be emitted)
		testCommon.VerifyEventFunc(t, sdk.Events{}, res.Events)
	}

	// refresh the event manager to verify upcoming events only
	ctx = ctx.WithEventManager(sdk.NewEventManager())

	// `onBehalfOf` disapprove 1 link from `proxy`
	{
		msgPxDisapproveCoins := types.NewMsgProxyDisapproveCoins(proxy, onBehalfOf, "link", sdk.OneInt())

		// should succeed
		res = h(ctx, msgPxDisapproveCoins)
		require.True(t, res.IsOK())

		// check emitted events
		e := sdk.Events{
			sdk.NewEvent(
				EventProxyDisapproveCoins,
				sdk.NewAttribute(AttributeKeyProxyAddress, msgPxDisapproveCoins.Proxy.String()),
				sdk.NewAttribute(AttributeKeyProxyOnBehalfOfAddress, msgPxDisapproveCoins.OnBehalfOf.String()),
				sdk.NewAttribute(AttributeKeyProxyDenom, msgPxDisapproveCoins.Denom),
				sdk.NewAttribute(AttributeKeyProxyAmount, msgPxDisapproveCoins.Amount.String()),
			),
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeySender, msgPxDisapproveCoins.OnBehalfOf.String()),
				sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
			),
		}
		testCommon.VerifyEventFunc(t, e, res.Events)
	}

	// refresh the event manager to verify upcoming events only
	ctx = ctx.WithEventManager(sdk.NewEventManager())

	// `proxy` sends 2 link to `receiver` on behalf of `onBehalfOf`
	sentAmount2 := sdk.NewInt(2)
	{
		msgPxSendCoinsFrom := types.NewMsgProxySendCoinsFrom(proxy, onBehalfOf, receiver, "link", sentAmount2)

		// should succeed
		res = h(ctx, msgPxSendCoinsFrom)
		require.True(t, res.IsOK())

		// check emitted events
		e := sdk.Events{
			sdk.NewEvent(
				"transfer",
				sdk.NewAttribute("recipient", msgPxSendCoinsFrom.ToAddress.String()),
				sdk.NewAttribute(sdk.AttributeKeyAmount, sentAmount2.String()+"link"),
			),
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeySender, msgPxSendCoinsFrom.OnBehalfOf.String()),
			),
			sdk.NewEvent(
				EventProxySendCoinsFrom,
				sdk.NewAttribute(AttributeKeyProxyAddress, msgPxSendCoinsFrom.Proxy.String()),
				sdk.NewAttribute(AttributeKeyProxyOnBehalfOfAddress, msgPxSendCoinsFrom.OnBehalfOf.String()),
				sdk.NewAttribute(AttributeKeyProxyToAddress, msgPxSendCoinsFrom.ToAddress.String()),
				sdk.NewAttribute(AttributeKeyProxyDenom, msgPxSendCoinsFrom.Denom),
				sdk.NewAttribute(AttributeKeyProxyAmount, msgPxSendCoinsFrom.Amount.String()),
			),
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeySender, msgPxSendCoinsFrom.OnBehalfOf.String()),
				sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
			),
		}
		testCommon.VerifyEventFunc(t, e, res.Events)
	}

	// refresh the event manager to verify upcoming events only
	ctx = ctx.WithEventManager(sdk.NewEventManager())

	// check balance of `onBehalfOf` and `receiver`
	{
		onBehalfOfBalance := ak.GetAccount(ctx, onBehalfOf).GetCoins().AmountOf("link")
		require.Equal(t, initialBalance.Sub(sentAmount1).Sub(sentAmount2), onBehalfOfBalance)
		receiverBalance := ak.GetAccount(ctx, receiver).GetCoins().AmountOf("link")
		require.Equal(t, sentAmount1.Add(sentAmount2), receiverBalance)
	}

	// refresh the event manager to verify upcoming events only
	ctx = ctx.WithEventManager(sdk.NewEventManager())

	// 'proxy' tries to send 1 link to `receiver` on behalf of `onBehalfOf`
	{
		msgPxSendCoinsFrom := types.NewMsgProxySendCoinsFrom(proxy, onBehalfOf, receiver, "link", sdk.OneInt())

		// should fail as there is no coin approved (all sent!)
		res = h(ctx, msgPxSendCoinsFrom)
		require.False(t, res.IsOK())
		require.Equal(t, res.Codespace, types.DefaultCodespace)
		require.Equal(t, res.Code, types.CodeProxyNotExist)

		// check emitted events (no event should be emitted)
		testCommon.VerifyEventFunc(t, sdk.Events{}, res.Events)
	}

	// refresh the event manager to verify upcoming events only
	ctx = ctx.WithEventManager(sdk.NewEventManager())

	// 'onBehalfOf' tries to disapprove 1 link from `proxy`
	{
		msgPxDisapproveCoins := types.NewMsgProxyDisapproveCoins(proxy, onBehalfOf, "link", sdk.OneInt())

		// should fail as there is no proxy anymore (all sent!)
		res = h(ctx, msgPxDisapproveCoins)
		require.False(t, res.IsOK())
		require.Equal(t, res.Codespace, types.DefaultCodespace)
		require.Equal(t, res.Code, types.CodeProxyNotExist)

		// check emitted events (no event should be emitted)
		testCommon.VerifyEventFunc(t, sdk.Events{}, res.Events)
	}
}
