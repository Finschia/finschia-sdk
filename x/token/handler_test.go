package token

import (
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	linktype "github.com/line/link/types"
	testCommon "github.com/line/link/x/token/internal/keeper"
	"github.com/line/link/x/token/internal/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

const (
	symbolCony = "cony"
	name       = "description"
	tokenuri   = "tokenuri"
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
		msg := types.NewMsgIssueCFT(name, symbol1, tokenuri, addr1, amount, decimals, true)
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())
	}
	{
		msg := types.NewMsgIssueCFT(name, symbol1, tokenuri, addr1, amount, decimals, true)
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())
	}
	{
		msg := types.NewMsgIssueCFT(name, symbol1, tokenuri, addr2, amount, decimals, true)
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
		msg := types.NewMsgIssueCFT(name, symbol1, tokenuri, addr2, amount, decimals, true)
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())
	}
	{
		msg := types.NewMsgRevokePermission(addr1, permission)
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())
	}
	{
		msg := types.NewMsgIssueCFT(name, symbol1, tokenuri, addr1, amount, decimals, true)
		res := h(ctx, msg)
		require.False(t, res.Code.IsOK())
		require.Equal(t, types.DefaultCodespace, res.Codespace)
		require.Equal(t, types.CodeTokenPermission, res.Code)
	}
}

