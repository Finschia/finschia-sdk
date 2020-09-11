package types

import (
	"testing"
	"unicode/utf8"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/line/link-modules/x/contract"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

const (
	ModifyMsgType = "modify_token"
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
	require.Equal(t, msg.ContractID, msg2.ContractID)
	require.Equal(t, msg.TokenIndex, msg2.TokenIndex)
	require.Equal(t, msg.TokenType, msg2.TokenType)
	require.Equal(t, msg.Changes, msg2.Changes)
	require.Equal(t, msg.Owner, msg2.Owner)
}

func TestMsgModify_ValidateBasic(t *testing.T) {
	t.Log("normal case")
	{
		msg := AMsgModify().Build()
		require.NoError(t, msg.ValidateBasic())
	}
	t.Log("empty contractID found")
	{
		msg := AMsgModify().Contract("").Build()
		require.EqualError(t, msg.ValidateBasic(), sdkerrors.Wrap(contract.ErrInvalidContractID, "ContractID: ").Error())
	}
	t.Log("empty owner")
	{
		msg := AMsgModify().Owner(nil).Build()
		require.EqualError(t, msg.ValidateBasic(), sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "owner address cannot be empty").Error())
	}
	t.Log("invalid contractID found")
	{
		msg := AMsgModify().Contract("0123456789001234567890").Build()
		require.EqualError(t,
			msg.ValidateBasic(),
			sdkerrors.Wrapf(contract.ErrInvalidContractID, "ContractID: %s", msg.ContractID).Error())
	}
	t.Log("img uri too long")
	{
		msg := AMsgModify().Changes(NewChangesWithMap(map[string]string{"base_img_uri": length1001String})).
			Build()

		require.EqualError(t, msg.ValidateBasic(), sdkerrors.Wrapf(ErrInvalidBaseImgURILength, "[%s] should be shorter than [%d] UTF-8 characters, current length: [%d]", length1001String, MaxBaseImgURILength, utf8.RuneCountInString(length1001String)).Error())
	}
	t.Log("name too long")
	{
		msg := AMsgModify().Changes(NewChangesWithMap(map[string]string{"name": length1001String})).Build()

		require.EqualError(t, msg.ValidateBasic(), sdkerrors.Wrapf(ErrInvalidNameLength, "[%s] should be shorter than [%d] UTF-8 characters, current length: [%d]", length1001String, MaxTokenNameLength, utf8.RuneCountInString(length1001String)).Error())
	}
	t.Log("invalid changes field")
	{
		msg := AMsgModify().Changes(NewChangesWithMap(map[string]string{"invalid_field": "val"})).Build()

		require.EqualError(t, msg.ValidateBasic(), sdkerrors.Wrap(ErrInvalidChangesField, "Field: invalid_field").Error())
	}
	t.Log("no token uri field")
	{
		msg := AMsgModify().Changes(NewChangesWithMap(map[string]string{"name": "new_name"})).Build()
		require.NoError(t, msg.ValidateBasic())
	}
	t.Log("Test with changes more than max")
	{
		// Given msg with changes more than max
		changeList := make([]Change, MaxChangeFieldsCount+1)
		msg := AMsgModify().Changes(changeList).Build()

		// When validate basic, Then error is occurred
		require.EqualError(t, msg.ValidateBasic(), sdkerrors.Wrapf(ErrInvalidChangesFieldCount, "You can not change fields more than [%d] at once, current count: [%d]", MaxChangeFieldsCount, len(changeList)).Error())
	}
	t.Log("Test with nft token type")
	{
		msg := AMsgModify().TokenType(defaultTokenType).Build()
		require.EqualError(t, msg.ValidateBasic(), sdkerrors.Wrap(ErrInvalidChangesField, "Field: base_img_uri").Error())

		msg = AMsgModify().TokenType(defaultTokenType).
			Changes(NewChangesWithMap(map[string]string{"name": "new_name"})).
			Build()
		require.NoError(t, msg.ValidateBasic())
	}
	t.Log("Test with nft token type and index")
	{
		msg := AMsgModify().TokenType(defaultTokenType).TokenIndex(defaultTokenIndex).Build()
		require.EqualError(t, msg.ValidateBasic(), sdkerrors.Wrap(ErrInvalidChangesField, "Field: base_img_uri").Error())

		msg = AMsgModify().TokenType(defaultTokenType).TokenIndex(defaultTokenIndex).
			Changes(NewChangesWithMap(map[string]string{"name": "new_name"})).
			Build()
		require.NoError(t, msg.ValidateBasic())
	}
	t.Log("Test with ft token type and index")
	{
		msg := AMsgModify().TokenType(defaultTokenTypeFT).TokenIndex(defaultTokenIndex).Build()
		require.EqualError(t, msg.ValidateBasic(), sdkerrors.Wrap(ErrInvalidChangesField, "Field: base_img_uri").Error())

		msg = AMsgModify().TokenType(defaultTokenTypeFT).TokenIndex(defaultTokenIndex).
			Changes(NewChangesWithMap(map[string]string{"name": "new_name"})).
			Build()
		require.NoError(t, msg.ValidateBasic())
	}
	t.Log("Test with ft token type and not index")
	{
		msg := AMsgModify().TokenType(defaultTokenTypeFT).Build()
		require.EqualError(t, msg.ValidateBasic(), sdkerrors.Wrap(ErrTokenTypeFTWithoutIndex, defaultTokenTypeFT).Error())
	}
	t.Log("Test with invalid token type")
	{
		invalidTokenType := "010101"
		msg := AMsgModify().TokenType(invalidTokenType).Build()
		require.EqualError(t, msg.ValidateBasic(), sdkerrors.Wrap(ErrInvalidTokenType, invalidTokenType).Error())
	}
	t.Log("Test with invalid token index")
	{
		invalidTokenIndex := "010101"
		msg := AMsgModify().TokenIndex(invalidTokenIndex).Build()
		require.EqualError(t, msg.ValidateBasic(), sdkerrors.Wrap(ErrInvalidTokenIndex, invalidTokenIndex).Error())
	}
	t.Log("Test with token index not token type")
	{
		msg := AMsgModify().TokenIndex(defaultTokenIndex).Build()
		require.EqualError(t, msg.ValidateBasic(), ErrTokenIndexWithoutType.Error())
	}
}

func AMsgModify() *MsgModifyBuilder {
	return &MsgModifyBuilder{
		msgModify: NewMsgModify(
			sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()),
			defaultContractID,
			"",
			"",
			NewChangesWithMap(map[string]string{
				"name":         "new_name",
				"base_img_uri": "new_base_img_uri",
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

func (b *MsgModifyBuilder) Contract(contractID string) *MsgModifyBuilder {
	b.msgModify.ContractID = contractID
	return b
}

func (b *MsgModifyBuilder) TokenType(tokenType string) *MsgModifyBuilder {
	b.msgModify.TokenType = tokenType
	return b
}

func (b *MsgModifyBuilder) TokenIndex(tokenIndex string) *MsgModifyBuilder {
	b.msgModify.TokenIndex = tokenIndex
	return b
}

func (b *MsgModifyBuilder) Changes(changes Changes) *MsgModifyBuilder {
	b.msgModify.Changes = changes
	return b
}
