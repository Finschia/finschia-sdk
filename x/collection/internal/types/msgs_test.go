package types

import (
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

func TestMsgBasics(t *testing.T) {
	cdc := ModuleCdc

	{
		length1001String := strings.Repeat("Eng글자日本語はスゲ", 91) // 11 * 91 = 1001
		msg := NewMsgIssueFT(addr1, defaultName, defaultSymbol, length1001String, sdk.NewInt(1), sdk.NewInt(8), true)
		require.EqualError(t, msg.ValidateBasic(), ErrInvalidTokenURILength(DefaultCodespace, length1001String).Error())
	}
	{
		msg := NewMsgIssueFT(addr1, defaultName, defaultSymbol, defaultTokenURI, sdk.NewInt(1), sdk.NewInt(8), true)
		require.Equal(t, "issue_ft", msg.Type())
		require.Equal(t, "collection", msg.Route())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		require.Equal(t, addr1, msg.GetSigners()[0])
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
		msg := NewMsgIssueNFT(addr1, defaultSymbol, defaultName)
		require.Equal(t, "issue_nft", msg.Type())
		require.Equal(t, "collection", msg.Route())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		require.Equal(t, addr1, msg.GetSigners()[0])
		require.NoError(t, msg.ValidateBasic())

		b := msg.GetSignBytes()

		msg2 := MsgIssueNFT{}

		err := cdc.UnmarshalJSON(b, &msg2)
		require.NoError(t, err)

		require.Equal(t, msg.Symbol, msg2.Symbol)
		require.Equal(t, msg.Owner, msg2.Owner)
		require.Equal(t, msg.Name, msg2.Name)
	}
	{
		msg := NewMsgMintNFT(addr1, addr1, defaultName, defaultSymbol, defaultTokenURI, defaultTokenType)
		require.Equal(t, "mint_nft", msg.Type())
		require.Equal(t, "collection", msg.Route())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		require.Equal(t, addr1, msg.GetSigners()[0])
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
		msg := NewMsgBurnNFT(addr1, defaultSymbol, defaultTokenID1)
		require.Equal(t, "burn_nft", msg.Type())
		require.Equal(t, "collection", msg.Route())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		require.Equal(t, addr1, msg.GetSigners()[0])
		require.NoError(t, msg.ValidateBasic())

		b := msg.GetSignBytes()

		msg2 := MsgBurnNFT{}

		err := cdc.UnmarshalJSON(b, &msg2)
		require.NoError(t, err)

		require.Equal(t, msg.Symbol, msg2.Symbol)
		require.Equal(t, msg.From, msg2.From)
		require.Equal(t, msg.TokenIDs, msg2.TokenIDs)
	}
	{
		addr2 := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
		msg := NewMsgGrantPermission(addr1, addr2, Permission{Action: "issue", Resource: "resource"})
		require.Equal(t, "grant_perm", msg.Type())
		require.Equal(t, "collection", msg.Route())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		require.Equal(t, addr1, msg.GetSigners()[0])
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
		msg := NewMsgRevokePermission(addr1, Permission{Action: "issue", Resource: "resource"})
		require.Equal(t, "revoke_perm", msg.Type())
		require.Equal(t, "collection", msg.Route())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		require.Equal(t, addr1, msg.GetSigners()[0])
		require.NoError(t, msg.ValidateBasic())

		b := msg.GetSignBytes()

		msg2 := MsgRevokePermission{}

		err := cdc.UnmarshalJSON(b, &msg2)
		require.NoError(t, err)

		require.Equal(t, msg.From, msg2.From)
		require.Equal(t, msg.Permission, msg2.Permission)
	}
	{
		msg := NewMsgTransferFT(addr1, addr2, defaultSymbol, NewCoin(defaultTokenIDFT, sdk.NewInt(defaultAmount)))
		require.Equal(t, "transfer_ft", msg.Type())
		require.Equal(t, "collection", msg.Route())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		require.Equal(t, addr1, msg.GetSigners()[0])
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
		msg := NewMsgTransferFT(nil, addr2, defaultSymbol, NewCoin(defaultTokenIDFT, sdk.NewInt(defaultAmount)))
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidAddress("From cannot be empty").Error())

		msg = NewMsgTransferFT(addr1, nil, defaultSymbol, NewCoin(defaultTokenIDFT, sdk.NewInt(defaultAmount)))
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidAddress("To cannot be empty").Error())

		msg = NewMsgTransferFT(addr1, addr2, "", NewCoin(defaultTokenIDFT, sdk.NewInt(defaultAmount)))
		require.EqualError(t, msg.ValidateBasic(), ErrInvalidTokenSymbol(DefaultCodespace, "symbol [] mismatched to [^[a-z][a-z0-9]{5,7}$]").Error())

		require.Panics(t, func() {
			NewMsgTransferFT(addr1, addr2, defaultSymbol, NewCoin("1", sdk.NewInt(defaultAmount)))
		}, "")

		require.Panics(t, func() {
			NewMsgTransferFT(addr1, addr2, defaultSymbol, NewCoin("1", sdk.NewInt(-1*defaultAmount)))
		}, "")
	}

	{
		msg := NewMsgTransferNFT(addr1, addr2, defaultSymbol, defaultTokenID1)
		require.Equal(t, "transfer_nft", msg.Type())
		require.Equal(t, "collection", msg.Route())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		require.Equal(t, addr1, msg.GetSigners()[0])
		require.NoError(t, msg.ValidateBasic())

		b := msg.GetSignBytes()

		msg2 := MsgTransferNFT{}

		err := cdc.UnmarshalJSON(b, &msg2)
		require.NoError(t, err)

		require.Equal(t, msg.From, msg2.From)
		require.Equal(t, msg.To, msg2.To)
		require.Equal(t, msg.Symbol, msg2.Symbol)
		require.Equal(t, msg.TokenIDs, msg2.TokenIDs)
	}

	{
		msg := NewMsgTransferNFT(nil, addr2, defaultSymbol, defaultTokenID1)
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidAddress("From cannot be empty").Error())

		msg = NewMsgTransferNFT(addr1, nil, defaultSymbol, defaultTokenID1)
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidAddress("To cannot be empty").Error())

		msg = NewMsgTransferNFT(addr1, addr2, "", defaultTokenID1)
		require.EqualError(t, msg.ValidateBasic(), ErrInvalidTokenSymbol(DefaultCodespace, "symbol [] mismatched to [^[a-z][a-z0-9]{5,7}$]").Error())

		msg = NewMsgTransferNFT(addr1, addr2, defaultSymbol, "1")
		require.EqualError(t, msg.ValidateBasic(), ErrInvalidTokenID(DefaultCodespace, "symbol [1] mismatched to [^[a-z0-9]{16}$]").Error())
	}

	{
		msg := NewMsgTransferFTFrom(addr1, addr2, addr2, defaultSymbol, NewCoin(defaultTokenIDFT, sdk.NewInt(defaultAmount)))
		require.Equal(t, "transfer_ft_from", msg.Type())
		require.Equal(t, "collection", msg.Route())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		require.Equal(t, addr1, msg.GetSigners()[0])
		require.NoError(t, msg.ValidateBasic())

		b := msg.GetSignBytes()

		msg2 := MsgTransferFTFrom{}

		err := cdc.UnmarshalJSON(b, &msg2)
		require.NoError(t, err)

		require.Equal(t, msg.Proxy, msg2.Proxy)
		require.Equal(t, msg.From, msg2.From)
		require.Equal(t, msg.To, msg2.To)
		require.Equal(t, msg.Symbol, msg2.Symbol)
		require.Equal(t, msg.Amount, msg2.Amount)
	}

	{
		msg := NewMsgTransferFTFrom(nil, addr2, addr2, defaultSymbol, NewCoin(defaultTokenIDFT, sdk.NewInt(defaultAmount)))
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidAddress("Proxy cannot be empty").Error())

		msg = NewMsgTransferFTFrom(addr1, nil, addr2, defaultSymbol, NewCoin(defaultTokenIDFT, sdk.NewInt(defaultAmount)))
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidAddress("From cannot be empty").Error())

		msg = NewMsgTransferFTFrom(addr1, addr2, nil, defaultSymbol, NewCoin(defaultTokenIDFT, sdk.NewInt(defaultAmount)))
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidAddress("To cannot be empty").Error())

		require.Panics(t, func() {
			NewMsgTransferFT(addr1, addr2, defaultSymbol, NewCoin("1", sdk.NewInt(defaultAmount)))
		}, "")

		require.Panics(t, func() {
			NewMsgTransferFT(addr1, addr2, defaultSymbol, NewCoin("1", sdk.NewInt(-1*defaultAmount)))
		}, "")
	}
	//nolint:dupl
	{
		msg := NewMsgTransferNFTFrom(addr1, addr2, addr2, defaultSymbol, defaultTokenID1)
		require.Equal(t, "transfer_nft_from", msg.Type())
		require.Equal(t, "collection", msg.Route())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		require.Equal(t, addr1, msg.GetSigners()[0])
		require.NoError(t, msg.ValidateBasic())

		b := msg.GetSignBytes()

		msg2 := MsgTransferNFTFrom{}

		err := cdc.UnmarshalJSON(b, &msg2)
		require.NoError(t, err)

		require.Equal(t, msg.Proxy, msg2.Proxy)
		require.Equal(t, msg.From, msg2.From)
		require.Equal(t, msg.To, msg2.To)
		require.Equal(t, msg.Symbol, msg2.Symbol)
		require.Equal(t, msg.TokenIDs, msg2.TokenIDs)
	}

	{
		msg := NewMsgTransferNFTFrom(nil, addr2, addr2, defaultSymbol, defaultTokenIDFT)
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidAddress("Proxy cannot be empty").Error())

		msg = NewMsgTransferNFTFrom(addr1, nil, addr2, defaultSymbol, defaultTokenIDFT)
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidAddress("From cannot be empty").Error())

		msg = NewMsgTransferNFTFrom(addr1, addr2, nil, defaultSymbol, defaultTokenIDFT)
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidAddress("To cannot be empty").Error())

		msg = NewMsgTransferNFTFrom(addr1, addr2, addr2, defaultSymbol, "1")
		require.EqualError(t, msg.ValidateBasic(), ErrInvalidTokenID(DefaultCodespace, "symbol [1] mismatched to [^[a-z0-9]{16}$]").Error())
	}

	{
		msg := NewMsgAttach(addr1, defaultSymbol, defaultTokenID1, defaultTokenID2)
		require.Equal(t, "attach", msg.Type())
		require.Equal(t, "collection", msg.Route())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		require.Equal(t, addr1, msg.GetSigners()[0])
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
		msg := NewMsgAttach(nil, defaultSymbol, defaultTokenID1, defaultTokenID2)
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidAddress("From cannot be empty").Error())

		msg = NewMsgAttach(addr1, "s", defaultTokenID1, defaultTokenID2)
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidCoins("s").Error())

		msg = NewMsgAttach(addr1, defaultSymbol, "1", defaultTokenID2)
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidCoins("1").Error())

		msg = NewMsgAttach(addr1, defaultSymbol, defaultTokenID1, "2")
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidCoins("2").Error())

		msg = NewMsgAttach(addr1, defaultSymbol, defaultTokenID1, defaultTokenID1)
		require.EqualError(t, msg.ValidateBasic(), ErrCannotAttachToItself(DefaultCodespace, defaultTokenID1).Error())
	}

	{
		msg := NewMsgDetach(addr1, defaultSymbol, defaultTokenID1)
		require.Equal(t, "detach", msg.Type())
		require.Equal(t, "collection", msg.Route())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		require.Equal(t, addr1, msg.GetSigners()[0])
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
		msg := NewMsgDetach(nil, defaultSymbol, "item0001")
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidAddress("From cannot be empty").Error())

		msg = NewMsgDetach(addr1, "s", "item0001")
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidCoins("s").Error())

		msg = NewMsgDetach(addr1, defaultSymbol, "1")
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidCoins("1").Error())
	}
	//nolint:dupl
	{
		msg := NewMsgAttachFrom(addr1, addr2, defaultSymbol, defaultTokenID1, defaultTokenID2)
		require.Equal(t, "attach_from", msg.Type())
		require.Equal(t, "collection", msg.Route())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		require.Equal(t, addr1, msg.GetSigners()[0])
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
		msg := NewMsgAttachFrom(nil, addr2, defaultSymbol, defaultTokenID1, defaultTokenID2)
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidAddress("Proxy cannot be empty").Error())

		msg = NewMsgAttachFrom(addr1, nil, defaultSymbol, defaultTokenID1, defaultTokenID2)
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidAddress("From cannot be empty").Error())

		msg = NewMsgAttachFrom(addr1, addr2, "s", defaultTokenID1, defaultTokenID2)
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidCoins("s").Error())

		msg = NewMsgAttachFrom(addr1, addr2, defaultSymbol, "1", defaultTokenID2)
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidCoins("1").Error())

		msg = NewMsgAttachFrom(addr1, addr2, defaultSymbol, defaultTokenID1, "2")
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidCoins("2").Error())

		msg = NewMsgAttachFrom(addr1, addr2, defaultSymbol, defaultTokenID1, defaultTokenID1)
		require.EqualError(t, msg.ValidateBasic(), ErrCannotAttachToItself(DefaultCodespace, defaultTokenID1).Error())
	}
	//nolint:dupl
	{
		msg := NewMsgDetachFrom(addr1, addr2, defaultSymbol, defaultTokenID1)
		require.Equal(t, "detach_from", msg.Type())
		require.Equal(t, "collection", msg.Route())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		require.Equal(t, addr1, msg.GetSigners()[0])
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
		msg := NewMsgDetachFrom(nil, addr2, defaultSymbol, defaultTokenID1)
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidAddress("Proxy cannot be empty").Error())

		msg = NewMsgDetachFrom(addr1, nil, defaultSymbol, defaultTokenID1)
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidAddress("From cannot be empty").Error())

		msg = NewMsgDetachFrom(addr1, addr2, "s", defaultTokenID1)
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidCoins("s").Error())

		msg = NewMsgDetachFrom(addr1, addr2, defaultSymbol, "1")
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidCoins("1").Error())
	}

	{
		msg := NewMsgApprove(addr1, addr2, defaultSymbol)
		require.Equal(t, "approve_collection", msg.Type())
		require.Equal(t, "collection", msg.Route())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		require.Equal(t, addr1, msg.GetSigners()[0])
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
		msg := NewMsgDisapprove(addr1, addr2, defaultSymbol)
		require.Equal(t, "disapprove_collection", msg.Type())
		require.Equal(t, "collection", msg.Route())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		require.Equal(t, addr1, msg.GetSigners()[0])
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
		msg := NewMsgBurnFTFrom(defaultSymbol, addr1, addr2, OneCoin(defaultTokenIDFT))
		require.Equal(t, "burn_ft_from", msg.Type())
		require.Equal(t, "collection", msg.Route())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		require.Equal(t, addr1, msg.GetSigners()[0])
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
		msg := NewMsgBurnFTFrom(defaultSymbol, addr1, addr1, OneCoin(defaultTokenIDFT))
		require.EqualError(t, msg.ValidateBasic(), ErrApproverProxySame(DefaultCodespace, addr1.String()).Error())

		msg = NewMsgBurnFTFrom(defaultSymbol, nil, addr1, OneCoin(defaultTokenIDFT))
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidAddress("Proxy cannot be empty").Error())

		msg = NewMsgBurnFTFrom(defaultSymbol, addr1, nil, OneCoin(defaultTokenIDFT))
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidAddress("From cannot be empty").Error())
	}

	{
		msg := NewMsgBurnNFTFrom(addr1, addr2, defaultSymbol, defaultTokenID1)
		require.Equal(t, "burn_nft_from", msg.Type())
		require.Equal(t, "collection", msg.Route())
		require.Equal(t, sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
		require.Equal(t, addr1, msg.GetSigners()[0])
		require.NoError(t, msg.ValidateBasic())

		b := msg.GetSignBytes()

		msg2 := MsgBurnNFTFrom{}

		err := cdc.UnmarshalJSON(b, &msg2)
		require.NoError(t, err)

		require.Equal(t, msg.Proxy, msg2.Proxy)
		require.Equal(t, msg.From, msg2.From)
		require.Equal(t, msg.Symbol, msg2.Symbol)
		require.Equal(t, msg.TokenIDs, msg2.TokenIDs)
	}

	{
		msg := NewMsgBurnNFTFrom(addr1, addr1, defaultSymbol, defaultTokenID1)
		require.EqualError(t, msg.ValidateBasic(), ErrApproverProxySame(DefaultCodespace, addr1.String()).Error())

		msg = NewMsgBurnNFTFrom(nil, addr1, defaultSymbol, defaultTokenID1)
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidAddress("Proxy cannot be empty").Error())

		msg = NewMsgBurnNFTFrom(addr1, nil, defaultSymbol, defaultTokenID1)
		require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidAddress("From cannot be empty").Error())
	}
}

