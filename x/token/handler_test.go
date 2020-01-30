package token

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	linktype "github.com/line/link/types"
	testCommon "github.com/line/link/x/token/internal/keeper"
	"github.com/line/link/x/token/internal/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"strings"
	"testing"
)

const (
	symbolCony  = "cony"
	symbolBrown = "brown"
	name        = "description"
	tokenuri    = "tokenuri"
)

var verifyEventFunc = func(t *testing.T, expected sdk.Events, actual sdk.Events) {
	require.Equal(t, sdk.StringifyEvents(expected.ToABCIEvents()).String(), sdk.StringifyEvents(actual.ToABCIEvents()).String())
}
var (
	amount   = sdk.NewInt(1000)
	decimals = sdk.NewInt(6)
	priv1    = secp256k1.GenPrivKey()
	addr1    = sdk.AccAddress(priv1.PubKey().Address())
	symbol1  = symbolCony + addr1.String()[len(addr1.String())-3:]
	symbol2  = symbolBrown + addr1.String()[len(addr1.String())-3:]
	tokenId1 = "item0001"
	tokenId2 = "item0002"
	priv2    = secp256k1.GenPrivKey()
	addr2    = sdk.AccAddress(priv2.PubKey().Address())
)

func TestHandlerUnrecognized(t *testing.T) {
	input := testCommon.SetupTestInput(t)
	ctx, keeper := input.Ctx, input.Keeper

	h := NewHandler(keeper)

	res := h(ctx, sdk.NewTestMsg())
	require.False(t, res.IsOK())
	require.True(t, strings.Contains(res.Log, "Unrecognized  Msg type"))
	require.False(t, res.Code.IsOK())
}

type handlerTestSuite struct {
	ti *testCommon.TestInput
}

//This prevents you from being affected by the previous Context.
//There is no need to explicitly clear the previous context every time.
//It means that you can focus on writing your testing code what you want to prove.
func (ht *handlerTestSuite) handleNewMsg(msg sdk.Msg) sdk.Result {
	ht.ti.Ctx = ht.ti.Ctx.WithEventManager(sdk.NewEventManager())
	h := NewHandler(ht.ti.Keeper)
	return h(ht.ti.Ctx, msg)
}

func TestHandlerModifyTokenURI(t *testing.T) {
	h := handlerTestSuite{testCommon.SetupTestInput(t)}
	modifyTokenURI := "modifyTokenURI"

	t.Log("token not exist")
	{
		res := h.handleNewMsg(types.NewMsgModifyTokenURI(addr1, symbol1, modifyTokenURI, "tokenId0"))
		require.False(t, res.Code.IsOK())
		require.Equal(t, types.DefaultCodespace, res.Codespace)
		require.Equal(t, types.CodeCollectionNotExist, res.Code)
		verifyEventFunc(t, nil, res.Events)
	}
	t.Log("modify token for FT")
	{
		require.True(t, h.handleNewMsg(types.NewMsgIssue(name, symbol1, tokenuri, addr1, amount, decimals, true)).Code.IsOK())
		res := h.handleNewMsg(types.NewMsgModifyTokenURI(addr1, symbol1, modifyTokenURI, ""))
		require.True(t, res.Code.IsOK())
		require.Equal(t, types.EventTypeModifyTokenURI, res.Events[0].Type)
		require.Equal(t, modifyTokenURI, string(res.Events[0].Attributes[4].Value))
	}

	t.Log("modify token for NFT")
	{
		require.True(t, h.handleNewMsg(types.NewMsgIssueNFT(name, symbol2, tokenuri, addr1)).Code.IsOK())
		res := h.handleNewMsg(types.NewMsgModifyTokenURI(addr1, symbol2, modifyTokenURI, ""))
		require.True(t, res.Code.IsOK())
		require.Equal(t, types.EventTypeModifyTokenURI, res.Events[0].Type)
		require.Equal(t, modifyTokenURI, string(res.Events[0].Attributes[4].Value))
	}
}

func TestHandlerIssueFT(t *testing.T) {
	input := testCommon.SetupTestInput(t)
	ctx, keeper := input.Ctx, input.Keeper

	h := NewHandler(keeper)

	{
		msg := types.NewMsgIssue(name, symbol1, tokenuri, addr1, amount, decimals, true)
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())
	}
	{
		msg := types.NewMsgIssue(name, symbol1, tokenuri, addr1, amount, decimals, true)
		res := h(ctx, msg)
		require.False(t, res.Code.IsOK())
		require.Equal(t, types.DefaultCodespace, res.Codespace)
		require.Equal(t, types.CodeTokenExist, res.Code)
	}
}