func TestHandlerIssueNFTCollection(t *testing.T) {
	input := testCommon.SetupTestInput(t)
	ctx, keeper := input.Ctx, input.Keeper

	h := NewHandler(keeper)

	{
		msg := types.NewMsgCreateCollection(name, symbol1, addr1)
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())
	}

	{
		//Expect token type is 1001
		{
			msg := types.NewMsgIssueCNFT(symbol1, addr1)
			res := h(ctx, msg)
			require.True(t, res.Code.IsOK())
		}
		//Expect token type is 1002
		{
			msg := types.NewMsgIssueCNFT(symbol1, addr1)
			res := h(ctx, msg)
			require.True(t, res.Code.IsOK())
		}
		{
			msg := types.NewMsgMintCNFT(name, symbol1, tokenuri, "1001", addr1, addr1)
			res := h(ctx, msg)
			require.True(t, res.Code.IsOK())
		}
		{
			msg := types.NewMsgMintCNFT(name, symbol1, tokenuri, "1001", addr1, addr2)
			res := h(ctx, msg)
			require.True(t, res.Code.IsOK())
		}
		{
			mintPermission := types.Permission{
				Action:   "mint",
				Resource: symbol1 + "1001",
			}
			{
				msg := types.NewMsgGrantPermission(addr1, addr2, mintPermission)
				res := h(ctx, msg)
				require.True(t, res.Code.IsOK())
			}
			{
				msg := types.NewMsgMintCNFT(name, symbol1, tokenuri, "1001", addr2, addr2)
				res := h(ctx, msg)
				require.True(t, res.Code.IsOK())
			}
			{
				msg := types.NewMsgRevokePermission(addr1, mintPermission)
				res := h(ctx, msg)
				require.True(t, res.Code.IsOK())
			}
			{
				msg := types.NewMsgMintCNFT(name, symbol1, tokenuri, "1001", addr1, addr1)
				res := h(ctx, msg)
				require.False(t, res.Code.IsOK())
				require.Equal(t, types.DefaultCodespace, res.Codespace)
				require.Equal(t, types.CodeTokenPermission, res.Code)
			}
		}
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

	//Expect token type is 1003
	{
		msg := types.NewMsgIssueCNFT(symbol1, addr2)
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())
	}
	{
		msg := types.NewMsgMintCNFT(name, symbol1, tokenuri, "1003", addr2, addr2)
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())
	}
	{
		msg := types.NewMsgRevokePermission(addr1, permission)
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())
	}
	{
		msg := types.NewMsgIssueCNFT(symbol1, addr1)
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

	{
		symbol := "t01" + suffixAddr1
		msg := types.NewMsgIssue(name, symbol, tokenuri, addr1, amount, decimals, true)
		require.NoError(t, msg.ValidateBasic())
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())

		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "token")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("to", addr1.String())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("perm_resource", symbol)),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("perm_action", types.ModifyAction)),
			sdk.NewEvent("issue", sdk.NewAttribute("name", name)),
			sdk.NewEvent("issue", sdk.NewAttribute("symbol", symbol)),
			sdk.NewEvent("issue", sdk.NewAttribute("token_id", "")),
			sdk.NewEvent("issue", sdk.NewAttribute("owner", addr1.String())),
			sdk.NewEvent("issue", sdk.NewAttribute("amount", amount.String())),
			sdk.NewEvent("issue", sdk.NewAttribute("mintable", "true")),
			sdk.NewEvent("issue", sdk.NewAttribute("decimals", decimals.String())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("to", addr1.String())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("perm_resource", symbol)),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("perm_action", "mint")),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("to", addr1.String())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("perm_resource", symbol)),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("perm_action", "burn")),
		}
		verifyEventFunc(t, e, res.Events)
	}
	{
		symbol := "t01" + suffixAddr1
		msg := types.NewMsgMint(addr1, addr1, sdk.NewCoins(sdk.NewInt64Coin(symbol, amount.Int64())))
		require.NoError(t, msg.ValidateBasic())
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())
		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "token")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("mint", sdk.NewAttribute("amount", amount.String()+symbol)),
			sdk.NewEvent("mint", sdk.NewAttribute("from", addr1.String())),
			sdk.NewEvent("mint", sdk.NewAttribute("to", addr1.String())),
		}
		verifyEventFunc(t, e, res.Events)
	}

	{
		symbol := "t01" + suffixAddr1
		msg := types.NewMsgBurn(addr1, sdk.NewCoins(sdk.NewInt64Coin(symbol, amount.Int64())))
		require.NoError(t, msg.ValidateBasic())
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())
		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "token")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("burn", sdk.NewAttribute("amount", amount.String()+symbol)),
			sdk.NewEvent("burn", sdk.NewAttribute("from", addr1.String())),
		}
		verifyEventFunc(t, e, res.Events)
	}

	{
		symbol := "t03" + suffixAddr1
		msg := types.NewMsgCreateCollection(name, symbol, addr1)
		require.NoError(t, msg.ValidateBasic())
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())

		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "token")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("to", addr1.String())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("perm_resource", symbol)),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("perm_action", "issue")),
			sdk.NewEvent("create_collection", sdk.NewAttribute("name", name)),
			sdk.NewEvent("create_collection", sdk.NewAttribute("symbol", symbol)),
			sdk.NewEvent("create_collection", sdk.NewAttribute("owner", addr1.String())),
		}
		verifyEventFunc(t, e, res.Events)
	}

	{
		symbol := "t03" + suffixAddr1
		symbolWithID := symbol + "00010000"
		msg := types.NewMsgIssueCFT(name, symbol, tokenuri, addr1, amount, decimals, true)
		require.NoError(t, msg.ValidateBasic())
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())

		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "token")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("issue_cft", sdk.NewAttribute("name", name)),
			sdk.NewEvent("issue_cft", sdk.NewAttribute("symbol", symbol)),
			sdk.NewEvent("issue_cft", sdk.NewAttribute("token_id", "00010000")),
			sdk.NewEvent("issue_cft", sdk.NewAttribute("owner", addr1.String())),
			sdk.NewEvent("issue_cft", sdk.NewAttribute("amount", amount.String())),
			sdk.NewEvent("issue_cft", sdk.NewAttribute("mintable", "true")),
			sdk.NewEvent("issue_cft", sdk.NewAttribute("decimals", decimals.String())),
			sdk.NewEvent("issue_cft", sdk.NewAttribute("token_uri", tokenuri)),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("to", addr1.String())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("perm_resource", symbolWithID)),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("perm_action", types.ModifyAction)),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("to", addr1.String())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("perm_resource", symbolWithID)),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("perm_action", "mint")),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("to", addr1.String())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("perm_resource", symbolWithID)),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("perm_action", "burn")),
		}
		verifyEventFunc(t, e, res.Events)
	}

	{
		symbol := "t03" + suffixAddr1
		symbolWithID := symbol + "00010000"
		msg := types.NewMsgMintCFT(addr1, addr1, linktype.NewCoinWithTokenIDs(linktype.NewCoinWithTokenID(symbol, "00010000", amount)))
		require.NoError(t, msg.ValidateBasic())
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())
		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "token")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("mint_cft", sdk.NewAttribute("amount", amount.String()+symbolWithID)),
			sdk.NewEvent("mint_cft", sdk.NewAttribute("from", addr1.String())),
			sdk.NewEvent("mint_cft", sdk.NewAttribute("to", addr1.String())),
		}
		verifyEventFunc(t, e, res.Events)
	}

	{
		symbol := "t03" + suffixAddr1
		symbolWithID := symbol + "00010000"
		msg := types.NewMsgBurnCFT(addr1, linktype.NewCoinWithTokenIDs(linktype.NewCoinWithTokenID(symbol, "00010000", amount)))
		require.NoError(t, msg.ValidateBasic())
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())
		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "token")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("burn_cft", sdk.NewAttribute("amount", amount.String()+symbolWithID)),
			sdk.NewEvent("burn_cft", sdk.NewAttribute("from", addr1.String())),
		}
		verifyEventFunc(t, e, res.Events)
	}

	{
		symbol := "t03" + suffixAddr1
		msg := types.NewMsgIssueCNFT(symbol, addr1)
		require.NoError(t, msg.ValidateBasic())
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())
		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "token")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("issue_cnft", sdk.NewAttribute("symbol", symbol)),
			sdk.NewEvent("issue_cnft", sdk.NewAttribute("token_type", "1001")),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("to", addr1.String())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("perm_resource", symbol+"1001")),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("perm_action", types.MintAction)),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("to", addr1.String())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("perm_resource", symbol+"1001")),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("perm_action", types.BurnAction)),
		}
		verifyEventFunc(t, e, res.Events)
	}

	{
		symbol := "t03" + suffixAddr1
		symbolWithID := symbol + "10010001"
		msg := types.NewMsgMintCNFT(name, symbol, tokenuri, "1001", addr1, addr1)
		require.NoError(t, msg.ValidateBasic())
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())

		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "token")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("mint_cnft", sdk.NewAttribute("name", name)),
			sdk.NewEvent("mint_cnft", sdk.NewAttribute("symbol", symbol)),
			sdk.NewEvent("mint_cnft", sdk.NewAttribute("token_id", "10010001")),
			sdk.NewEvent("mint_cnft", sdk.NewAttribute("from", addr1.String())),
			sdk.NewEvent("mint_cnft", sdk.NewAttribute("to", addr1.String())),
			sdk.NewEvent("mint_cnft", sdk.NewAttribute("token_uri", tokenuri)),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("to", addr1.String())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("perm_resource", symbolWithID)),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("perm_action", types.ModifyAction)),
		}
		verifyEventFunc(t, e, res.Events)
	}

	permission := types.Permission{
		Action:   "issue",
		Resource: "t03" + suffixAddr1,
	}

	{
		msg := types.NewMsgGrantPermission(addr1, addr2, permission)
		require.NoError(t, msg.ValidateBasic())
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())

		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "token")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("from", addr1.String())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("to", addr2.String())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("perm_resource", permission.GetResource())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("perm_action", permission.GetAction())),
		}
		verifyEventFunc(t, e, res.Events)
	}
	{
		msg := types.NewMsgRevokePermission(addr1, permission)
		require.NoError(t, msg.ValidateBasic())
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())
		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "token")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("revoke_perm", sdk.NewAttribute("from", addr1.String())),
			sdk.NewEvent("revoke_perm", sdk.NewAttribute("perm_resource", permission.GetResource())),
			sdk.NewEvent("revoke_perm", sdk.NewAttribute("perm_action", permission.GetAction())),
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
	e := sdk.Events{
		sdk.NewEvent("message", sdk.NewAttribute("module", "token")),
		sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
		sdk.NewEvent("transfer", sdk.NewAttribute("from", addr1.String())),
		sdk.NewEvent("transfer", sdk.NewAttribute("to", addr2.String())),
		sdk.NewEvent("transfer", sdk.NewAttribute("symbol", symbol1)),
		sdk.NewEvent("transfer", sdk.NewAttribute("amount", sdk.NewInt(10).String())),
	}
	verifyEventFunc(t, e, res.Events)
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
		msg := types.NewMsgIssueCFT(name, symbol1, tokenuri, addr1, amount, decimals, true)
		res = h(ctx, msg)
		require.True(t, res.Code.IsOK())
	}

	msg := types.NewMsgTransferCFT(addr1, addr2, symbol1, "00010000", amount)
	res := h(ctx, msg)
	require.True(t, res.Code.IsOK())
	e := sdk.Events{
		sdk.NewEvent("message", sdk.NewAttribute("module", "token")),
		sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
		sdk.NewEvent("transfer_cft", sdk.NewAttribute("from", addr1.String())),
		sdk.NewEvent("transfer_cft", sdk.NewAttribute("to", addr2.String())),
		sdk.NewEvent("transfer_cft", sdk.NewAttribute("symbol", symbol1)),
		sdk.NewEvent("transfer_cft", sdk.NewAttribute("token_id", "00010000")),
		sdk.NewEvent("transfer_cft", sdk.NewAttribute("amount", amount.String())),
	}
	verifyEventFunc(t, e, res.Events)
}