func TestMsgModifyTokenURI_ValidateBasicMsgBasics(t *testing.T) {
	cdc := ModuleCdc
	addr := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
	const ModifyActionName = "modify_token"
	t.Log("normal case")
	{
		msg := NewMsgModifyTokenURI(addr, defaultSymbol, defaultTokenURI, defaultTokenID1)
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
		msg := NewMsgModifyTokenURI(addr, "", defaultTokenURI, defaultTokenID1)
		require.Error(t, msg.ValidateBasic())
	}
	t.Log("empty owner")
	{
		msg := NewMsgModifyTokenURI(nil, defaultSymbol, defaultTokenURI, defaultTokenID1)
		require.Error(t, msg.ValidateBasic())
	}
	t.Log("tokenURI too long")
	{
		length1001String := strings.Repeat("Eng글자日本語はスゲ", 91) // 11 * 91 = 1001
		msg := NewMsgModifyTokenURI(addr, defaultSymbol, length1001String, defaultTokenID1)
		require.EqualError(t, msg.ValidateBasic(), ErrInvalidTokenURILength(DefaultCodespace, length1001String).Error())
	}
	t.Log("invalid symbol found")
	{
		msg := NewMsgModifyTokenURI(addr, "invalidsymbol2198721987", defaultTokenURI, defaultTokenID1)
		require.Error(t, msg.ValidateBasic())
	}
	t.Log("invalid tokenid found")
	{
		msg := NewMsgModifyTokenURI(addr, defaultSymbol, defaultTokenURI, "tokenid")
		require.Error(t, msg.ValidateBasic())
	}
}