func TestHandlerIssueNFT(t *testing.T) {
	input := testCommon.SetupTestInput(t)
	ctx, keeper := input.Ctx, input.Keeper

	h := NewHandler(keeper)

	{
		msg := types.NewMsgIssueNFT(name, symbol1, tokenuri, addr1)
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())
	}
	{
		msg := types.NewMsgIssueNFT(name, symbol1, tokenuri, addr1)
		res := h(ctx, msg)
		require.False(t, res.Code.IsOK())
		require.Equal(t, types.DefaultCodespace, res.Codespace)
		require.Equal(t, types.CodeTokenExist, res.Code)
	}
}

func TestHandlerIssueFTCollection(t *testing.T) {
	input := testCommon.SetupTestInput(t)
	ctx, keeper := input.Ctx, input.Keeper

	h := NewHandler(keeper)

	{
		msg := types.NewMsgCreateCollection(name, symbol1, addr1)
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())
	}

	{
		msg := types.NewMsgIssueCollection(name, symbol1, tokenuri, addr1, amount, decimals, true, "00000001")
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())
	}
	{
		msg := types.NewMsgIssueCollection(name, symbol1, tokenuri, addr1, amount, decimals, true, "00000001")
		res := h(ctx, msg)
		require.False(t, res.Code.IsOK())
		require.Equal(t, types.DefaultCodespace, res.Codespace)
		require.Equal(t, types.CodeTokenExist, res.Code)
	}
	{
		msg := types.NewMsgIssueCollection(name, symbol1, tokenuri, addr1, amount, decimals, true, "00000002")
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())
	}
	{
		msg := types.NewMsgIssueCollection(name, symbol1, tokenuri, addr2, amount, decimals, true, "00000003")
		res := h(ctx, msg)
		require.False(t, res.Code.IsOK())
		require.Equal(t, types.DefaultCodespace, res.Codespace)
		require.Equal(t, types.CodeTokenPermission, res.Code)
	}

	permission := types.Permission{
		Action:   "issue",
		Resource: symbol1,
	}

	{
		msg := types.NewMsgGrantPermission(addr1, addr2, permission)
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())
	}
	{
		msg := types.NewMsgIssueCollection(name, symbol1, tokenuri, addr2, amount, decimals, true, "00000003")
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())
	}
	{
		msg := types.NewMsgRevokePermission(addr1, permission)
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())
	}
	{
		msg := types.NewMsgIssueCollection(name, symbol1, tokenuri, addr1, amount, decimals, true, "00000004")
		res := h(ctx, msg)
		require.False(t, res.Code.IsOK())
		require.Equal(t, types.DefaultCodespace, res.Codespace)
		require.Equal(t, types.CodeTokenPermission, res.Code)
	}
}

