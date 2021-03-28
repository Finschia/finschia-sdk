package handler

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/line/lbm-sdk/v2/x/contract"

	"github.com/line/lbm-sdk/v2/x/token/internal/types"
	"github.com/stretchr/testify/require"
)

func TestHandleMsgTransfer(t *testing.T) {
	ctx, h := cacheKeeper()

	t.Log("Prepare Token Issued")
	{
		k.NewContractID(ctx)
		token := types.NewToken(defaultContractID, defaultName, defaultSymbol, defaultMeta, defaultImageURI, sdk.NewInt(defaultDecimals), true)
		err := k.IssueToken(ctx, token, sdk.NewInt(defaultAmount), addr1, addr1)
		require.NoError(t, err)
	}

	t.Log("Transfer Token")
	{
		msg := types.NewMsgTransfer(addr1, addr2, defaultContractID, sdk.NewInt(10))
		res, err := h(ctx, msg)
		require.NoError(t, err)
		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "token")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("transfer", sdk.NewAttribute("from", addr1.String())),
			sdk.NewEvent("transfer", sdk.NewAttribute("to", addr2.String())),
			sdk.NewEvent("transfer", sdk.NewAttribute("contract_id", defaultContractID)),
			sdk.NewEvent("transfer", sdk.NewAttribute("amount", sdk.NewInt(10).String())),
		}
		verifyEventFunc(t, e, res.Events)
	}
	t.Log("Transfer Coin. Expect Fail")
	{
		msg := types.NewMsgTransfer(addr1, addr2, defaultCoin, sdk.NewInt(10))
		require.EqualError(t, msg.ValidateBasic(), sdkerrors.Wrapf(contract.ErrInvalidContractID, "ContractID: %s", defaultCoin).Error())
		_, err := h(ctx, msg)
		require.Error(t, err)
	}
}

func TestHandleTransferFrom(t *testing.T) {
	ctx, h := cacheKeeper()

	var contractID string
	{
		msgIssue := types.NewMsgIssue(addr1, addr1, defaultName, defaultSymbol, defaultMeta, defaultImageURI, sdk.NewInt(defaultAmount), sdk.NewInt(defaultDecimals), true)
		res, err := h(ctx, msgIssue)
		require.NoError(t, err)

		contractID = GetMadeContractID(res.Events)

		msgApprove := types.NewMsgApprove(addr1, contractID, addr2)
		_, err = h(ctx, msgApprove)
		require.NoError(t, err)
	}

	msg := types.NewMsgTransferFrom(addr2, contractID, addr1, addr2, sdk.NewInt(defaultAmount))
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