func TestHandleTransferCFTFrom(t *testing.T) {
	input := testCommon.SetupTestInput(t)
	ctx, keeper := input.Ctx, input.Keeper

	h := NewHandler(keeper)

	{
		createMsg := types.NewMsgCreateCollection(name, symbol1, addr1)
		res := h(ctx, createMsg)
		require.True(t, res.Code.IsOK())
		msg := types.NewMsgIssueCFT(name, symbol1, tokenuri, addr1, amount, decimals, true)
		res = h(ctx, msg)
		require.True(t, res.Code.IsOK())
		msg2 := types.NewMsgApproveCollection(addr1, addr2, symbol1)
		res = h(ctx, msg2)
		require.True(t, res.Code.IsOK())
	}

	msg := types.NewMsgTransferCFTFrom(addr2, addr1, addr2, symbol1, "00010000", amount)
	res := h(ctx, msg)
	require.True(t, res.Code.IsOK())
	e := sdk.Events{
		sdk.NewEvent("message", sdk.NewAttribute("module", "token")),
		sdk.NewEvent("message", sdk.NewAttribute("sender", addr2.String())),
		sdk.NewEvent("transfer_cft_from", sdk.NewAttribute("proxy", addr2.String())),
		sdk.NewEvent("transfer_cft_from", sdk.NewAttribute("from", addr1.String())),
		sdk.NewEvent("transfer_cft_from", sdk.NewAttribute("to", addr2.String())),
		sdk.NewEvent("transfer_cft_from", sdk.NewAttribute("symbol", symbol1)),
		sdk.NewEvent("transfer_cft_from", sdk.NewAttribute("token_id", "00010000")),
		sdk.NewEvent("transfer_cft_from", sdk.NewAttribute("amount", amount.String())),
	}
	verifyEventFunc(t, e, res.Events)
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
		msg := types.NewMsgIssueCNFT(symbol1, addr1)
		res = h(ctx, msg)
		require.True(t, res.Code.IsOK())
		msg2 := types.NewMsgMintCNFT(name, symbol1, tokenuri, "1001", addr1, addr1)
		res = h(ctx, msg2)
		require.True(t, res.Code.IsOK())
	}

	msg := types.NewMsgTransferCNFT(addr1, addr2, symbol1, "10010001")
	res := h(ctx, msg)
	require.True(t, res.Code.IsOK())
	e := sdk.Events{
		sdk.NewEvent("message", sdk.NewAttribute("module", "token")),
		sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
		sdk.NewEvent("transfer_cnft", sdk.NewAttribute("from", addr1.String())),
		sdk.NewEvent("transfer_cnft", sdk.NewAttribute("to", addr2.String())),
		sdk.NewEvent("transfer_cnft", sdk.NewAttribute("symbol", symbol1)),
		sdk.NewEvent("transfer_cnft", sdk.NewAttribute("token_id", "10010001")),
	}
	verifyEventFunc(t, e, res.Events)
}

