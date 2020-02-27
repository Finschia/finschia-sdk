package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	linktype "github.com/line/link/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

const (
	ModifyMsgType = "modify_token"
	DefaultSymbol = "symbol"
)

func TestNewMsgModify(t *testing.T) {
	msg := AMsgModify().Build()

	require.Equal(t, ModifyMsgType, msg.Type())
	require.Equal(t, ModuleName, msg.Route())
	require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
	require.Equal(t, msg.Owner, msg.GetSigners()[0])
}

func TestMarshalMsgModify(t *testing.T) {
	// Given
	msg := AMsgModify().Build()

	// When marshal and unmarshal it
	msg2 := MsgModify{}
	err := ModuleCdc.UnmarshalJSON(msg.GetSignBytes(), &msg2)
	require.NoError(t, err)

	// Then they are equal
	require.Equal(t, msg.Symbol, msg2.Symbol)
	require.Equal(t, msg.Changes, msg2.Changes)
	require.Equal(t, msg.Owner, msg2.Owner)
}

func TestMsgModify_ValidateBasic(t *testing.T) {
	t.Log("normal case")
	{
		msg := AMsgModify().Build()
		require.NoError(t, msg.ValidateBasic())
	}
	t.Log("empty symbol found")
	{
		msg := AMsgModify().Symbol("").Build()
		require.EqualError(t, msg.ValidateBasic(), ErrInvalidTokenSymbol(DefaultCodespace, "").Error())
	}
	t.Log("empty owner")
	{
		msg := AMsgModify().Owner(nil).Build()
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidAddress("owner address cannot be empty").Error())
	}
	t.Log("invalid symbol found")
	{
		invalidSymbol := "invalidsymbol2198721987"
		msg := AMsgModify().Symbol(invalidSymbol).Build()
		require.EqualError(t, msg.ValidateBasic(), ErrInvalidTokenSymbol(DefaultCodespace, invalidSymbol).Error())
	}
	t.Log("tokenURI too long")
	{
		msg := AMsgModify().Changes(linktype.NewChangesWithMap(map[string]string{"token_uri": length1001String})).Build()

		require.EqualError(t, msg.ValidateBasic(), ErrInvalidTokenURILength(DefaultCodespace, length1001String).Error())
	}
	t.Log("name too long")
	{
		msg := AMsgModify().Changes(linktype.NewChangesWithMap(map[string]string{"name": length1001String})).Build()

		require.EqualError(t, msg.ValidateBasic(), ErrInvalidNameLength(DefaultCodespace, length1001String).Error())
	}
	t.Log("invalid changes field")
	{
		msg := AMsgModify().Changes(linktype.NewChangesWithMap(map[string]string{"invalid_field": "val"})).Build()

		require.EqualError(t, msg.ValidateBasic(), ErrInvalidChangesField(DefaultCodespace, "invalid_field").Error())
	}
	t.Log("no token uri field")
	{
		msg := AMsgModify().Changes(linktype.NewChangesWithMap(map[string]string{"name": "new_name"})).Build()
		require.NoError(t, msg.ValidateBasic())
	}
	t.Log("Test with changes more than max")
	{
		// Given changes more than max
		changeList := make([]linktype.Change, MaxChangeFieldsCount+1)
		msg := AMsgModify().Changes(changeList).Build()

		require.EqualError(t, msg.ValidateBasic(), ErrInvalidChangesFieldCount(DefaultCodespace, len(changeList)).Error())
	}
}

func AMsgModify() *MsgModifyBuilder {
	return &MsgModifyBuilder{
		msgModify: NewMsgModify(
			sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()),
			DefaultSymbol,
			linktype.NewChangesWithMap(map[string]string{
				"name":      "new_name",
				"token_uri": "new_torken_uri",
			}),
		),
	}
}

type MsgModifyBuilder struct {
	msgModify MsgModify
}

func (b *MsgModifyBuilder) Build() MsgModify {
	return b.msgModify
}

func (b *MsgModifyBuilder) Owner(owner sdk.AccAddress) *MsgModifyBuilder {
	b.msgModify.Owner = owner
	return b
}

func (b *MsgModifyBuilder) Symbol(symbol string) *MsgModifyBuilder {
	b.msgModify.Symbol = symbol
	return b
}

func (b *MsgModifyBuilder) Changes(changes linktype.Changes) *MsgModifyBuilder {
	b.msgModify.Changes = changes
	return b
}