func TestEvents(t *testing.T) {
	input := testCommon.SetupTestInput(t)
	ctx, keeper := input.Ctx, input.Keeper
	h := NewHandler(keeper)
	suffixAddr1 := addr1.String()[len(addr1.String())-3:]
	tokenAddr := sdk.AccAddress(crypto.AddressHash([]byte("token")))

	ctx = ctx.WithEventManager(sdk.NewEventManager())
	{
		symbol := "t01" + suffixAddr1
		msg := types.NewMsgIssue(name, symbol, tokenuri, addr1, amount, decimals, true)
		require.NoError(t, msg.ValidateBasic())
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())

		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("sender", tokenAddr.String())),
			sdk.NewEvent("message", sdk.NewAttribute("module", "token")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("grant_perm_token", sdk.NewAttribute("to", addr1.String())),
			sdk.NewEvent("grant_perm_token", sdk.NewAttribute("perm_resource", symbol)),
			sdk.NewEvent("grant_perm_token", sdk.NewAttribute("perm_action", "mint")),
			sdk.NewEvent(types.EventTypeModifyTokenURIPermToken, sdk.NewAttribute("to", addr1.String())),
			sdk.NewEvent(types.EventTypeModifyTokenURIPermToken, sdk.NewAttribute("perm_resource", symbol)),
			sdk.NewEvent(types.EventTypeModifyTokenURIPermToken, sdk.NewAttribute("perm_action", types.ModifyAction)),
			sdk.NewEvent("issue_token", sdk.NewAttribute("name", name)),
			sdk.NewEvent("issue_token", sdk.NewAttribute("symbol", symbol)),
			sdk.NewEvent("issue_token", sdk.NewAttribute("denom", symbol)),
			sdk.NewEvent("issue_token", sdk.NewAttribute("owner", addr1.String())),
			sdk.NewEvent("issue_token", sdk.NewAttribute("amount", amount.String())),
			sdk.NewEvent("issue_token", sdk.NewAttribute("mintable", "true")),
			sdk.NewEvent("issue_token", sdk.NewAttribute("decimals", decimals.String())),
			sdk.NewEvent("issue_token", sdk.NewAttribute("token_type", "ft")),
			sdk.NewEvent("mint_token", sdk.NewAttribute("amount", amount.String()+symbol)),
			sdk.NewEvent("mint_token", sdk.NewAttribute("to", addr1.String())),
			sdk.NewEvent("transfer", sdk.NewAttribute("recipient", addr1.String())),
			sdk.NewEvent("transfer", sdk.NewAttribute("amount", amount.String()+symbol)),
		}
		verifyEventFunc(t, e, res.Events)
	}
	ctx = ctx.WithEventManager(sdk.NewEventManager())
	{
		symbol := "t01" + suffixAddr1
		msg := types.NewMsgMint(addr1, addr1, sdk.NewCoins(sdk.NewInt64Coin(symbol, amount.Int64())))
		require.NoError(t, msg.ValidateBasic())
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())
		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("sender", tokenAddr.String())),
			sdk.NewEvent("message", sdk.NewAttribute("module", "token")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("mint_token", sdk.NewAttribute("amount", amount.String()+symbol)),
			sdk.NewEvent("mint_token", sdk.NewAttribute("to", addr1.String())),
			sdk.NewEvent("transfer", sdk.NewAttribute("recipient", addr1.String())),
			sdk.NewEvent("transfer", sdk.NewAttribute("amount", amount.String()+symbol)),
		}
		verifyEventFunc(t, e, res.Events)
	}

	ctx = ctx.WithEventManager(sdk.NewEventManager())
	{
		symbol := "t01" + suffixAddr1
		msg := types.NewMsgBurn(addr1, sdk.NewCoins(sdk.NewInt64Coin(symbol, amount.Int64())))
		require.NoError(t, msg.ValidateBasic())
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())
		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("message", sdk.NewAttribute("module", "token")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("burn_token", sdk.NewAttribute("amount", amount.String()+symbol)),
			sdk.NewEvent("burn_token", sdk.NewAttribute("from", addr1.String())),
			sdk.NewEvent("transfer", sdk.NewAttribute("recipient", tokenAddr.String())),
			sdk.NewEvent("transfer", sdk.NewAttribute("amount", amount.String()+symbol)),
		}
		verifyEventFunc(t, e, res.Events)
	}

	ctx = ctx.WithEventManager(sdk.NewEventManager())
	{
		symbol := "t02" + suffixAddr1
		msg := types.NewMsgIssueNFT(name, symbol, tokenuri, addr1)
		require.NoError(t, msg.ValidateBasic())
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())

		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("sender", tokenAddr.String())),
			sdk.NewEvent("message", sdk.NewAttribute("module", "token")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("issue_token", sdk.NewAttribute("name", name)),
			sdk.NewEvent("issue_token", sdk.NewAttribute("symbol", symbol)),
			sdk.NewEvent("issue_token", sdk.NewAttribute("denom", symbol)),
			sdk.NewEvent("issue_token", sdk.NewAttribute("owner", addr1.String())),
			sdk.NewEvent("issue_token", sdk.NewAttribute("token_uri", tokenuri)),
			sdk.NewEvent("issue_token", sdk.NewAttribute("token_type", "nft")),
			sdk.NewEvent("mint_token", sdk.NewAttribute("amount", sdk.NewInt(1).String()+symbol)),
			sdk.NewEvent("mint_token", sdk.NewAttribute("to", addr1.String())),
			sdk.NewEvent(types.EventTypeModifyTokenURIPermToken, sdk.NewAttribute("to", addr1.String())),
			sdk.NewEvent(types.EventTypeModifyTokenURIPermToken, sdk.NewAttribute("perm_resource", symbol)),
			sdk.NewEvent(types.EventTypeModifyTokenURIPermToken, sdk.NewAttribute("perm_action", types.ModifyAction)),
			sdk.NewEvent("transfer", sdk.NewAttribute("recipient", addr1.String())),
			sdk.NewEvent("transfer", sdk.NewAttribute("amount", sdk.NewInt(1).String()+symbol)),
		}
		verifyEventFunc(t, e, res.Events)
	}

	ctx = ctx.WithEventManager(sdk.NewEventManager())
	{
		symbol := "t03" + suffixAddr1
		msg := types.NewMsgCreateCollection(name, symbol, addr1)
		require.NoError(t, msg.ValidateBasic())
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())

		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "token")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("grant_perm_token", sdk.NewAttribute("to", addr1.String())),
			sdk.NewEvent("grant_perm_token", sdk.NewAttribute("perm_resource", symbol)),
			sdk.NewEvent("grant_perm_token", sdk.NewAttribute("perm_action", "issue")),
			sdk.NewEvent("create_collection", sdk.NewAttribute("name", name)),
			sdk.NewEvent("create_collection", sdk.NewAttribute("symbol", symbol)),
			sdk.NewEvent("create_collection", sdk.NewAttribute("owner", addr1.String())),
		}
		verifyEventFunc(t, e, res.Events)
	}

	ctx = ctx.WithEventManager(sdk.NewEventManager())
	{
		symbol := "t03" + suffixAddr1
		symbolWithID := symbol + "00000001"
		msg := types.NewMsgIssueCollection(name, symbol, tokenuri, addr1, amount, decimals, true, "00000001")
		require.NoError(t, msg.ValidateBasic())
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())

		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("sender", tokenAddr.String())),
			sdk.NewEvent("message", sdk.NewAttribute("module", "token")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("grant_perm_token", sdk.NewAttribute("to", addr1.String())),
			sdk.NewEvent("grant_perm_token", sdk.NewAttribute("perm_resource", symbolWithID)),
			sdk.NewEvent("grant_perm_token", sdk.NewAttribute("perm_action", "mint")),
			sdk.NewEvent("issue_token", sdk.NewAttribute("name", name)),
			sdk.NewEvent("issue_token", sdk.NewAttribute("symbol", symbol)),
			sdk.NewEvent("issue_token", sdk.NewAttribute("denom", symbolWithID)),
			sdk.NewEvent("issue_token", sdk.NewAttribute("owner", addr1.String())),
			sdk.NewEvent("issue_token", sdk.NewAttribute("amount", amount.String())),
			sdk.NewEvent("issue_token", sdk.NewAttribute("mintable", "true")),
			sdk.NewEvent("issue_token", sdk.NewAttribute("decimals", decimals.String())),
			sdk.NewEvent("issue_token", sdk.NewAttribute("token_type", "cft")),
			sdk.NewEvent("mint_token", sdk.NewAttribute("amount", amount.String()+symbolWithID)),
			sdk.NewEvent("mint_token", sdk.NewAttribute("to", addr1.String())),
			sdk.NewEvent(types.EventTypeModifyTokenURIPermToken, sdk.NewAttribute("to", addr1.String())),
			sdk.NewEvent(types.EventTypeModifyTokenURIPermToken, sdk.NewAttribute("perm_resource", symbolWithID)),
			sdk.NewEvent(types.EventTypeModifyTokenURIPermToken, sdk.NewAttribute("perm_action", types.ModifyAction)),
			sdk.NewEvent("transfer", sdk.NewAttribute("recipient", addr1.String())),
			sdk.NewEvent("transfer", sdk.NewAttribute("amount", amount.String()+symbolWithID)),
		}
		verifyEventFunc(t, e, res.Events)
	}

	ctx = ctx.WithEventManager(sdk.NewEventManager())
	{
		symbol := "t03" + suffixAddr1
		symbolWithID := symbol + "00000001"
		msg := types.NewMsgCollectionTokenMint(addr1, addr1, linktype.NewCoinWithTokenIDs(linktype.NewCoinWithTokenID(symbol, "00000001", amount)))
		require.NoError(t, msg.ValidateBasic())
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())
		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("sender", tokenAddr.String())),
			sdk.NewEvent("message", sdk.NewAttribute("module", "token")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("mint_token", sdk.NewAttribute("amount", amount.String()+symbolWithID)),
			sdk.NewEvent("mint_token", sdk.NewAttribute("to", addr1.String())),
			sdk.NewEvent("transfer", sdk.NewAttribute("recipient", addr1.String())),
			sdk.NewEvent("transfer", sdk.NewAttribute("amount", amount.String()+symbolWithID)),
		}
		verifyEventFunc(t, e, res.Events)
	}

	ctx = ctx.WithEventManager(sdk.NewEventManager())
	{
		symbol := "t03" + suffixAddr1
		symbolWithID := symbol + "00000001"
		msg := types.NewMsgCollectionTokenBurn(addr1, linktype.NewCoinWithTokenIDs(linktype.NewCoinWithTokenID(symbol, "00000001", amount)))
		require.NoError(t, msg.ValidateBasic())
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())
		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("message", sdk.NewAttribute("module", "token")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("burn_token", sdk.NewAttribute("amount", amount.String()+symbolWithID)),
			sdk.NewEvent("burn_token", sdk.NewAttribute("from", addr1.String())),
			sdk.NewEvent("transfer", sdk.NewAttribute("recipient", tokenAddr.String())),
			sdk.NewEvent("transfer", sdk.NewAttribute("amount", amount.String()+symbolWithID)),
		}
		verifyEventFunc(t, e, res.Events)
	}

	ctx = ctx.WithEventManager(sdk.NewEventManager())
	{
		symbol := "t03" + suffixAddr1
		symbolWithID := symbol + "a0000002"
		msg := types.NewMsgIssueNFTCollection(name, symbol, tokenuri, addr1, "a0000002")
		require.NoError(t, msg.ValidateBasic())
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())

		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("sender", tokenAddr.String())),
			sdk.NewEvent("message", sdk.NewAttribute("module", "token")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("issue_token", sdk.NewAttribute("name", name)),
			sdk.NewEvent("issue_token", sdk.NewAttribute("symbol", symbol)),
			sdk.NewEvent("issue_token", sdk.NewAttribute("denom", symbolWithID)),
			sdk.NewEvent("issue_token", sdk.NewAttribute("owner", addr1.String())),
			sdk.NewEvent("issue_token", sdk.NewAttribute("token_uri", tokenuri)),
			sdk.NewEvent("issue_token", sdk.NewAttribute("token_type", "cnft")),
			sdk.NewEvent("mint_token", sdk.NewAttribute("amount", sdk.NewInt(1).String()+symbolWithID)),
			sdk.NewEvent("mint_token", sdk.NewAttribute("to", addr1.String())),
			sdk.NewEvent(types.EventTypeModifyTokenURIPermToken, sdk.NewAttribute("to", addr1.String())),
			sdk.NewEvent(types.EventTypeModifyTokenURIPermToken, sdk.NewAttribute("perm_resource", symbolWithID)),
			sdk.NewEvent(types.EventTypeModifyTokenURIPermToken, sdk.NewAttribute("perm_action", types.ModifyAction)),
			sdk.NewEvent("transfer", sdk.NewAttribute("recipient", addr1.String())),
			sdk.NewEvent("transfer", sdk.NewAttribute("amount", sdk.NewInt(1).String()+symbolWithID)),
		}
		verifyEventFunc(t, e, res.Events)
	}

	permission := types.Permission{
		Action:   "issue",
		Resource: "t03" + suffixAddr1,
	}

	ctx = ctx.WithEventManager(sdk.NewEventManager())
	{
		msg := types.NewMsgGrantPermission(addr1, addr2, permission)
		require.NoError(t, msg.ValidateBasic())
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())

		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "token")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("grant_perm_token", sdk.NewAttribute("from", addr1.String())),
			sdk.NewEvent("grant_perm_token", sdk.NewAttribute("to", addr2.String())),
			sdk.NewEvent("grant_perm_token", sdk.NewAttribute("perm_resource", permission.GetResource())),
			sdk.NewEvent("grant_perm_token", sdk.NewAttribute("perm_action", permission.GetAction())),
		}
		verifyEventFunc(t, e, res.Events)
	}
	ctx = ctx.WithEventManager(sdk.NewEventManager())
	{
		msg := types.NewMsgRevokePermission(addr1, permission)
		require.NoError(t, msg.ValidateBasic())
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())
		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "token")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("revoke_perm_token", sdk.NewAttribute("from", addr1.String())),
			sdk.NewEvent("revoke_perm_token", sdk.NewAttribute("perm_resource", permission.GetResource())),
			sdk.NewEvent("revoke_perm_token", sdk.NewAttribute("perm_action", permission.GetAction())),
		}
		verifyEventFunc(t, e, res.Events)
	}
}

