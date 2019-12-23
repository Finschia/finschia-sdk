package token

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	description = "description"
	tokenuri    = "tokenuri"
)

var (
	amount      = sdk.NewInt(1000)
	decimals    = sdk.NewInt(6)
	priv1       = secp256k1.GenPrivKey()
	addr1       = sdk.AccAddress(priv1.PubKey().Address())
	suffixAddr1 = symbolCony + addr1.String()[len(addr1.String())-3:]
	priv2       = secp256k1.GenPrivKey()
	addr2       = sdk.AccAddress(priv2.PubKey().Address())
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

func TestHandlerIssueFT(t *testing.T) {
	input := testCommon.SetupTestInput(t)
	ctx, keeper := input.Ctx, input.Keeper

	h := NewHandler(keeper)

	{
		msg := types.NewMsgIssue(description, suffixAddr1, tokenuri, addr1, amount, decimals, true)
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())
	}
	{
		msg := types.NewMsgIssue(description, suffixAddr1, tokenuri, addr1, amount, decimals, true)
		res := h(ctx, msg)
		require.False(t, res.Code.IsOK())
	}
}

func TestHandlerIssueNFT(t *testing.T) {
	input := testCommon.SetupTestInput(t)
	ctx, keeper := input.Ctx, input.Keeper

	h := NewHandler(keeper)

	{
		msg := types.NewMsgIssueNFT(description, suffixAddr1, tokenuri, addr1)
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())
	}
	{
		msg := types.NewMsgIssueNFT(description, suffixAddr1, tokenuri, addr1)
		res := h(ctx, msg)
		require.False(t, res.Code.IsOK())
	}
}

func TestHandlerIssueFTCollection(t *testing.T) {
	input := testCommon.SetupTestInput(t)
	ctx, keeper := input.Ctx, input.Keeper

	h := NewHandler(keeper)

	{
		msg := types.NewMsgIssueCollection(description, suffixAddr1, tokenuri, addr1, amount, decimals, true, "00000001")
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())
	}
	{
		msg := types.NewMsgIssueCollection(description, suffixAddr1, tokenuri, addr1, amount, decimals, true, "00000001")
		res := h(ctx, msg)
		require.False(t, res.Code.IsOK())
	}
	{
		msg := types.NewMsgIssueCollection(description, suffixAddr1, tokenuri, addr1, amount, decimals, true, "00000002")
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())
	}
	{
		msg := types.NewMsgIssueCollection(description, suffixAddr1, tokenuri, addr2, amount, decimals, true, "00000003")
		res := h(ctx, msg)
		require.False(t, res.Code.IsOK())
	}

	permission := types.Permission{
		Action:   "issue",
		Resource: suffixAddr1,
	}

	{
		msg := types.NewMsgGrantPermission(addr1, addr2, permission)
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())
	}
	{
		msg := types.NewMsgIssueCollection(description, suffixAddr1, tokenuri, addr2, amount, decimals, true, "00000003")
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())
	}
	{
		msg := types.NewMsgRevokePermission(addr1, permission)
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())
	}
	{
		msg := types.NewMsgIssueCollection(description, suffixAddr1, tokenuri, addr1, amount, decimals, true, "00000004")
		res := h(ctx, msg)
		require.False(t, res.Code.IsOK())
	}

}

type attr struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type event struct {
	Type       string `json:"type"`
	Attributes []attr `json:"attributes"`
}

func getEvents(rawEvents sdk.Events) []event {
	var events []event
	for _, re := range rawEvents {

		var attrs []attr
		for _, a := range re.Attributes {
			attrs = append(attrs, attr{
				Key:   string(a.GetKey()),
				Value: string(a.GetValue()),
			},
			)
		}

		events = append(events, event{
			Type:       re.Type,
			Attributes: attrs,
		})
	}
	return events
}

func mustMarshalJson(v interface{}) string {
	bz, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		panic(err)
	}
	return string(bz)
}

