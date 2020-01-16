package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"testing"
)

func TestMsgModifyTokenURI_ValidateBasicMsgBasics(t *testing.T) {
	cdc := ModuleCdc
	addr := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	t.Log("normal case")
	{
		msg := NewMsgModifyTokenURI(addr, "symbol", "tokenURI", "tokenid0")
		require.Equal(t, ModifyActionName, msg.Type())
		require.Equal(t, ModuleName, msg.Route())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		require.Equal(t, addr, msg.GetSigners()[0])
		require.NoError(t, msg.ValidateBasic())

		b := msg.GetSignBytes()
		msg2 := MsgModifyTokenURI{}
		err := cdc.UnmarshalJSON(b, &msg2)
		require.NoError(t, err)

		require.Equal(t, msg.Symbol, msg2.Symbol)
		require.Equal(t, msg.TokenURI, msg2.TokenURI)
		require.Equal(t, msg.Owner, msg2.Owner)
		require.Equal(t, msg.TokenID, msg2.TokenID)
	}
	t.Log("empty symbol found")
	{
		msg := NewMsgModifyTokenURI(addr, "", "tokenURI", "tokenid0")
		require.Error(t, msg.ValidateBasic())
	}
	t.Log("empty owner")
	{
		msg := NewMsgModifyTokenURI(nil, "symbol", "tokenURI", "tokenid0")
		require.Error(t, msg.ValidateBasic())
	}
	t.Log("invalid symbol found")
	{
		msg := NewMsgModifyTokenURI(addr, "invalidsymbol2198721987", "tokenURI", "tokenid0")
		require.Error(t, msg.ValidateBasic())
	}
	t.Log("invalid tokenid found")
	{
		msg := NewMsgModifyTokenURI(addr, "symbol", "tokenURI", "tokenid")
		require.Error(t, msg.ValidateBasic())
	}
}
