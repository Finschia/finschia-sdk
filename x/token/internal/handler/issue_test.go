package handler

import (
	"testing"

	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/token/internal/types"
	"github.com/stretchr/testify/require"
)

func GetMadeContractID(events sdk.Events) string {
	for _, event := range events.ToABCIEvents() {
		for _, attr := range event.Attributes {
			if string(attr.Key) == types.AttributeKeyContractID {
				return string(attr.Value)
			}
		}
	}
	return ""
}

func TestHandleMsgIssue(t *testing.T) {
	ctx, h := cacheKeeper()

	t.Log("Issue Token")
	{
		msg := types.NewMsgIssue(addr1, addr1, defaultName, defaultSymbol, defaultMeta, defaultImageURI, sdk.NewInt(defaultAmount), sdk.NewInt(defaultDecimals), true)
		require.NoError(t, msg.ValidateBasic())
		res, err := h(ctx, msg)
		require.NoError(t, err)
		contractID := GetMadeContractID(res.Events)

		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "token")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("to", addr1.String())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("contract_id", contractID)),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("perm", types.ModifyAction)),
			sdk.NewEvent("issue", sdk.NewAttribute("contract_id", contractID)),
			sdk.NewEvent("issue", sdk.NewAttribute("name", defaultName)),
			sdk.NewEvent("issue", sdk.NewAttribute("symbol", defaultSymbol)),
			sdk.NewEvent("issue", sdk.NewAttribute("owner", addr1.String())),
			sdk.NewEvent("issue", sdk.NewAttribute("to", addr1.String())),
			sdk.NewEvent("issue", sdk.NewAttribute("amount", sdk.NewInt(defaultAmount).String())),
			sdk.NewEvent("issue", sdk.NewAttribute("mintable", "true")),
			sdk.NewEvent("issue", sdk.NewAttribute("decimals", sdk.NewInt(defaultDecimals).String())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("to", addr1.String())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("contract_id", contractID)),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("perm", "mint")),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("to", addr1.String())),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("contract_id", contractID)),
			sdk.NewEvent("grant_perm", sdk.NewAttribute("perm", "burn")),
		}
		verifyEventFunc(t, e, res.Events)
	}

	t.Log("Issue Token Again Expect Success")
	{
		msg := types.NewMsgIssue(addr1, addr1, defaultName, defaultSymbol, defaultMeta, defaultImageURI, sdk.NewInt(defaultAmount), sdk.NewInt(defaultDecimals), true)
		_, err := h(ctx, msg)
		require.NoError(t, err)
	}
}