func TestHandleTransfer(t *testing.T) {
	input := testCommon.SetupTestInput(t)
	ctx, keeper, ak, bk := input.Ctx, input.Keeper, input.Ak, input.Bk

	h := NewHandler(keeper)

	bk.SetSendEnabled(ctx, true)

	acc := ak.NewAccountWithAddress(ctx, addr1)
	err := acc.SetCoins(sdk.NewCoins(sdk.NewInt64Coin(symbol1, 100)))
	require.NoError(t, err)
	ak.SetAccount(ctx, acc)

	msg := types.NewMsgTransferFT(addr1, addr2, symbol1, sdk.NewInt(10))
	res := h(ctx, msg)
	require.True(t, res.Code.IsOK())
}

func TestHandleTransferError(t *testing.T) {
	input := testCommon.SetupTestInput(t)
	ctx, ak, bk := input.Ctx, input.Ak, input.Bk

	bk.SetSendEnabled(ctx, true)

	acc := ak.NewAccountWithAddress(ctx, addr1)
	err := acc.SetCoins(sdk.NewCoins(sdk.NewInt64Coin("link", 100)))
	require.NoError(t, err)
	ak.SetAccount(ctx, acc)

	msg := types.NewMsgTransferFT(addr1, addr2, "link", sdk.NewInt(10))
	require.EqualError(t, msg.ValidateBasic(), sdk.ErrInvalidCoins("Only user defined token is possible: link").Error())
}

