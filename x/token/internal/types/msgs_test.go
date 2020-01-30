package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"testing"
)

func TestMsgBasics(t *testing.T) {
	cdc := ModuleCdc
	addr := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	addr2 := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())

	addrSuffix := types.AccAddrSuffix(addr)
	{
		msg := NewMsgIssue("name", "symb"+addrSuffix, "tokenuri", addr, sdk.NewInt(1), sdk.NewInt(8), true)
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
		require.Equal(t, msg.TokenURI, msg2.TokenURI)
		require.Equal(t, msg.Owner, msg2.Owner)
		require.Equal(t, msg.Amount, msg.Amount)
		require.Equal(t, msg.Decimals, msg2.Decimals)
		require.Equal(t, msg.Mintable, msg2.Mintable)
	}
	{
		msg := NewMsgIssueNFT("name", "symb"+addrSuffix, "tokenuri", addr)
		require.Equal(t, "issue_nft", msg.Type())
		require.Equal(t, "token", msg.Route())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		require.Equal(t, addr, msg.GetSigners()[0])
		require.NoError(t, msg.ValidateBasic())

		b := msg.GetSignBytes()

		msg2 := MsgIssueNFT{}

		err := cdc.UnmarshalJSON(b, &msg2)
		require.NoError(t, err)

		require.Equal(t, msg.Name, msg2.Name)
		require.Equal(t, msg.Symbol, msg2.Symbol)
		require.Equal(t, msg.TokenURI, msg2.TokenURI)
		require.Equal(t, msg.Owner, msg2.Owner)
	}
	{
		msg := NewMsgIssueCollection("name", "symb"+addrSuffix, "tokenuri", addr, sdk.NewInt(1), sdk.NewInt(8), true, "0okenid0")
		require.Equal(t, "issue_token_collection", msg.Type())
		require.Equal(t, "token", msg.Route())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		require.Equal(t, addr, msg.GetSigners()[0])
		require.NoError(t, msg.ValidateBasic())

		b := msg.GetSignBytes()

		msg2 := MsgIssueCollection{}

		err := cdc.UnmarshalJSON(b, &msg2)
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
		require.Equal(t, "issue_nft_collection", msg.Type())
		require.Equal(t, "token", msg.Route())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		require.Equal(t, addr, msg.GetSigners()[0])
		require.NoError(t, msg.ValidateBasic())

		b := msg.GetSignBytes()

		msg2 := MsgIssueNFTCollection{}

		err := cdc.UnmarshalJSON(b, &msg2)
		require.NoError(t, err)

		require.Equal(t, msg.Name, msg2.Name)
		require.Equal(t, msg.Symbol, msg2.Symbol)
		require.Equal(t, msg.TokenURI, msg2.TokenURI)
		require.Equal(t, msg.Owner, msg2.Owner)
		require.Equal(t, msg.TokenID, msg2.TokenID)
	}
	{
		msg := NewMsgMint(addr, addr, sdk.NewCoins(sdk.NewCoin("link", sdk.NewInt(1))))
		require.Equal(t, "mint_token", msg.Type())
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
		msg := NewMsgBurn(addr, sdk.NewCoins(sdk.NewCoin("link", sdk.NewInt(1))))
		require.Equal(t, "burn_token", msg.Type())
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
		require.Equal(t, "grant_permission", msg.Type())
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
		require.Equal(t, "revoke_permission", msg.Type())
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
		msg := NewMsgTransferFT(addr, addr, "mytoken", sdk.NewInt(4))
		require.Equal(t, "transfer-ft", msg.Type())
		require.Equal(t, "token", msg.Route())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		require.Equal(t, addr, msg.GetSigners()[0])
		require.NoError(t, msg.ValidateBasic())

		b := msg.GetSignBytes()

		msg2 := MsgTransferFT{}

		err := cdc.UnmarshalJSON(b, &msg2)
		require.NoError(t, err)

		require.Equal(t, msg.FromAddress, msg2.FromAddress)
		require.Equal(t, msg.ToAddress, msg2.ToAddress)
		require.Equal(t, msg.Symbol, msg2.Symbol)
		require.Equal(t, msg.Amount, msg2.Amount)
	}

	{
		msg := NewMsgTransferCFT(addr, addr2, "symbol", "00000001", sdk.NewInt(1))
		require.Equal(t, "transfer-cft", msg.Type())
		require.Equal(t, "token", msg.Route())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		require.Equal(t, addr, msg.GetSigners()[0])
		require.NoError(t, msg.ValidateBasic())

		b := msg.GetSignBytes()

		msg2 := MsgTransferCFT{}

		err := cdc.UnmarshalJSON(b, &msg2)
		require.NoError(t, err)

		require.Equal(t, msg.FromAddress, msg2.FromAddress)
		require.Equal(t, msg.ToAddress, msg2.ToAddress)
		require.Equal(t, msg.Symbol, msg2.Symbol)
		require.Equal(t, msg.TokenID, msg2.TokenID)
		require.Equal(t, msg.Amount, msg2.Amount)
	}

	{
		msg := NewMsgTransferNFT(addr, addr2, "symbol")
		require.Equal(t, "transfer-nft", msg.Type())
		require.Equal(t, "token", msg.Route())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		require.Equal(t, addr, msg.GetSigners()[0])
		require.NoError(t, msg.ValidateBasic())

		b := msg.GetSignBytes()

		msg2 := MsgTransferNFT{}

		err := cdc.UnmarshalJSON(b, &msg2)
		require.NoError(t, err)

		require.Equal(t, msg.FromAddress, msg2.FromAddress)
		require.Equal(t, msg.ToAddress, msg2.ToAddress)
		require.Equal(t, msg.Symbol, msg2.Symbol)
	}

	{
		msg := NewMsgTransferCNFT(addr, addr2, "symbol", "00000001")
		require.Equal(t, "transfer-cnft", msg.Type())
		require.Equal(t, "token", msg.Route())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		require.Equal(t, addr, msg.GetSigners()[0])
		require.NoError(t, msg.ValidateBasic())

		b := msg.GetSignBytes()

		msg2 := MsgTransferCNFT{}

		err := cdc.UnmarshalJSON(b, &msg2)
		require.NoError(t, err)

		require.Equal(t, msg.FromAddress, msg2.FromAddress)
		require.Equal(t, msg.ToAddress, msg2.ToAddress)
		require.Equal(t, msg.Symbol, msg2.Symbol)
		require.Equal(t, msg.TokenID, msg2.TokenID)
	}

	{
		msg := NewMsgAttach(addr, "symbol", "item0001", "item0002")
		require.Equal(t, "attach", msg.Type())
		require.Equal(t, "token", msg.Route())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		require.Equal(t, addr, msg.GetSigners()[0])
		require.NoError(t, msg.ValidateBasic())

		b := msg.GetSignBytes()

		msg2 := MsgAttach{}

		err := cdc.UnmarshalJSON(b, &msg2)
		require.NoError(t, err)

		require.Equal(t, msg.FromAddress, msg2.FromAddress)
		require.Equal(t, msg.ToTokenID, msg2.ToTokenID)
		require.Equal(t, msg.Symbol, msg2.Symbol)
		require.Equal(t, msg.TokenID, msg2.TokenID)
	}

	{
		msg := NewMsgDetach(addr, addr2, "symbol", "item0001")
		require.Equal(t, "detach", msg.Type())
		require.Equal(t, "token", msg.Route())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		require.Equal(t, addr, msg.GetSigners()[0])
		require.NoError(t, msg.ValidateBasic())

		b := msg.GetSignBytes()

		msg2 := MsgDetach{}

		err := cdc.UnmarshalJSON(b, &msg2)
		require.NoError(t, err)

		require.Equal(t, msg.FromAddress, msg2.FromAddress)
		require.Equal(t, msg.ToAddress, msg2.ToAddress)
		require.Equal(t, msg.Symbol, msg2.Symbol)
		require.Equal(t, msg.TokenID, msg2.TokenID)
	}
}

func TestMsgModifyTokenURI_ValidateBasicMsgBasics(t *testing.T) {
	cdc := ModuleCdc
	addr := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	const ModifyActionName = "modify_token"
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
