package types

import (
	"strings"
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
		length1001String := strings.Repeat("Eng글자日本語はスゲ", 91) // 11 * 91 = 1001
		msg := NewMsgIssueFT(addr, "name", "symb"+addrSuffix, length1001String, sdk.NewInt(1), sdk.NewInt(8), true)
		require.EqualError(t, msg.ValidateBasic(), ErrInvalidTokenURILength(DefaultCodespace, length1001String).Error())
	}
	{
		msg := NewMsgIssueFT(addr, "name", "symb"+addrSuffix, "tokenuri", sdk.NewInt(1), sdk.NewInt(8), true)
		require.Equal(t, "issue_ft", msg.Type())
		require.Equal(t, "collection", msg.Route())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		require.Equal(t, addr, msg.GetSigners()[0])
		require.NoError(t, msg.ValidateBasic())

		b := msg.GetSignBytes()

		msg2 := MsgIssueFT{}

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
		msg := NewMsgIssueNFT(addr, "symb"+addrSuffix)
		require.Equal(t, "issue_nft", msg.Type())
		require.Equal(t, "collection", msg.Route())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		require.Equal(t, addr, msg.GetSigners()[0])
		require.NoError(t, msg.ValidateBasic())

		b := msg.GetSignBytes()

		msg2 := MsgIssueNFT{}

		err := cdc.UnmarshalJSON(b, &msg2)
		require.NoError(t, err)

		require.Equal(t, msg.Symbol, msg2.Symbol)
		require.Equal(t, msg.Owner, msg2.Owner)
	}
	{
		msg := NewMsgMintNFT(addr, addr, "name", "symb"+addrSuffix, "tokenuri", "toke")
		require.Equal(t, "mint_nft", msg.Type())
		require.Equal(t, "collection", msg.Route())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		require.Equal(t, addr, msg.GetSigners()[0])
		require.NoError(t, msg.ValidateBasic())

		b := msg.GetSignBytes()

		msg2 := MsgMintNFT{}

		err := cdc.UnmarshalJSON(b, &msg2)
		require.NoError(t, err)

		require.Equal(t, msg.Name, msg2.Name)
		require.Equal(t, msg.Symbol, msg2.Symbol)
		require.Equal(t, msg.TokenURI, msg2.TokenURI)
		require.Equal(t, msg.From, msg2.From)
		require.Equal(t, msg.TokenType, msg2.TokenType)
	}
	{
		msg := NewMsgBurnNFT(addr, "symb"+addrSuffix, "10010001")
		require.Equal(t, "burn_nft", msg.Type())
		require.Equal(t, "collection", msg.Route())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		require.Equal(t, addr, msg.GetSigners()[0])
		require.NoError(t, msg.ValidateBasic())

		b := msg.GetSignBytes()

		msg2 := MsgBurnNFT{}

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
		msg := NewMsgTransferFT(addr, addr2, "symbol", "00000001", sdk.NewInt(1))
		require.Equal(t, "transfer_ft", msg.Type())
		require.Equal(t, "collection", msg.Route())
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
		require.Equal(t, msg.TokenID, msg2.TokenID)
		require.Equal(t, msg.Amount, msg2.Amount)
	}

	{
		msg := NewMsgTransferFT(nil, addr2, "symbol", "00000001", sdk.NewInt(1))
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidAddress("From cannot be empty").Error())

		msg = NewMsgTransferFT(addr, nil, "symbol", "00000001", sdk.NewInt(1))
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidAddress("To cannot be empty").Error())

		msg = NewMsgTransferFT(addr, addr2, "", "00000001", sdk.NewInt(1))
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidCoins("Token symbol is empty").Error())

		msg = NewMsgTransferFT(addr, addr2, "symbol", "1", sdk.NewInt(1))
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidCoins("Only user defined token is possible: symbol1").Error())

		msg = NewMsgTransferFT(addr, addr2, "symbol", "00000001", sdk.NewInt(-1))
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInsufficientCoins("send amount must be positive").Error())
	}

	{
		msg := NewMsgTransferNFT(addr, addr2, "symbol", "00000001")
		require.Equal(t, "transfer_nft", msg.Type())
		require.Equal(t, "collection", msg.Route())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		require.Equal(t, addr, msg.GetSigners()[0])
		require.NoError(t, msg.ValidateBasic())

		b := msg.GetSignBytes()

		msg2 := MsgTransferNFT{}

		err := cdc.UnmarshalJSON(b, &msg2)
		require.NoError(t, err)

		require.Equal(t, msg.From, msg2.From)
		require.Equal(t, msg.To, msg2.To)
		require.Equal(t, msg.Symbol, msg2.Symbol)
		require.Equal(t, msg.TokenID, msg2.TokenID)
	}

	{
		msg := NewMsgTransferNFT(nil, addr2, "symbol", "00000001")
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidAddress("From cannot be empty").Error())

		msg = NewMsgTransferNFT(addr, nil, "symbol", "00000001")
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidAddress("To cannot be empty").Error())

		msg = NewMsgTransferNFT(addr, addr2, "", "00000001")
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidCoins("Token symbol is empty").Error())

		msg = NewMsgTransferNFT(addr, addr2, "symbol", "1")
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidCoins("Only user defined token is possible: symbol1").Error())
	}

	{
		msg := NewMsgTransferFTFrom(addr, addr2, addr2, "symbol", "00000001", sdk.NewInt(1))
		require.Equal(t, "transfer_ft_from", msg.Type())
		require.Equal(t, "collection", msg.Route())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		require.Equal(t, addr, msg.GetSigners()[0])
		require.NoError(t, msg.ValidateBasic())

		b := msg.GetSignBytes()

		msg2 := MsgTransferFTFrom{}

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
		msg := NewMsgTransferFTFrom(nil, addr2, addr2, "symbol", "00000001", sdk.NewInt(1))
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidAddress("Proxy cannot be empty").Error())

		msg = NewMsgTransferFTFrom(addr, nil, addr2, "symbol", "00000001", sdk.NewInt(1))
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidAddress("From cannot be empty").Error())

		msg = NewMsgTransferFTFrom(addr, addr2, nil, "symbol", "00000001", sdk.NewInt(1))
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidAddress("To cannot be empty").Error())

		msg = NewMsgTransferFTFrom(addr, addr2, addr2, "symbol", "1", sdk.NewInt(1))
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidCoins("Only user defined token is possible: symbol1").Error())

		msg = NewMsgTransferFTFrom(addr, addr2, addr2, "symbol", "00000001", sdk.NewInt(-1))
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInsufficientCoins("send amount must be positive").Error())
	}
	//nolint:dupl
	{
		msg := NewMsgTransferNFTFrom(addr, addr2, addr2, "symbol", "00000001")
		require.Equal(t, "transfer_nft_from", msg.Type())
		require.Equal(t, "collection", msg.Route())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		require.Equal(t, addr, msg.GetSigners()[0])
		require.NoError(t, msg.ValidateBasic())

		b := msg.GetSignBytes()

		msg2 := MsgTransferNFTFrom{}

		err := cdc.UnmarshalJSON(b, &msg2)
		require.NoError(t, err)

		require.Equal(t, msg.Proxy, msg2.Proxy)
		require.Equal(t, msg.From, msg2.From)
		require.Equal(t, msg.To, msg2.To)
		require.Equal(t, msg.Symbol, msg2.Symbol)
		require.Equal(t, msg.TokenID, msg2.TokenID)
	}

	{
		msg := NewMsgTransferNFTFrom(nil, addr2, addr2, "symbol", "00000001")
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidAddress("Proxy cannot be empty").Error())

		msg = NewMsgTransferNFTFrom(addr, nil, addr2, "symbol", "00000001")
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidAddress("From cannot be empty").Error())

		msg = NewMsgTransferNFTFrom(addr, addr2, nil, "symbol", "00000001")
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidAddress("To cannot be empty").Error())

		msg = NewMsgTransferNFTFrom(addr, addr2, addr2, "symbol", "1")
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
		msg := NewMsgDetach(addr, "symbol", "item0001")
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
		require.Equal(t, msg.Symbol, msg2.Symbol)
		require.Equal(t, msg.TokenID, msg2.TokenID)
	}

	{
		msg := NewMsgDetach(nil, "symbol", "item0001")
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidAddress("From cannot be empty").Error())

		msg = NewMsgDetach(addr, "s", "item0001")
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidCoins("s").Error())

		msg = NewMsgDetach(addr, "symbol", "1")
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
		msg := NewMsgDetachFrom(addr, addr2, "symbol", "item0001")
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
		require.Equal(t, msg.Symbol, msg2.Symbol)
		require.Equal(t, msg.TokenID, msg2.TokenID)
	}

	{
		msg := NewMsgDetachFrom(nil, addr2, "symbol", "item0001")
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidAddress("Proxy cannot be empty").Error())

		msg = NewMsgDetachFrom(addr, nil, "symbol", "item0001")
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidAddress("From cannot be empty").Error())

		msg = NewMsgDetachFrom(addr, addr2, "s", "item0001")
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidCoins("s").Error())

		msg = NewMsgDetachFrom(addr, addr2, "symbol", "1")
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
		msg := NewMsgBurnFTFrom(addr, addr2, types.NewCoinWithTokenIDs(types.NewCoinWithTokenID("symbol", "00000001", sdk.NewInt(1))))
		require.Equal(t, "burn_ft_from", msg.Type())
		require.Equal(t, "collection", msg.Route())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		require.Equal(t, addr, msg.GetSigners()[0])
		require.NoError(t, msg.ValidateBasic())

		b := msg.GetSignBytes()

		msg2 := MsgBurnFTFrom{}

		err := cdc.UnmarshalJSON(b, &msg2)
		require.NoError(t, err)

		require.Equal(t, msg.Proxy, msg2.Proxy)
		require.Equal(t, msg.From, msg2.From)
		require.Equal(t, msg.Amount, msg2.Amount)
	}

	{
		msg := NewMsgBurnFTFrom(addr, addr, types.NewCoinWithTokenIDs(types.NewCoinWithTokenID("symbol", "00000001", sdk.NewInt(1))))
		require.EqualError(t, msg.ValidateBasic(), ErrApproverProxySame(DefaultCodespace, addr.String()).Error())

		msg = NewMsgBurnFTFrom(nil, addr, types.NewCoinWithTokenIDs(types.NewCoinWithTokenID("symbol", "00000001", sdk.NewInt(1))))
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidAddress("Proxy cannot be empty").Error())

		msg = NewMsgBurnFTFrom(addr, nil, types.NewCoinWithTokenIDs(types.NewCoinWithTokenID("symbol", "00000001", sdk.NewInt(1))))
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidAddress("From cannot be empty").Error())
	}

	{
		msg := NewMsgBurnNFTFrom(addr, addr2, "symbol", "item0001")
		require.Equal(t, "burn_nft_from", msg.Type())
		require.Equal(t, "collection", msg.Route())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		require.Equal(t, addr, msg.GetSigners()[0])
		require.NoError(t, msg.ValidateBasic())

		b := msg.GetSignBytes()

		msg2 := MsgBurnNFTFrom{}

		err := cdc.UnmarshalJSON(b, &msg2)
		require.NoError(t, err)

		require.Equal(t, msg.Proxy, msg2.Proxy)
		require.Equal(t, msg.From, msg2.From)
		require.Equal(t, msg.Symbol, msg2.Symbol)
		require.Equal(t, msg.TokenID, msg2.TokenID)
	}

	{
		msg := NewMsgBurnNFTFrom(addr, addr, "symbol", "item0001")
		require.EqualError(t, msg.ValidateBasic(), ErrApproverProxySame(DefaultCodespace, addr.String()).Error())

		msg = NewMsgBurnNFTFrom(nil, addr, "symbol", "item0001")
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidAddress("Proxy cannot be empty").Error())

		msg = NewMsgBurnNFTFrom(addr, nil, "symbol", "item0001")
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
	t.Log("tokenURI too long")
	{
		length1001String := strings.Repeat("Eng글자日本語はスゲ", 91) // 11 * 91 = 1001
		msg := NewMsgModifyTokenURI(addr, "symbol", length1001String, "tokenid0")
		require.EqualError(t, msg.ValidateBasic(), ErrInvalidTokenURILength(DefaultCodespace, length1001String).Error())
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