func TestHandleTransferCFT(t *testing.T) {
	input := testCommon.SetupTestInput(t)
	ctx, keeper := input.Ctx, input.Keeper

	h := NewHandler(keeper)

	{
		createMsg := types.NewMsgCreateCollection(name, symbol1, addr1)
		res := h(ctx, createMsg)
		require.True(t, res.Code.IsOK())
		msg := types.NewMsgIssueCollection(name, symbol1, tokenuri, addr1, amount, decimals, true, tokenId1)
		res = h(ctx, msg)
		require.True(t, res.Code.IsOK())
	}

	msg := types.NewMsgTransferCFT(addr1, addr2, symbol1, tokenId1, amount)
	res := h(ctx, msg)
	require.True(t, res.Code.IsOK())
}

func TestHandleTransferNFT(t *testing.T) {
	input := testCommon.SetupTestInput(t)
	ctx, keeper, bk := input.Ctx, input.Keeper, input.Bk

	h := NewHandler(keeper)

	bk.SetSendEnabled(ctx, true)

	{
		msg := types.NewMsgIssueNFT(name, symbol1, tokenuri, addr1)
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())
	}

	msg := types.NewMsgTransferNFT(addr1, addr2, symbol1)
	res := h(ctx, msg)
	require.True(t, res.Code.IsOK())
}

