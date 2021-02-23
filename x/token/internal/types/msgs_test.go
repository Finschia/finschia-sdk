package types

import (
	"strings"
	"testing"
	"unicode/utf8"

	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/contract"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

// nolint:dupl
func TestMsgBasics(t *testing.T) {
	cdc := ModuleCdc
	addr := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	length1001String := strings.Repeat("Eng글자日本語はスゲ", 91) // 11 * 91 = 1001

	{
		msg := NewMsgIssue(addr, addr, "name", "BTC", "{}", "imageuri", sdk.NewInt(1), sdk.NewInt(8), true)
		require.Equal(t, "issue_token", msg.Type())
		require.Equal(t, "token", msg.Route())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		require.Equal(t, addr, msg.GetSigners()[0])
		require.NoError(t, msg.ValidateBasic())

		b := msg.GetSignBytes()

		msg2 := MsgIssue{}

		err := cdc.UnmarshalJSON(b, &msg2)
		require.NoError(t, err)

		require.Equal(t, msg.Name, msg2.Name)
		require.Equal(t, msg.Symbol, msg2.Symbol)
		require.Equal(t, msg.ImageURI, msg2.ImageURI)
		require.Equal(t, msg.Owner, msg2.Owner)
		require.Equal(t, msg.Amount, msg.Amount)
		require.Equal(t, msg.Decimals, msg2.Decimals)
		require.Equal(t, msg.Mintable, msg2.Mintable)
	}
	{
		msg := NewMsgIssue(addr, addr, "name", "BTC", "", length1001String, sdk.NewInt(1), sdk.NewInt(8), true)
		require.EqualError(t, msg.ValidateBasic(), sdkerrors.Wrapf(ErrInvalidImageURILength, "[%s] should be shorter than [%d] UTF-8 characters, current length: [%d]", length1001String, MaxImageURILength, utf8.RuneCountInString(length1001String)).Error())
	}
	{
		msg := NewMsgIssue(addr, addr, "name", "", "", length1001String, sdk.NewInt(1), sdk.NewInt(8), true)
		require.EqualError(t, msg.ValidateBasic(), sdkerrors.Wrap(ErrInvalidTokenSymbol, "Symbol: ").Error())
	}
	{
		msg := NewMsgIssue(addr, addr, "name", "123456789012345678901", "", length1001String, sdk.NewInt(1), sdk.NewInt(8), true)
		require.EqualError(t, msg.ValidateBasic(), sdkerrors.Wrap(ErrInvalidTokenSymbol, "Symbol: 123456789012345678901").Error())
	}
	{
		msg := NewMsgIssue(addr, addr, "name", "BCD_A", "", length1001String, sdk.NewInt(1), sdk.NewInt(8), true)
		require.EqualError(t, msg.ValidateBasic(), sdkerrors.Wrap(ErrInvalidTokenSymbol, "Symbol: BCD_A").Error())
	}
	{
		msg := NewMsgIssue(addr, addr, "name", "12", "", length1001String, sdk.NewInt(1), sdk.NewInt(8), true)
		require.EqualError(t, msg.ValidateBasic(), sdkerrors.Wrap(ErrInvalidTokenSymbol, "Symbol: 12").Error())
	}
	{
		msg := NewMsgIssue(addr, addr, length1001String, "BTC", "", "tokenuri", sdk.NewInt(1), sdk.NewInt(8), true)
		require.EqualError(t, msg.ValidateBasic(), sdkerrors.Wrapf(ErrInvalidNameLength, "[%s] should be shorter than [%d] UTF-8 characters, current length: [%d]", length1001String, MaxTokenNameLength, utf8.RuneCountInString(length1001String)).Error())
	}
	{
		msg := NewMsgIssue(addr, addr, "", "BTC", "", "tokenuri", sdk.NewInt(1), sdk.NewInt(8), true)
		require.EqualError(t, msg.ValidateBasic(), ErrInvalidTokenName.Error())
	}
	{
		msg := NewMsgIssue(addr, addr, "name", "BTC", "", "tokenuri", sdk.NewInt(1), sdk.NewInt(19), true)
		require.EqualError(t, msg.ValidateBasic(), sdkerrors.Wrapf(ErrInvalidTokenDecimals, "Decimals: %s", sdk.NewInt(19).String()).Error())
	}
	{
		msg := NewMsgMint(addr, contract.SampleContractID, addr, sdk.NewInt(1))
		require.Equal(t, "mint", msg.Type())
		require.Equal(t, "token", msg.Route())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		require.Equal(t, addr, msg.GetSigners()[0])
		require.NoError(t, msg.ValidateBasic())

		b := msg.GetSignBytes()

		msg2 := MsgMint{}

		err := cdc.UnmarshalJSON(b, &msg2)
		require.NoError(t, err)

		require.Equal(t, msg.To, msg2.To)
		require.Equal(t, msg.Amount, msg2.Amount)
	}
	{
		msg := NewMsgBurn(addr, contract.SampleContractID, sdk.NewInt(1))
		require.Equal(t, "burn", msg.Type())
		require.Equal(t, "token", msg.Route())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		require.Equal(t, addr, msg.GetSigners()[0])
		require.NoError(t, msg.ValidateBasic())

		b := msg.GetSignBytes()

		msg2 := MsgBurn{}

		err := cdc.UnmarshalJSON(b, &msg2)
		require.NoError(t, err)

		require.Equal(t, msg.From, msg2.From)
		require.Equal(t, msg.Amount, msg2.Amount)
	}
	{
		msg := NewMsgBurnFrom(addr1, contract.SampleContractID, addr2, sdk.NewInt(1))
		require.Equal(t, "burn_from", msg.Type())
		require.Equal(t, "token", msg.Route())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		require.Equal(t, addr1, msg.GetSigners()[0])
		require.NoError(t, msg.ValidateBasic())

		b := msg.GetSignBytes()

		msg2 := MsgBurnFrom{}

		err := cdc.UnmarshalJSON(b, &msg2)
		require.NoError(t, err)

		require.Equal(t, msg.Proxy, msg2.Proxy)
		require.Equal(t, msg.From, msg2.From)
		require.Equal(t, msg.Amount, msg2.Amount)
	}

	{
		msg := NewMsgBurnFrom(addr1, contract.SampleContractID, addr1, sdk.NewInt(1))
		require.EqualError(t, msg.ValidateBasic(), sdkerrors.Wrapf(ErrApproverProxySame, "Approver: %s", addr1.String()).Error())

		msg = NewMsgBurnFrom(nil, contract.SampleContractID, addr1, sdk.NewInt(1))
		require.EqualError(t, msg.ValidateBasic(), sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "Proxy cannot be empty").Error())

		msg = NewMsgBurnFrom(addr1, contract.SampleContractID, nil, sdk.NewInt(1))
		require.EqualError(t, msg.ValidateBasic(), sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "From cannot be empty").Error())

		msg = NewMsgBurnFrom(addr1, contract.SampleContractID, addr2, sdk.NewInt(-1))
		require.EqualError(t, msg.ValidateBasic(), sdkerrors.Wrap(ErrInvalidAmount, "-1").Error())
	}
	{
		addr2 := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
		msg := NewMsgGrantPermission(addr, contract.SampleContractID, addr2, NewMintPermission())
		require.Equal(t, "grant_perm", msg.Type())
		require.Equal(t, "token", msg.Route())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		require.Equal(t, addr, msg.GetSigners()[0])
		require.NoError(t, msg.ValidateBasic())

		b := msg.GetSignBytes()

		msg2 := MsgGrantPermission{}

		err := cdc.UnmarshalJSON(b, &msg2)
		require.NoError(t, err)

		require.Equal(t, msg.From, msg2.From)
		require.Equal(t, msg.To, msg2.To)
		require.Equal(t, msg.Permission, msg2.Permission)
	}

	{
		msg := NewMsgRevokePermission(addr, contract.SampleContractID, NewMintPermission())
		require.Equal(t, "revoke_perm", msg.Type())
		require.Equal(t, "token", msg.Route())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		require.Equal(t, addr, msg.GetSigners()[0])
		require.NoError(t, msg.ValidateBasic())

		b := msg.GetSignBytes()

		msg2 := MsgRevokePermission{}

		err := cdc.UnmarshalJSON(b, &msg2)
		require.NoError(t, err)

		require.Equal(t, msg.From, msg2.From)
		require.Equal(t, msg.Permission, msg2.Permission)
	}

	{
		msg := NewMsgTransfer(addr, addr, contract.SampleContractID, sdk.NewInt(4))
		require.Equal(t, "transfer_ft", msg.Type())
		require.Equal(t, "token", msg.Route())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		require.Equal(t, addr, msg.GetSigners()[0])
		require.NoError(t, msg.ValidateBasic())

		b := msg.GetSignBytes()

		msg2 := MsgTransfer{}

		err := cdc.UnmarshalJSON(b, &msg2)
		require.NoError(t, err)

		require.Equal(t, msg.From, msg2.From)
		require.Equal(t, msg.ContractID, msg2.ContractID)
		require.Equal(t, msg.To, msg2.To)
		require.Equal(t, msg.Amount, msg2.Amount)
	}

	{
		msg := NewMsgTransfer(nil, addr, contract.SampleContractID, sdk.NewInt(4))
		require.EqualError(t, msg.ValidateBasic(), sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "from cannot be empty").Error())

		msg = NewMsgTransfer(addr, nil, contract.SampleContractID, sdk.NewInt(4))
		require.EqualError(t, msg.ValidateBasic(), sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "to cannot be empty").Error())

		msg = NewMsgTransfer(addr, addr, "m", sdk.NewInt(4))
		require.EqualError(t, msg.ValidateBasic(), sdkerrors.Wrap(contract.ErrInvalidContractID, "ContractID: m").Error())

		msg = NewMsgTransfer(addr, addr, contract.SampleContractID, sdk.NewInt(-1))
		require.EqualError(t, msg.ValidateBasic(), sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "send amount must be positive").Error())
	}

	{
		msg := NewMsgTransferFrom(addr1, contract.SampleContractID, addr2, addr2, sdk.NewInt(defaultAmount))
		require.Equal(t, "transfer_from", msg.Type())
		require.Equal(t, "token", msg.Route())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		require.Equal(t, addr1, msg.GetSigners()[0])
		require.NoError(t, msg.ValidateBasic())

		b := msg.GetSignBytes()

		msg2 := MsgTransferFrom{}

		err := cdc.UnmarshalJSON(b, &msg2)
		require.NoError(t, err)

		require.Equal(t, msg.Proxy, msg2.Proxy)
		require.Equal(t, msg.From, msg2.From)
		require.Equal(t, msg.To, msg2.To)
		require.Equal(t, msg.ContractID, msg2.ContractID)
		require.Equal(t, msg.Amount, msg2.Amount)
	}

	{
		msg := NewMsgTransferFrom(nil, contract.SampleContractID, addr2, addr2, sdk.NewInt(defaultAmount))
		require.EqualError(t, msg.ValidateBasic(), sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "Proxy cannot be empty").Error())

		msg = NewMsgTransferFrom(addr1, contract.SampleContractID, nil, addr2, sdk.NewInt(defaultAmount))
		require.EqualError(t, msg.ValidateBasic(), sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "From cannot be empty").Error())

		msg = NewMsgTransferFrom(addr1, contract.SampleContractID, addr2, nil, sdk.NewInt(defaultAmount))
		require.EqualError(t, msg.ValidateBasic(), sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "To cannot be empty").Error())
	}
}