func TestEvents(t *testing.T) {
	input := testCommon.SetupTestInput(t)

	ctx, keeper := input.Ctx, input.Keeper

	h := NewHandler(keeper)

	suffixAddr1 = addr1.String()[len(addr1.String())-3:]

	tokenAddr := sdk.AccAddress(crypto.AddressHash([]byte("token")))

	verifyEventFunc := func(t *testing.T, expected sdk.Events, actual sdk.Events) {
		require.Equal(t, sdk.StringifyEvents(expected.ToABCIEvents()).String(), sdk.StringifyEvents(actual.ToABCIEvents()).String())
	}

	ctx = ctx.WithEventManager(sdk.NewEventManager())
	{
		symbol := "t01" + suffixAddr1
		msg := types.NewMsgIssue(description, symbol, tokenuri, addr1, amount, decimals, true)
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())

		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("sender", tokenAddr.String())),
			sdk.NewEvent("message", sdk.NewAttribute("module", "token")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("grant_perm_token", sdk.NewAttribute("to", addr1.String())),
			sdk.NewEvent("grant_perm_token", sdk.NewAttribute("perm_resource", symbol)),
			sdk.NewEvent("grant_perm_token", sdk.NewAttribute("perm_action", "issue")),
			sdk.NewEvent("grant_perm_token", sdk.NewAttribute("to", addr1.String())),
			sdk.NewEvent("grant_perm_token", sdk.NewAttribute("perm_resource", symbol)),
			sdk.NewEvent("grant_perm_token", sdk.NewAttribute("perm_action", "mint")),
			sdk.NewEvent("issue_token", sdk.NewAttribute("name", description)),
			sdk.NewEvent("issue_token", sdk.NewAttribute("symbol", symbol)),
			sdk.NewEvent("issue_token", sdk.NewAttribute("owner", addr1.String())),
			sdk.NewEvent("issue_token", sdk.NewAttribute("amount", amount.String())),
			sdk.NewEvent("issue_token", sdk.NewAttribute("mintable", "true")),
			sdk.NewEvent("issue_token", sdk.NewAttribute("decimals", decimals.String())),
			sdk.NewEvent("issue_token", sdk.NewAttribute("token_uri", "")),
			sdk.NewEvent("issue_token", sdk.NewAttribute("token_type", "ft")),
			sdk.NewEvent("mint_token", sdk.NewAttribute("amount", amount.String()+symbol)),
			sdk.NewEvent("mint_token", sdk.NewAttribute("to", addr1.String())),
			sdk.NewEvent("occupy_symbol", sdk.NewAttribute("symbol", symbol)),
			sdk.NewEvent("occupy_symbol", sdk.NewAttribute("owner", addr1.String())),
			sdk.NewEvent("transfer", sdk.NewAttribute("recipient", addr1.String())),
			sdk.NewEvent("transfer", sdk.NewAttribute("amount", amount.String()+symbol)),
		}
		verifyEventFunc(t, e, res.Events)
	}

	ctx = ctx.WithEventManager(sdk.NewEventManager())
	{
		symbol := "t02" + suffixAddr1
		msg := types.NewMsgIssueNFT(description, symbol, tokenuri, addr1)
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())

		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("sender", tokenAddr.String())),
			sdk.NewEvent("message", sdk.NewAttribute("module", "token")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("grant_perm_token", sdk.NewAttribute("to", addr1.String())),
			sdk.NewEvent("grant_perm_token", sdk.NewAttribute("perm_resource", symbol)),
			sdk.NewEvent("grant_perm_token", sdk.NewAttribute("perm_action", "issue")),
			sdk.NewEvent("issue_token", sdk.NewAttribute("name", description)),
			sdk.NewEvent("issue_token", sdk.NewAttribute("symbol", symbol)),
			sdk.NewEvent("issue_token", sdk.NewAttribute("owner", addr1.String())),
			sdk.NewEvent("issue_token", sdk.NewAttribute("amount", sdk.NewInt(1).String())),
			sdk.NewEvent("issue_token", sdk.NewAttribute("mintable", "false")),
			sdk.NewEvent("issue_token", sdk.NewAttribute("decimals", sdk.NewInt(0).String())),
			sdk.NewEvent("issue_token", sdk.NewAttribute("token_uri", tokenuri)),
			sdk.NewEvent("issue_token", sdk.NewAttribute("token_type", "nft")),
			sdk.NewEvent("mint_token", sdk.NewAttribute("amount", sdk.NewInt(1).String()+symbol)),
			sdk.NewEvent("mint_token", sdk.NewAttribute("to", addr1.String())),
			sdk.NewEvent("occupy_symbol", sdk.NewAttribute("symbol", symbol)),
			sdk.NewEvent("occupy_symbol", sdk.NewAttribute("owner", addr1.String())),
			sdk.NewEvent("transfer", sdk.NewAttribute("recipient", addr1.String())),
			sdk.NewEvent("transfer", sdk.NewAttribute("amount", sdk.NewInt(1).String()+symbol)),
		}
		verifyEventFunc(t, e, res.Events)
	}
	ctx = ctx.WithEventManager(sdk.NewEventManager())
	{
		symbol := "t03" + suffixAddr1
		symbolWithID := symbol + "00000001"
		msg := types.NewMsgIssueCollection(description, symbol, tokenuri, addr1, amount, decimals, true, "00000001")
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())

		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("sender", tokenAddr.String())),
			sdk.NewEvent("message", sdk.NewAttribute("module", "token")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("grant_perm_token", sdk.NewAttribute("to", addr1.String())),
			sdk.NewEvent("grant_perm_token", sdk.NewAttribute("perm_resource", symbol)),
			sdk.NewEvent("grant_perm_token", sdk.NewAttribute("perm_action", "issue")),
			sdk.NewEvent("grant_perm_token", sdk.NewAttribute("to", addr1.String())),
			sdk.NewEvent("grant_perm_token", sdk.NewAttribute("perm_resource", symbolWithID)),
			sdk.NewEvent("grant_perm_token", sdk.NewAttribute("perm_action", "mint")),
			sdk.NewEvent("issue_token", sdk.NewAttribute("name", description)),
			sdk.NewEvent("issue_token", sdk.NewAttribute("symbol", symbolWithID)),
			sdk.NewEvent("issue_token", sdk.NewAttribute("owner", addr1.String())),
			sdk.NewEvent("issue_token", sdk.NewAttribute("amount", amount.String())),
			sdk.NewEvent("issue_token", sdk.NewAttribute("mintable", "true")),
			sdk.NewEvent("issue_token", sdk.NewAttribute("decimals", decimals.String())),
			sdk.NewEvent("issue_token", sdk.NewAttribute("token_uri", "")),
			sdk.NewEvent("issue_token", sdk.NewAttribute("token_type", "cft")),
			sdk.NewEvent("mint_token", sdk.NewAttribute("amount", amount.String()+symbolWithID)),
			sdk.NewEvent("mint_token", sdk.NewAttribute("to", addr1.String())),
			sdk.NewEvent("occupy_symbol", sdk.NewAttribute("symbol", symbol)),
			sdk.NewEvent("occupy_symbol", sdk.NewAttribute("owner", addr1.String())),
			sdk.NewEvent("transfer", sdk.NewAttribute("recipient", addr1.String())),
			sdk.NewEvent("transfer", sdk.NewAttribute("amount", amount.String()+symbolWithID)),
		}
		verifyEventFunc(t, e, res.Events)
	}
	ctx = ctx.WithEventManager(sdk.NewEventManager())
	{
		symbol := "t04" + suffixAddr1
		symbolWithID := symbol + "00000001"
		msg := types.NewMsgIssueNFTCollection(description, symbol, tokenuri, addr1, "00000001")
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())

		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("sender", tokenAddr.String())),
			sdk.NewEvent("message", sdk.NewAttribute("module", "token")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("grant_perm_token", sdk.NewAttribute("to", addr1.String())),
			sdk.NewEvent("grant_perm_token", sdk.NewAttribute("perm_resource", symbol)),
			sdk.NewEvent("grant_perm_token", sdk.NewAttribute("perm_action", "issue")),
			sdk.NewEvent("issue_token", sdk.NewAttribute("name", description)),
			sdk.NewEvent("issue_token", sdk.NewAttribute("symbol", symbolWithID)),
			sdk.NewEvent("issue_token", sdk.NewAttribute("owner", addr1.String())),
			sdk.NewEvent("issue_token", sdk.NewAttribute("amount", sdk.NewInt(1).String())),
			sdk.NewEvent("issue_token", sdk.NewAttribute("mintable", "false")),
			sdk.NewEvent("issue_token", sdk.NewAttribute("decimals", sdk.NewInt(0).String())),
			sdk.NewEvent("issue_token", sdk.NewAttribute("token_uri", tokenuri)),
			sdk.NewEvent("issue_token", sdk.NewAttribute("token_type", "cnft")),
			sdk.NewEvent("mint_token", sdk.NewAttribute("amount", sdk.NewInt(1).String()+symbolWithID)),
			sdk.NewEvent("mint_token", sdk.NewAttribute("to", addr1.String())),
			sdk.NewEvent("occupy_symbol", sdk.NewAttribute("symbol", symbol)),
			sdk.NewEvent("occupy_symbol", sdk.NewAttribute("owner", addr1.String())),
			sdk.NewEvent("transfer", sdk.NewAttribute("recipient", addr1.String())),
			sdk.NewEvent("transfer", sdk.NewAttribute("amount", sdk.NewInt(1).String()+symbolWithID)),
		}
		verifyEventFunc(t, e, res.Events)
	}

	permission := types.Permission{
		Action:   "issue",
		Resource: "t01" + suffixAddr1,
	}

	ctx = ctx.WithEventManager(sdk.NewEventManager())
	{
		msg := types.NewMsgGrantPermission(addr1, addr2, permission)
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
