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
		msg := NewMsgIssueCFT(addr, "name", "symb"+addrSuffix, "tokenuri", sdk.NewInt(1), sdk.NewInt(8), true)
		require.Equal(t, "issue_cft", msg.Type())
		require.Equal(t, "collection", msg.Route())
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
		msg := NewMsgIssueCNFT(addr, "symb"+addrSuffix)
		require.Equal(t, "issue_cnft", msg.Type())
		require.Equal(t, "collection", msg.Route())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		require.Equal(t, addr, msg.GetSigners()[0])
		require.NoError(t, msg.ValidateBasic())

		b := msg.GetSignBytes()

		msg2 := MsgIssueCNFT{}

		err := cdc.UnmarshalJSON(b, &msg2)
		require.NoError(t, err)

		require.Equal(t, msg.Symbol, msg2.Symbol)
		require.Equal(t, msg.Owner, msg2.Owner)
	}
	{
		msg := NewMsgMintCNFT(addr, addr, "name", "symb"+addrSuffix, "tokenuri", "toke")
		require.Equal(t, "mint_cnft", msg.Type())
		require.Equal(t, "collection", msg.Route())
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
		msg := NewMsgBurnCNFT(addr, "symb"+addrSuffix, "10010001")
		require.Equal(t, "burn_cnft", msg.Type())
		require.Equal(t, "collection", msg.Route())
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
		addr2 := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
		msg := NewMsgGrantPermission(addr, addr2, Permission{Action: "issue", Resource: "resource"})
		require.Equal(t, "grant_perm", msg.Type())
		require.Equal(t, "collection", msg.Route())
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
		require.Equal(t, "collection", msg.Route())
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
		msg := NewMsgTransferCFT(addr, addr2, "symbol", "00000001", sdk.NewInt(1))
		require.Equal(t, "transfer_cft", msg.Type())
		require.Equal(t, "collection", msg.Route())
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
		require.Equal(t, "collection", msg.Route())
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
		require.Equal(t, "collection", msg.Route())
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
	//nolint:dupl
	{
		msg := NewMsgTransferCNFTFrom(addr, addr2, addr2, "symbol", "00000001")
		require.Equal(t, "transfer_cnft_from", msg.Type())
		require.Equal(t, "collection", msg.Route())
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
		require.Equal(t, "collection", msg.Route())
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
		require.Equal(t, "collection", msg.Route())
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
		require.Equal(t, "collection", msg.Route())
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
	//nolint:dupl
	{
		msg := NewMsgDetachFrom(addr, addr2, addr2, "symbol", "item0001")
		require.Equal(t, "detach_from", msg.Type())
		require.Equal(t, "collection", msg.Route())
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
		msg := NewMsgApprove(addr, addr2, "symbol")
		require.Equal(t, "approve_collection", msg.Type())
		require.Equal(t, "collection", msg.Route())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		require.Equal(t, addr, msg.GetSigners()[0])
		require.NoError(t, msg.ValidateBasic())

		b := msg.GetSignBytes()

		msg2 := MsgApprove{}

		err := cdc.UnmarshalJSON(b, &msg2)
		require.NoError(t, err)

		require.Equal(t, msg.Proxy, msg2.Proxy)
		require.Equal(t, msg.Approver, msg2.Approver)
		require.Equal(t, msg.Symbol, msg2.Symbol)
	}

	{
		msg := NewMsgDisapprove(addr, addr2, "symbol")
		require.Equal(t, "disapprove_collection", msg.Type())
		require.Equal(t, "collection", msg.Route())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		require.Equal(t, addr, msg.GetSigners()[0])
		require.NoError(t, msg.ValidateBasic())

		b := msg.GetSignBytes()

		msg2 := MsgDisapprove{}

		err := cdc.UnmarshalJSON(b, &msg2)
		require.NoError(t, err)

		require.Equal(t, msg.Proxy, msg2.Proxy)
		require.Equal(t, msg.Approver, msg2.Approver)
		require.Equal(t, msg.Symbol, msg2.Symbol)
	}

	{
		msg := NewMsgBurnCFTFrom(addr, addr2, types.NewCoinWithTokenIDs(types.NewCoinWithTokenID("symbol", "00000001", sdk.NewInt(1))))
		require.Equal(t, "burn_cft_from", msg.Type())
		require.Equal(t, "collection", msg.Route())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		require.Equal(t, addr, msg.GetSigners()[0])
		require.NoError(t, msg.ValidateBasic())

		b := msg.GetSignBytes()

		msg2 := MsgBurnCFTFrom{}

		err := cdc.UnmarshalJSON(b, &msg2)
		require.NoError(t, err)

		require.Equal(t, msg.Proxy, msg2.Proxy)
		require.Equal(t, msg.From, msg2.From)
		require.Equal(t, msg.Amount, msg2.Amount)
	}

	{
		msg := NewMsgBurnCFTFrom(addr, addr, types.NewCoinWithTokenIDs(types.NewCoinWithTokenID("symbol", "00000001", sdk.NewInt(1))))
		require.EqualError(t, msg.ValidateBasic(), ErrApproverProxySame(DefaultCodespace, addr.String()).Error())

		msg = NewMsgBurnCFTFrom(nil, addr, types.NewCoinWithTokenIDs(types.NewCoinWithTokenID("symbol", "00000001", sdk.NewInt(1))))
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidAddress("Proxy cannot be empty").Error())

		msg = NewMsgBurnCFTFrom(addr, nil, types.NewCoinWithTokenIDs(types.NewCoinWithTokenID("symbol", "00000001", sdk.NewInt(1))))
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidAddress("From cannot be empty").Error())
	}

	{
		msg := NewMsgBurnCNFTFrom(addr, addr2, "symbol", "item0001")
		require.Equal(t, "burn_cnft_from", msg.Type())
		require.Equal(t, "collection", msg.Route())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		require.Equal(t, addr, msg.GetSigners()[0])
		require.NoError(t, msg.ValidateBasic())

		b := msg.GetSignBytes()

		msg2 := MsgBurnCNFTFrom{}

		err := cdc.UnmarshalJSON(b, &msg2)
		require.NoError(t, err)

		require.Equal(t, msg.Proxy, msg2.Proxy)
		require.Equal(t, msg.From, msg2.From)
		require.Equal(t, msg.Symbol, msg2.Symbol)
		require.Equal(t, msg.TokenID, msg2.TokenID)
	}

	{
		msg := NewMsgBurnCNFTFrom(addr, addr, "symbol", "item0001")
		require.EqualError(t, msg.ValidateBasic(), ErrApproverProxySame(DefaultCodespace, addr.String()).Error())

		msg = NewMsgBurnCNFTFrom(nil, addr, "symbol", "item0001")
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidAddress("Proxy cannot be empty").Error())

		msg = NewMsgBurnCNFTFrom(addr, nil, "symbol", "item0001")
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidAddress("From cannot be empty").Error())
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
