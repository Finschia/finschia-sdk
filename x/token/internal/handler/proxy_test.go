package handler

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/lbm-sdk/v2/x/token/internal/types"
	"github.com/stretchr/testify/require"
)

func TestHandleApprove(t *testing.T) {
	ctx, h := cacheKeeper()

	var contractID string
	{
		msgIssue := types.NewMsgIssue(addr1, addr1, defaultName, defaultSymbol, defaultMeta, defaultImageURI, sdk.NewInt(defaultAmount), sdk.NewInt(defaultDecimals), true)
		res, err := h(ctx, msgIssue)
		require.NoError(t, err)

		contractID = GetMadeContractID(res.Events)

		msgMint := types.NewMsgMint(addr1, contractID, addr1, sdk.NewInt(defaultAmount))
		_, err = h(ctx, msgMint)
		require.NoError(t, err)
	}

	msg := types.NewMsgTransferFrom(addr2, contractID, addr1, addr2, sdk.NewInt(defaultAmount))
	_, err := h(ctx, msg)
	require.Error(t, err)

	{
		msgApprove := types.NewMsgApprove(addr1, contractID, addr2)
		_, err := h(ctx, msgApprove)
		require.NoError(t, err)
	}

	msg = types.NewMsgTransferFrom(addr2, contractID, addr1, addr2, sdk.NewInt(defaultAmount))
	res, err := h(ctx, msg)
	require.NoError(t, err)

	e := sdk.Events{
		sdk.NewEvent("message", sdk.NewAttribute("module", "token")),
		sdk.NewEvent("message", sdk.NewAttribute("sender", addr2.String())),
		sdk.NewEvent("transfer_from", sdk.NewAttribute("contract_id", contractID)),
		sdk.NewEvent("transfer_from", sdk.NewAttribute("proxy", addr2.String())),
		sdk.NewEvent("transfer_from", sdk.NewAttribute("from", addr1.String())),
		sdk.NewEvent("transfer_from", sdk.NewAttribute("to", addr2.String())),
		sdk.NewEvent("transfer_from", sdk.NewAttribute("amount", sdk.NewInt(defaultAmount).String())),
	}
	verifyEventFunc(t, e, res.Events)
}