func TestHandleTransferCNFTFrom(t *testing.T) {
	input := testCommon.SetupTestInput(t)
	ctx, keeper, bk := input.Ctx, input.Keeper, input.Bk

	h := NewHandler(keeper)

	bk.SetSendEnabled(ctx, true)

	{
		createMsg := types.NewMsgCreateCollection(name, symbol1, addr1)
		res := h(ctx, createMsg)
		require.True(t, res.Code.IsOK())
		msg := types.NewMsgIssueCNFT(symbol1, addr1)
		res = h(ctx, msg)
		require.True(t, res.Code.IsOK())
		msg2 := types.NewMsgMintCNFT(name, symbol1, tokenuri, "1001", addr1, addr1)
		res = h(ctx, msg2)
		require.True(t, res.Code.IsOK())
		msg3 := types.NewMsgApproveCollection(addr1, addr2, symbol1)
		res = h(ctx, msg3)
		require.True(t, res.Code.IsOK())
	}

	msg := types.NewMsgTransferCNFTFrom(addr2, addr1, addr2, symbol1, "10010001")
	res := h(ctx, msg)
	require.True(t, res.Code.IsOK())
	e := sdk.Events{
		sdk.NewEvent("message", sdk.NewAttribute("module", "token")),
		sdk.NewEvent("message", sdk.NewAttribute("sender", addr2.String())),
		sdk.NewEvent("transfer_cnft_from", sdk.NewAttribute("proxy", addr2.String())),
		sdk.NewEvent("transfer_cnft_from", sdk.NewAttribute("from", addr1.String())),
		sdk.NewEvent("transfer_cnft_from", sdk.NewAttribute("to", addr2.String())),
		sdk.NewEvent("transfer_cnft_from", sdk.NewAttribute("symbol", symbol1)),
		sdk.NewEvent("transfer_cnft_from", sdk.NewAttribute("token_id", "10010001")),
	}
	verifyEventFunc(t, e, res.Events)
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
		msg := types.NewMsgIssueCNFT(symbol1, addr1)
		res = h(ctx, msg)
		require.True(t, res.Code.IsOK())
		msg2 := types.NewMsgMintCNFT(name, symbol1, tokenuri, "1001", addr1, addr1)
		res = h(ctx, msg2)
		require.True(t, res.Code.IsOK())
		msg2 = types.NewMsgMintCNFT(name, symbol1, tokenuri, "1001", addr1, addr1)
		res = h(ctx, msg2)
		require.True(t, res.Code.IsOK())
	}

	{
		msg := types.NewMsgAttach(addr1, symbol1, "10010001", "10010002")
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())
		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "token")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("attach", sdk.NewAttribute("from", addr1.String())),
			sdk.NewEvent("attach", sdk.NewAttribute("symbol", symbol1)),
			sdk.NewEvent("attach", sdk.NewAttribute("to_token_id", "10010001")),
			sdk.NewEvent("attach", sdk.NewAttribute("token_id", "10010002")),
		}
		verifyEventFunc(t, e, res.Events)
	}

	{
		msg2 := types.NewMsgDetach(addr1, addr1, symbol1, "10010002")
		res2 := h(ctx, msg2)
		require.True(t, res2.Code.IsOK())
		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "token")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("detach", sdk.NewAttribute("from", addr1.String())),
			sdk.NewEvent("detach", sdk.NewAttribute("to", addr1.String())),
			sdk.NewEvent("detach", sdk.NewAttribute("symbol", symbol1)),
			sdk.NewEvent("detach", sdk.NewAttribute("token_id", "10010002")),
		}
		verifyEventFunc(t, e, res2.Events)
	}

	//Attach again
	{
		msg := types.NewMsgAttach(addr1, symbol1, "10010001", "10010002")
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())
		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "token")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("attach", sdk.NewAttribute("from", addr1.String())),
			sdk.NewEvent("attach", sdk.NewAttribute("symbol", symbol1)),
			sdk.NewEvent("attach", sdk.NewAttribute("to_token_id", "10010001")),
			sdk.NewEvent("attach", sdk.NewAttribute("token_id", "10010002")),
		}
		verifyEventFunc(t, e, res.Events)
	}
	//Burn token
	{
		msg := types.NewMsgBurnCNFT(symbol1, "10010001", addr1)
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())
		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "token")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("burn_cnft", sdk.NewAttribute("symbol", symbol1)),
			sdk.NewEvent("burn_cnft", sdk.NewAttribute("token_id", "10010001")),
			sdk.NewEvent("burn_cnft", sdk.NewAttribute("from", addr1.String())),
		}
		verifyEventFunc(t, e, res.Events)
	}
}

