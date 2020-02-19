package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

//nolint:dupl
func TestMsgBasics(t *testing.T) {
	cdc := ModuleCdc
	proxy := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	onBehalfOf := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	receiver := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())

	// MsgProxyApproveCoins
	{
		msg := NewMsgProxyApproveCoins(proxy, onBehalfOf, "link", sdk.NewInt(1))
		require.Equal(t, MsgTypeProxyApproveCoins, msg.Type())
		require.Equal(t, RouterKey, msg.Route())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		require.Equal(t, onBehalfOf, msg.GetSigners()[0])
		require.NoError(t, msg.ValidateBasic())

		b := msg.GetSignBytes()

		msg2 := MsgProxyApproveCoins{}

		err := cdc.UnmarshalJSON(b, &msg2)
		require.NoError(t, err)

		require.Equal(t, msg.Proxy, msg2.Proxy)
		require.Equal(t, msg.OnBehalfOf, msg2.OnBehalfOf)
		require.Equal(t, msg.Denom, msg2.Denom)
		require.Equal(t, msg.Amount, msg.Amount)
	}

	// MsgProxyDisapproveCoins
	{
		msg := NewMsgProxyDisapproveCoins(proxy, onBehalfOf, "link", sdk.NewInt(1))
		require.Equal(t, MsgTypeProxyDisapproveCoins, msg.Type())
		require.Equal(t, RouterKey, msg.Route())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		require.Equal(t, onBehalfOf, msg.GetSigners()[0])
		require.NoError(t, msg.ValidateBasic())

		b := msg.GetSignBytes()

		msg2 := MsgProxyDisapproveCoins{}

		err := cdc.UnmarshalJSON(b, &msg2)
		require.NoError(t, err)

		require.Equal(t, msg.Proxy, msg2.Proxy)
		require.Equal(t, msg.OnBehalfOf, msg2.OnBehalfOf)
		require.Equal(t, msg.Denom, msg2.Denom)
		require.Equal(t, msg.Amount, msg.Amount)
	}

	// MsgProxySendCoinsFrom
	{
		msg := NewMsgProxySendCoinsFrom(proxy, onBehalfOf, receiver, "link", sdk.NewInt(1))
		require.Equal(t, MsgTypeProxySendCoinsFrom, msg.Type())
		require.Equal(t, RouterKey, msg.Route())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		require.Equal(t, proxy, msg.GetSigners()[0])
		require.NoError(t, msg.ValidateBasic())

		b := msg.GetSignBytes()

		msg2 := MsgProxySendCoinsFrom{}

		err := cdc.UnmarshalJSON(b, &msg2)
		require.NoError(t, err)

		require.Equal(t, msg.Proxy, msg2.Proxy)
		require.Equal(t, msg.OnBehalfOf, msg2.OnBehalfOf)
		require.Equal(t, msg.Denom, msg2.Denom)
		require.Equal(t, msg.Amount, msg.Amount)
		require.Equal(t, msg.ToAddress, msg.ToAddress)
	}
}
