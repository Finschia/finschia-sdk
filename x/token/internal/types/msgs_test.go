package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
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
		msg := NewMsgIssue("name", "s", "tokenuri", addr, sdk.NewInt(1), sdk.NewInt(8), true)
		require.EqualError(t, msg.ValidateBasic(), ErrInvalidTokenSymbol(DefaultCodespace, "symbol [s] mismatched to [^[a-z][a-z0-9]{5,7}$]").Error())
	}
	{
		msg := NewMsgIssue("", "symb"+addrSuffix, "tokenuri", addr, sdk.NewInt(1), sdk.NewInt(8), true)
		require.EqualError(t, msg.ValidateBasic(), ErrInvalidTokenName(DefaultCodespace, "").Error())
	}
	{
		msg := NewMsgIssue("name", "symb"+addrSuffix, "tokenuri", addr, sdk.NewInt(1), sdk.NewInt(19), true)
		require.EqualError(t, msg.ValidateBasic(), ErrInvalidTokenDecimals(DefaultCodespace, sdk.NewInt(19)).Error())
	}
	{
		msg := NewMsgIssue("name", "symb"+addrSuffix, "tokenuri", addr, sdk.NewInt(1), sdk.NewInt(0), false)
		require.EqualError(t, msg.ValidateBasic(), ErrInvalidIssueFT(DefaultCodespace).Error())
	}
	{
		msg := NewMsgIssueCFT("name", "symb"+addrSuffix, "tokenuri", addr, sdk.NewInt(1), sdk.NewInt(8), true)
		require.Equal(t, "issue_cft", msg.Type())
		require.Equal(t, "token", msg.Route())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		require.Equal(t, addr, msg.GetSigners()[0])
		require.NoError(t, msg.ValidateBasic())

		b := msg.GetSignBytes()

		msg2 := MsgIssueCFT{}

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
		msg := NewMsgMintCNFT("name", "symb"+addrSuffix, "tokenuri", "toke", addr, addr)
		require.Equal(t, "mint_cnft", msg.Type())
		require.Equal(t, "token", msg.Route())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		require.Equal(t, addr, msg.GetSigners()[0])
		require.NoError(t, msg.ValidateBasic())

		b := msg.GetSignBytes()

		msg2 := MsgMintCNFT{}

		err := cdc.UnmarshalJSON(b, &msg2)
		require.NoError(t, err)

		require.Equal(t, msg.Name, msg2.Name)
		require.Equal(t, msg.Symbol, msg2.Symbol)
		require.Equal(t, msg.TokenURI, msg2.TokenURI)
		require.Equal(t, msg.From, msg2.From)
		require.Equal(t, msg.TokenType, msg2.TokenType)
	}
	{
		msg := NewMsgBurnCNFT("symb"+addrSuffix, "10010001", addr)
		require.Equal(t, "burn_cnft", msg.Type())
		require.Equal(t, "token", msg.Route())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		require.Equal(t, addr, msg.GetSigners()[0])
		require.NoError(t, msg.ValidateBasic())

		b := msg.GetSignBytes()

		msg2 := MsgBurnCNFT{}

		err := cdc.UnmarshalJSON(b, &msg2)
		require.NoError(t, err)

		require.Equal(t, msg.Symbol, msg2.Symbol)
		require.Equal(t, msg.From, msg2.From)
		require.Equal(t, msg.TokenID, msg2.TokenID)
	}
	{
		msg := NewMsgMint(addr, addr, sdk.NewCoins(sdk.NewCoin("link", sdk.NewInt(1))))
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
		msg := NewMsgBurn(addr, sdk.NewCoins(sdk.NewCoin("link", sdk.NewInt(1))))
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
		msg := NewMsgTransferFT(addr, addr, "mytoken", sdk.NewInt(4))
		require.Equal(t, "transfer_ft", msg.Type())
		require.Equal(t, "token", msg.Route())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		require.Equal(t, addr, msg.GetSigners()[0])
		require.NoError(t, msg.ValidateBasic())

		b := msg.GetSignBytes()

		msg2 := MsgTransferFT{}

		err := cdc.UnmarshalJSON(b, &msg2)
		require.NoError(t, err)

		require.Equal(t, msg.From, msg2.From)
		require.Equal(t, msg.To, msg2.To)
		require.Equal(t, msg.Symbol, msg2.Symbol)
		require.Equal(t, msg.Amount, msg2.Amount)
	}

	{
		msg := NewMsgTransferFT(nil, addr, "mytoken", sdk.NewInt(4))
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidAddress("missing sender address").Error())

		msg = NewMsgTransferFT(addr, nil, "mytoken", sdk.NewInt(4))
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidAddress("missing recipient address").Error())

		msg = NewMsgTransferFT(addr, addr, "m", sdk.NewInt(4))
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidCoins("Only user defined token is possible: m").Error())

		msg = NewMsgTransferFT(addr, addr, "mytoken", sdk.NewInt(-1))
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInsufficientCoins("send amount must be positive").Error())
	}

	{
		msg := NewMsgTransferCFT(addr, addr2, "symbol", "00000001", sdk.NewInt(1))
		require.Equal(t, "transfer_cft", msg.Type())
		require.Equal(t, "token", msg.Route())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		require.Equal(t, addr, msg.GetSigners()[0])
		require.NoError(t, msg.ValidateBasic())

		b := msg.GetSignBytes()

		msg2 := MsgTransferCFT{}

		err := cdc.UnmarshalJSON(b, &msg2)
		require.NoError(t, err)

		require.Equal(t, msg.From, msg2.From)
		require.Equal(t, msg.To, msg2.To)
		require.Equal(t, msg.Symbol, msg2.Symbol)
		require.Equal(t, msg.TokenID, msg2.TokenID)
		require.Equal(t, msg.Amount, msg2.Amount)
	}

	{
		msg := NewMsgTransferCFT(nil, addr2, "symbol", "00000001", sdk.NewInt(1))
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidAddress("From cannot be empty").Error())

		msg = NewMsgTransferCFT(addr, nil, "symbol", "00000001", sdk.NewInt(1))
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidAddress("To cannot be empty").Error())

		msg = NewMsgTransferCFT(addr, addr2, "", "00000001", sdk.NewInt(1))
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidCoins("Token symbol is empty").Error())

		msg = NewMsgTransferCFT(addr, addr2, "symbol", "1", sdk.NewInt(1))
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidCoins("Only user defined token is possible: symbol1").Error())

		msg = NewMsgTransferCFT(addr, addr2, "symbol", "00000001", sdk.NewInt(-1))
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInsufficientCoins("send amount must be positive").Error())
	}

	{
		msg := NewMsgTransferCNFT(addr, addr2, "symbol", "00000001")
		require.Equal(t, "transfer_cnft", msg.Type())
		require.Equal(t, "token", msg.Route())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		require.Equal(t, addr, msg.GetSigners()[0])
		require.NoError(t, msg.ValidateBasic())

		b := msg.GetSignBytes()

		msg2 := MsgTransferCNFT{}

		err := cdc.UnmarshalJSON(b, &msg2)
		require.NoError(t, err)

		require.Equal(t, msg.From, msg2.From)
		require.Equal(t, msg.To, msg2.To)
		require.Equal(t, msg.Symbol, msg2.Symbol)
		require.Equal(t, msg.TokenID, msg2.TokenID)
	}

	{
		msg := NewMsgTransferCNFT(nil, addr2, "symbol", "00000001")
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidAddress("From cannot be empty").Error())

		msg = NewMsgTransferCNFT(addr, nil, "symbol", "00000001")
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidAddress("To cannot be empty").Error())

		msg = NewMsgTransferCNFT(addr, addr2, "", "00000001")
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidCoins("Token symbol is empty").Error())

		msg = NewMsgTransferCNFT(addr, addr2, "symbol", "1")
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidCoins("Only user defined token is possible: symbol1").Error())
	}

	{
		msg := NewMsgTransferCFTFrom(addr, addr2, addr2, "symbol", "00000001", sdk.NewInt(1))
		require.Equal(t, "transfer_cft_from", msg.Type())
		require.Equal(t, "token", msg.Route())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		require.Equal(t, addr, msg.GetSigners()[0])
		require.NoError(t, msg.ValidateBasic())

		b := msg.GetSignBytes()

		msg2 := MsgTransferCFTFrom{}

		err := cdc.UnmarshalJSON(b, &msg2)
		require.NoError(t, err)

		require.Equal(t, msg.Proxy, msg2.Proxy)
		require.Equal(t, msg.From, msg2.From)
		require.Equal(t, msg.To, msg2.To)
		require.Equal(t, msg.Symbol, msg2.Symbol)
		require.Equal(t, msg.TokenID, msg2.TokenID)
		require.Equal(t, msg.Amount, msg2.Amount)
	}

	{
		msg := NewMsgTransferCFTFrom(nil, addr2, addr2, "symbol", "00000001", sdk.NewInt(1))
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidAddress("Proxy cannot be empty").Error())

		msg = NewMsgTransferCFTFrom(addr, nil, addr2, "symbol", "00000001", sdk.NewInt(1))
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidAddress("From cannot be empty").Error())

		msg = NewMsgTransferCFTFrom(addr, addr2, nil, "symbol", "00000001", sdk.NewInt(1))
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidAddress("To cannot be empty").Error())

		msg = NewMsgTransferCFTFrom(addr, addr2, addr2, "symbol", "1", sdk.NewInt(1))
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidCoins("Only user defined token is possible: symbol1").Error())

		msg = NewMsgTransferCFTFrom(addr, addr2, addr2, "symbol", "00000001", sdk.NewInt(-1))
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInsufficientCoins("send amount must be positive").Error())
	}

	{
		msg := NewMsgTransferCNFTFrom(addr, addr2, addr2, "symbol", "00000001")
		require.Equal(t, "transfer_cnft_from", msg.Type())
		require.Equal(t, "token", msg.Route())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		require.Equal(t, addr, msg.GetSigners()[0])
		require.NoError(t, msg.ValidateBasic())

		b := msg.GetSignBytes()

		msg2 := MsgTransferCNFTFrom{}

		err := cdc.UnmarshalJSON(b, &msg2)
		require.NoError(t, err)

		require.Equal(t, msg.Proxy, msg2.Proxy)
		require.Equal(t, msg.From, msg2.From)
		require.Equal(t, msg.To, msg2.To)
		require.Equal(t, msg.Symbol, msg2.Symbol)
		require.Equal(t, msg.TokenID, msg2.TokenID)
	}

	{
		msg := NewMsgTransferCNFTFrom(nil, addr2, addr2, "symbol", "00000001")
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidAddress("Proxy cannot be empty").Error())

		msg = NewMsgTransferCNFTFrom(addr, nil, addr2, "symbol", "00000001")
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidAddress("From cannot be empty").Error())

		msg = NewMsgTransferCNFTFrom(addr, addr2, nil, "symbol", "00000001")
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidAddress("To cannot be empty").Error())

		msg = NewMsgTransferCNFTFrom(addr, addr2, addr2, "symbol", "1")
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidCoins("Only user defined token is possible: symbol1").Error())
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

		require.Equal(t, msg.From, msg2.From)
		require.Equal(t, msg.ToTokenID, msg2.ToTokenID)
		require.Equal(t, msg.Symbol, msg2.Symbol)
		require.Equal(t, msg.TokenID, msg2.TokenID)
	}

	{
		msg := NewMsgAttach(nil, "symbol", "item0001", "item0002")
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidAddress("From cannot be empty").Error())

		msg = NewMsgAttach(addr, "s", "item0001", "item0002")
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidCoins("s").Error())

		msg = NewMsgAttach(addr, "symbol", "1", "item0002")
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidCoins("1").Error())

		msg = NewMsgAttach(addr, "symbol", "item0001", "2")
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidCoins("2").Error())

		msg = NewMsgAttach(addr, "symbol", "item0001", "item0001")
		require.EqualError(t, msg.ValidateBasic(), ErrCannotAttachToItself(DefaultCodespace, "symbol"+"item0001").Error())
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

		require.Equal(t, msg.From, msg2.From)
		require.Equal(t, msg.To, msg2.To)
		require.Equal(t, msg.Symbol, msg2.Symbol)
		require.Equal(t, msg.TokenID, msg2.TokenID)
	}

	{
		msg := NewMsgDetach(nil, addr2, "symbol", "item0001")
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidAddress("From cannot be empty").Error())

		msg = NewMsgDetach(addr, nil, "symbol", "item0001")
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidAddress("To cannot be empty").Error())

		msg = NewMsgDetach(addr, addr2, "s", "item0001")
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidCoins("s").Error())

		msg = NewMsgDetach(addr, addr2, "symbol", "1")
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidCoins("1").Error())
	}

	{
		msg := NewMsgAttachFrom(addr, addr2, "symbol", "item0001", "item0002")
		require.Equal(t, "attach_from", msg.Type())
		require.Equal(t, "token", msg.Route())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		require.Equal(t, addr, msg.GetSigners()[0])
		require.NoError(t, msg.ValidateBasic())

		b := msg.GetSignBytes()

		msg2 := MsgAttachFrom{}

		err := cdc.UnmarshalJSON(b, &msg2)
		require.NoError(t, err)

		require.Equal(t, msg.Proxy, msg2.Proxy)
		require.Equal(t, msg.From, msg2.From)
		require.Equal(t, msg.ToTokenID, msg2.ToTokenID)
		require.Equal(t, msg.Symbol, msg2.Symbol)
		require.Equal(t, msg.TokenID, msg2.TokenID)
	}

	{
		msg := NewMsgAttachFrom(nil, addr2, "symbol", "item0001", "item0002")
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidAddress("Proxy cannot be empty").Error())

		msg = NewMsgAttachFrom(addr, nil, "symbol", "item0001", "item0002")
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidAddress("From cannot be empty").Error())

		msg = NewMsgAttachFrom(addr, addr2, "s", "item0001", "item0002")
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidCoins("s").Error())

		msg = NewMsgAttachFrom(addr, addr2, "symbol", "1", "item0002")
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidCoins("1").Error())

		msg = NewMsgAttachFrom(addr, addr2, "symbol", "item0001", "2")
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidCoins("2").Error())

		msg = NewMsgAttachFrom(addr, addr2, "symbol", "item0001", "item0001")
		require.EqualError(t, msg.ValidateBasic(), ErrCannotAttachToItself(DefaultCodespace, "symbolitem0001").Error())
	}

	{
		msg := NewMsgDetachFrom(addr, addr2, addr2, "symbol", "item0001")
		require.Equal(t, "detach_from", msg.Type())
		require.Equal(t, "token", msg.Route())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		require.Equal(t, addr, msg.GetSigners()[0])
		require.NoError(t, msg.ValidateBasic())

		b := msg.GetSignBytes()

		msg2 := MsgDetachFrom{}

		err := cdc.UnmarshalJSON(b, &msg2)
		require.NoError(t, err)

		require.Equal(t, msg.Proxy, msg2.Proxy)
		require.Equal(t, msg.From, msg2.From)
		require.Equal(t, msg.To, msg2.To)
		require.Equal(t, msg.Symbol, msg2.Symbol)
		require.Equal(t, msg.TokenID, msg2.TokenID)
	}

	{
		msg := NewMsgDetachFrom(nil, addr2, addr2, "symbol", "item0001")
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidAddress("Proxy cannot be empty").Error())

		msg = NewMsgDetachFrom(addr, nil, addr2, "symbol", "item0001")
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidAddress("From cannot be empty").Error())

		msg = NewMsgDetachFrom(addr, addr2, nil, "symbol", "item0001")
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidAddress("To cannot be empty").Error())

		msg = NewMsgDetachFrom(addr, addr2, addr2, "s", "item0001")
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidCoins("s").Error())

		msg = NewMsgDetachFrom(addr, addr2, addr2, "symbol", "1")
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidCoins("1").Error())
	}

	{
		msg := NewMsgApproveCollection(addr, addr2, "symbol")
		require.Equal(t, "approve_collection", msg.Type())
		require.Equal(t, "token", msg.Route())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		require.Equal(t, addr, msg.GetSigners()[0])
		require.NoError(t, msg.ValidateBasic())

		b := msg.GetSignBytes()

		msg2 := MsgApproveCollection{}

		err := cdc.UnmarshalJSON(b, &msg2)
		require.NoError(t, err)

		require.Equal(t, msg.Proxy, msg2.Proxy)
		require.Equal(t, msg.Approver, msg2.Approver)
		require.Equal(t, msg.Symbol, msg2.Symbol)
	}

	{
		msg := NewMsgDisapproveCollection(addr, addr2, "symbol")
		require.Equal(t, "disapprove_collection", msg.Type())
		require.Equal(t, "token", msg.Route())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		require.Equal(t, addr, msg.GetSigners()[0])
		require.NoError(t, msg.ValidateBasic())

		b := msg.GetSignBytes()

		msg2 := MsgDisapproveCollection{}

		err := cdc.UnmarshalJSON(b, &msg2)
		require.NoError(t, err)

		require.Equal(t, msg.Proxy, msg2.Proxy)
		require.Equal(t, msg.Approver, msg2.Approver)
		require.Equal(t, msg.Symbol, msg2.Symbol)
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