func TestHandleAttachFromDetachFrom(t *testing.T) {
	input := testCommon.SetupTestInput(t)
	ctx, keeper := input.Ctx, input.Keeper

	h := NewHandler(keeper)

	{
		createMsg := types.NewMsgCreateCollection(name, symbol1, addr1)
		res := h(ctx, createMsg)
		require.True(t, res.Code.IsOK())
		msg := types.NewMsgIssueCNFT(symbol1, addr1)
		res = h(ctx, msg)
		require.True(t, res.Code.IsOK())
		msg2 := types.NewMsgMintCNFT(name, symbol1, tokenuri, "1001", addr1, addr1)
		res = h(ctx, msg2)
		require.True(t, res.Code.IsOK())
		msg2 = types.NewMsgMintCNFT(name, symbol1, tokenuri, "1001", addr1, addr1)
		res = h(ctx, msg2)
		require.True(t, res.Code.IsOK())
		msg3 := types.NewMsgApproveCollection(addr1, addr2, symbol1)
		res = h(ctx, msg3)
		require.True(t, res.Code.IsOK())
	}

	msg := types.NewMsgAttachFrom(addr2, addr1, symbol1, "10010001", "10010002")
	res := h(ctx, msg)
	require.True(t, res.Code.IsOK())
	e := sdk.Events{
		sdk.NewEvent("message", sdk.NewAttribute("module", "token")),
		sdk.NewEvent("message", sdk.NewAttribute("sender", addr2.String())),
		sdk.NewEvent("attach_from", sdk.NewAttribute("proxy", addr2.String())),
		sdk.NewEvent("attach_from", sdk.NewAttribute("from", addr1.String())),
		sdk.NewEvent("attach_from", sdk.NewAttribute("symbol", symbol1)),
		sdk.NewEvent("attach_from", sdk.NewAttribute("to_token_id", "10010001")),
		sdk.NewEvent("attach_from", sdk.NewAttribute("token_id", "10010002")),
	}
	verifyEventFunc(t, e, res.Events)

	msg2 := types.NewMsgDetachFrom(addr2, addr1, addr1, symbol1, "10010002")
	res2 := h(ctx, msg2)
	require.True(t, res2.Code.IsOK())
	e = sdk.Events{
		sdk.NewEvent("message", sdk.NewAttribute("module", "token")),
		sdk.NewEvent("message", sdk.NewAttribute("sender", addr2.String())),
		sdk.NewEvent("detach_from", sdk.NewAttribute("proxy", addr2.String())),
		sdk.NewEvent("detach_from", sdk.NewAttribute("from", addr1.String())),
		sdk.NewEvent("detach_from", sdk.NewAttribute("to", addr1.String())),
		sdk.NewEvent("detach_from", sdk.NewAttribute("symbol", symbol1)),
		sdk.NewEvent("detach_from", sdk.NewAttribute("token_id", "10010002")),
	}
	verifyEventFunc(t, e, res2.Events)
}

