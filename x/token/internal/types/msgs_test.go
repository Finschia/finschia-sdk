package types

import (
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/contract"
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
		require.EqualError(t, msg.ValidateBasic(), ErrInvalidImageURILength(DefaultCodespace, length1001String).Error())
	}
	{
		msg := NewMsgIssue(addr, addr, "name", "", "", length1001String, sdk.NewInt(1), sdk.NewInt(8), true)
		require.EqualError(t, msg.ValidateBasic(), ErrInvalidTokenSymbol(DefaultCodespace, "").Error())
	}
	{
		msg := NewMsgIssue(addr, addr, "name", "123456789012345678901", "", length1001String, sdk.NewInt(1), sdk.NewInt(8), true)
		require.EqualError(t, msg.ValidateBasic(), ErrInvalidTokenSymbol(DefaultCodespace, "123456789012345678901").Error())
	}
	{
		msg := NewMsgIssue(addr, addr, "name", "BCD_A", "", length1001String, sdk.NewInt(1), sdk.NewInt(8), true)
		require.EqualError(t, msg.ValidateBasic(), ErrInvalidTokenSymbol(DefaultCodespace, "BCD_A").Error())
	}
	{
		msg := NewMsgIssue(addr, addr, "name", "12", "", length1001String, sdk.NewInt(1), sdk.NewInt(8), true)
		require.EqualError(t, msg.ValidateBasic(), ErrInvalidTokenSymbol(DefaultCodespace, "12").Error())
	}
	{
		msg := NewMsgIssue(addr, addr, length1001String, "BTC", "", "tokenuri", sdk.NewInt(1), sdk.NewInt(8), true)
		require.EqualError(t, msg.ValidateBasic(), ErrInvalidNameLength(DefaultCodespace, length1001String).Error())
	}
	{
		msg := NewMsgIssue(addr, addr, "", "BTC", "", "tokenuri", sdk.NewInt(1), sdk.NewInt(8), true)
		require.EqualError(t, msg.ValidateBasic(), ErrInvalidTokenName(DefaultCodespace, "").Error())
	}
	{
		msg := NewMsgIssue(addr, addr, "name", "BTC", "", "tokenuri", sdk.NewInt(1), sdk.NewInt(19), true)
		require.EqualError(t, msg.ValidateBasic(), ErrInvalidTokenDecimals(DefaultCodespace, sdk.NewInt(19)).Error())
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
		addr2 := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
		msg := NewMsgGrantPermission(addr, addr2, Permission{Action: "issue", Resource: "resource"})
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
		msg := NewMsgRevokePermission(addr, Permission{Action: "issue", Resource: "resource"})
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
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidAddress("from cannot be empty").Error())

		msg = NewMsgTransfer(addr, nil, contract.SampleContractID, sdk.NewInt(4))
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidAddress("to cannot be empty").Error())

		msg = NewMsgTransfer(addr, addr, "m", sdk.NewInt(4))
		require.EqualError(t, msg.ValidateBasic(), contract.ErrInvalidContractID(contract.ContractCodeSpace, "m").Error())

		msg = NewMsgTransfer(addr, addr, contract.SampleContractID, sdk.NewInt(-1))
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInsufficientCoins("send amount must be positive").Error())
	}
}
