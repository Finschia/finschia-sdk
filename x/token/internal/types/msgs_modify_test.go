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
	ModifyMsgType     = "modify_token"
	DefaultContractID = "abcd1234"
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
		msg := AMsgModify().ContractID("").Build()
		require.EqualError(t, msg.ValidateBasic(), sdkerrors.Wrap(contract.ErrInvalidContractID, "ContractID: ").Error())
	}
	t.Log("empty owner")
	{
		msg := AMsgModify().Owner(nil).Build()
		require.EqualError(t, msg.ValidateBasic(), sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "owner address cannot be empty").Error())
	}
	t.Log("invalid contractID found")
	{
		invalidContractID := "invalid2198721987"
		msg := AMsgModify().ContractID(invalidContractID).Build()
		require.EqualError(t, msg.ValidateBasic(), sdkerrors.Wrapf(contract.ErrInvalidContractID, "ContractID: %s", invalidContractID).Error())
	}
	t.Log("image uri too long")
	{
		msg := AMsgModify().Changes(NewChangesWithMap(map[string]string{"img_uri": length1001String})).Build()

		require.EqualError(t, msg.ValidateBasic(), sdkerrors.Wrapf(ErrInvalidImageURILength, "[%s] should be shorter than [%d] UTF-8 characters, current length: [%d]", length1001String, MaxImageURILength, utf8.RuneCountInString(length1001String)).Error())
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
		// Given changes more than max
		changeList := make([]Change, MaxChangeFieldsCount+1)
		msg := AMsgModify().Changes(changeList).Build()

		require.EqualError(t, msg.ValidateBasic(), sdkerrors.Wrapf(ErrInvalidChangesFieldCount, "You can not change fields more than [%d] at once, current count: [%d]", MaxChangeFieldsCount, len(changeList)).Error())
	}
}

func AMsgModify() *MsgModifyBuilder {
	return &MsgModifyBuilder{
		msgModify: NewMsgModify(
			sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address()),
			DefaultContractID,
			NewChangesWithMap(map[string]string{
				"name":    "new_name",
				"img_uri": "new_img_uri",
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

func (b *MsgModifyBuilder) ContractID(contractID string) *MsgModifyBuilder {
	b.msgModify.ContractID = contractID
	return b
}

func (b *MsgModifyBuilder) Changes(changes Changes) *MsgModifyBuilder {
	b.msgModify.Changes = changes
	return b
}