func TestHandleApproveDisapprove(t *testing.T) {
	input := testCommon.SetupTestInput(t)
	ctx, keeper, bk := input.Ctx, input.Keeper, input.Bk

	h := NewHandler(keeper)

	bk.SetSendEnabled(ctx, true)

	{
		createMsg := types.NewMsgCreateCollection(name, symbol1, addr1)
		res := h(ctx, createMsg)
		require.True(t, res.Code.IsOK())
		msg := types.NewMsgIssueCNFT(symbol1, addr1)
		res = h(ctx, msg)
		require.True(t, res.Code.IsOK())
		msg2 := types.NewMsgMintCNFT(name, symbol1, tokenuri, "1001", addr1, addr1)
		res = h(ctx, msg2)
		require.True(t, res.Code.IsOK())
	}

	msg := types.NewMsgTransferCNFTFrom(addr2, addr1, addr2, symbol1, "10010001")
	res := h(ctx, msg)
	require.False(t, res.Code.IsOK())

	{
		msg3 := types.NewMsgApproveCollection(addr1, addr2, symbol1)
		res = h(ctx, msg3)
		require.True(t, res.Code.IsOK())
	}

	msg = types.NewMsgTransferCNFTFrom(addr2, addr1, addr2, symbol1, "10010001")
	res = h(ctx, msg)
	require.True(t, res.Code.IsOK())

	e := sdk.Events{
		sdk.NewEvent("message", sdk.NewAttribute("module", "token")),
		sdk.NewEvent("message", sdk.NewAttribute("sender", addr2.String())),
		sdk.NewEvent("transfer_cnft_from", sdk.NewAttribute("proxy", addr2.String())),
		sdk.NewEvent("transfer_cnft_from", sdk.NewAttribute("from", addr1.String())),
		sdk.NewEvent("transfer_cnft_from", sdk.NewAttribute("to", addr2.String())),
		sdk.NewEvent("transfer_cnft_from", sdk.NewAttribute("symbol", symbol1)),
		sdk.NewEvent("transfer_cnft_from", sdk.NewAttribute("token_id", "10010001")),
	}
	verifyEventFunc(t, e, res.Events)

	{
		msg3 := types.NewMsgDisapproveCollection(addr1, addr2, symbol1)
		res = h(ctx, msg3)
		require.True(t, res.Code.IsOK())
	}

	msg = types.NewMsgTransferCNFTFrom(addr2, addr1, addr2, symbol1, "10010001")
	res = h(ctx, msg)
	require.False(t, res.Code.IsOK())
}
