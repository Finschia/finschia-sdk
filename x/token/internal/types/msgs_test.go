package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"testing"
)

func TestMarshalJSON(t *testing.T) {

	cdc := ModuleCdc
	addr := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	addrSuffix := types.AccAddrSuffix(addr)
	{
		msg := NewMsgIssue("name", "symb"+addrSuffix, "tokenuri", addr, sdk.NewInt(1), sdk.NewInt(8), true)
		b, err := cdc.MarshalJSON(msg)
		require.NoError(t, err)

		err = msg.ValidateBasic()
		require.NoError(t, err)

		msg2 := MsgIssue{}

		err = cdc.UnmarshalJSON(b, &msg2)
		require.NoError(t, err)

		require.Equal(t, msg.Name, msg2.Name)
		require.Equal(t, msg.Symbol, msg2.Symbol)
		require.Equal(t, msg.TokenURI, msg2.TokenURI)
		require.Equal(t, msg.Owner, msg2.Owner)
		require.Equal(t, msg.Amount, msg.Amount)
		require.Equal(t, msg.Decimals, msg2.Decimals)
		require.Equal(t, msg.Mintable, msg2.Mintable)
	}
	{
		msg := NewMsgIssueNFT("name", "symb"+addrSuffix, "tokenuri", addr)
		b, err := cdc.MarshalJSON(msg)
		require.NoError(t, err)

		err = msg.ValidateBasic()
		require.NoError(t, err)

		msg2 := MsgIssueNFT{}

		err = cdc.UnmarshalJSON(b, &msg2)
		require.NoError(t, err)

		require.Equal(t, msg.Name, msg2.Name)
		require.Equal(t, msg.Symbol, msg2.Symbol)
		require.Equal(t, msg.TokenURI, msg2.TokenURI)
		require.Equal(t, msg.Owner, msg2.Owner)
	}
	{
		msg := NewMsgIssueCollection("name", "symb"+addrSuffix, "tokenuri", addr, sdk.NewInt(1), sdk.NewInt(8), true, "tokenid0")
		b, err := cdc.MarshalJSON(msg)
		require.NoError(t, err)

		err = msg.ValidateBasic()
		require.NoError(t, err)

		msg2 := MsgIssueCollection{}

		err = cdc.UnmarshalJSON(b, &msg2)
		require.NoError(t, err)

		require.Equal(t, msg.Name, msg2.Name)
		require.Equal(t, msg.Symbol, msg2.Symbol)
		require.Equal(t, msg.TokenURI, msg2.TokenURI)
		require.Equal(t, msg.Owner, msg2.Owner)
		require.Equal(t, msg.Amount, msg.Amount)
		require.Equal(t, msg.Decimals, msg2.Decimals)
		require.Equal(t, msg.Mintable, msg2.Mintable)
		require.Equal(t, msg.TokenID, msg2.TokenID)
	}
	{
		msg := NewMsgIssueNFTCollection("name", "symb"+addrSuffix, "tokenuri", addr, "tokenid0")
		b, err := cdc.MarshalJSON(msg)
		require.NoError(t, err)

		err = msg.ValidateBasic()
		require.NoError(t, err)

		msg2 := MsgIssueNFTCollection{}

		err = cdc.UnmarshalJSON(b, &msg2)
		require.NoError(t, err)

		require.Equal(t, msg.Name, msg2.Name)
		require.Equal(t, msg.Symbol, msg2.Symbol)
		require.Equal(t, msg.TokenURI, msg2.TokenURI)
		require.Equal(t, msg.Owner, msg2.Owner)
		require.Equal(t, msg.TokenID, msg2.TokenID)
	}

}