func TestHandleTransferCNFT(t *testing.T) {
	input := testCommon.SetupTestInput(t)
	ctx, keeper, bk := input.Ctx, input.Keeper, input.Bk

	h := NewHandler(keeper)

	bk.SetSendEnabled(ctx, true)

	{
		createMsg := types.NewMsgCreateCollection(name, symbol1, addr1)
		res := h(ctx, createMsg)
		require.True(t, res.Code.IsOK())
		msg := types.NewMsgIssueNFTCollection(name, symbol1, tokenuri, addr1, tokenId1)
		res = h(ctx, msg)
		require.True(t, res.Code.IsOK())
	}

	msg := types.NewMsgTransferCNFT(addr1, addr2, symbol1, tokenId1)
	res := h(ctx, msg)
	require.True(t, res.Code.IsOK())
}

func TestHandleAttachDetach(t *testing.T) {
	input := testCommon.SetupTestInput(t)
	ctx, keeper, bk := input.Ctx, input.Keeper, input.Bk

	h := NewHandler(keeper)

	bk.SetSendEnabled(ctx, true)

	{
		createMsg := types.NewMsgCreateCollection(name, symbol1, addr1)
		res := h(ctx, createMsg)
		require.True(t, res.Code.IsOK())
		msg := types.NewMsgIssueNFTCollection(name, symbol1, tokenuri, addr1, tokenId1)
		res = h(ctx, msg)
		require.True(t, res.Code.IsOK())
		msg = types.NewMsgIssueNFTCollection(name, symbol1, tokenuri, addr1, tokenId2)
		res = h(ctx, msg)
		require.True(t, res.Code.IsOK())
	}

	msg := types.NewMsgAttach(addr1, symbol1, tokenId1, tokenId2)
	res := h(ctx, msg)
	require.True(t, res.Code.IsOK())

	msg2 := types.NewMsgDetach(addr1, addr1, symbol1, tokenId2)
	res2 := h(ctx, msg2)
	require.True(t, res2.Code.IsOK())
}
