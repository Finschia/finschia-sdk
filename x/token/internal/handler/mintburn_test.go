package handler

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/lbm-sdk/v2/x/token/internal/types"
	"github.com/stretchr/testify/require"
)

func TestHandleMsgMint(t *testing.T) {
	ctx, h := cacheKeeper()

	t.Log("Prepare Token Issued")
	{
		k.NewContractID(ctx)
		token := types.NewToken(defaultContractID, defaultName, defaultSymbol, defaultMeta, defaultImageURI, sdk.NewInt(defaultDecimals), true)
		err := k.IssueToken(ctx, token, sdk.NewInt(defaultAmount), addr1, addr1)
		require.NoError(t, err)
	}

	t.Log("Burn Tokens")
	{
		msg := types.NewMsgMint(addr1, defaultContractID, addr1, sdk.NewInt(defaultAmount))
		require.NoError(t, msg.ValidateBasic())
		res, err := h(ctx, msg)
		require.NoError(t, err)
		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "token")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("mint", sdk.NewAttribute("contract_id", defaultContractID)),
			sdk.NewEvent("mint", sdk.NewAttribute("amount", sdk.NewInt(defaultAmount).String())),
			sdk.NewEvent("mint", sdk.NewAttribute("from", addr1.String())),
			sdk.NewEvent("mint", sdk.NewAttribute("to", addr1.String())),
		}
		verifyEventFunc(t, e, res.Events)
	}
}

func TestHandleMsgBurn(t *testing.T) {
	ctx, h := cacheKeeper()

	t.Log("Prepare Token Issued")
	{
		k.NewContractID(ctx)
		token := types.NewToken(defaultContractID, defaultName, defaultSymbol, defaultMeta, defaultImageURI, sdk.NewInt(defaultDecimals), true)
		err := k.IssueToken(ctx, token, sdk.NewInt(defaultAmount), addr1, addr1)
		require.NoError(t, err)
	}

	t.Log("Mint Tokens")
	{
		msg := types.NewMsgBurn(addr1, defaultContractID, sdk.NewInt(defaultAmount))
		require.NoError(t, msg.ValidateBasic())
		res, err := h(ctx, msg)
		require.NoError(t, err)
		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "token")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr1.String())),
			sdk.NewEvent("burn", sdk.NewAttribute("contract_id", defaultContractID)),
			sdk.NewEvent("burn", sdk.NewAttribute("amount", sdk.NewInt(defaultAmount).String())),
			sdk.NewEvent("burn", sdk.NewAttribute("from", addr1.String())),
		}
		verifyEventFunc(t, e, res.Events)
	}
}

func TestHandleMsgBurnFTFrom(t *testing.T) {
	ctx, h := cacheKeeper()

	t.Log("Prepare Token Issued")
	{
		k.NewContractID(ctx)
		token := types.NewToken(defaultContractID, defaultName, defaultSymbol, defaultMeta, defaultImageURI, sdk.NewInt(defaultDecimals), true)
		err := k.IssueToken(ctx, token, sdk.NewInt(defaultAmount), addr1, addr1)
		require.NoError(t, err)
	}

	t.Log("fail to burn when addr2 is not approved ")
	{
		burnMsg := types.NewMsgBurnFrom(addr2, defaultContractID, addr1, sdk.NewInt(100))
		_, err := h(ctx, burnMsg)
		require.Error(t, err)
	}

	t.Log("Approve addr2")
	{
		msgApprove := types.NewMsgApprove(addr1, defaultContractID, addr2)
		_, err := h(ctx, msgApprove)
		require.NoError(t, err)
	}

	t.Log("give permission to addr2")
	{
		permission := types.NewBurnPermission()
		msg := types.NewMsgGrantPermission(addr1, defaultContractID, addr2, permission)
		require.NoError(t, msg.ValidateBasic())
		_, err := h(ctx, msg)
		require.NoError(t, err)
	}

	t.Log("Error when invalid user")
	{
		burnMsg := types.NewMsgBurnFrom(addr2, defaultContractID, addr2, sdk.NewInt(100))
		_, err := h(ctx, burnMsg)
		require.Error(t, err)
	}

	t.Log("fail to burn over the being supplied")
	{
		burnMsg := types.NewMsgBurnFrom(addr2, defaultContractID, addr1, sdk.NewInt(1001))
		_, err := h(ctx, burnMsg)
		require.Error(t, err)
	}

	t.Log("fail to burn when invalid ContractID")
	{
		burnMsg := types.NewMsgBurnFrom(addr2, "abcd11234", addr1, sdk.NewInt(1000))
		_, err := h(ctx, burnMsg)
		require.Error(t, err)
	}

	t.Log("Succeed to burn")
	{
		burnMsg := types.NewMsgBurnFrom(addr2, defaultContractID, addr1, sdk.NewInt(100))
		res, err := h(ctx, burnMsg)
		require.NoError(t, err)

		e := sdk.Events{
			sdk.NewEvent("message", sdk.NewAttribute("module", "token")),
			sdk.NewEvent("message", sdk.NewAttribute("sender", addr2.String())),
			sdk.NewEvent("burn_from", sdk.NewAttribute("contract_id", defaultContractID)),
			sdk.NewEvent("burn_from", sdk.NewAttribute("proxy", addr2.String())),
			sdk.NewEvent("burn_from", sdk.NewAttribute("from", addr1.String())),
			sdk.NewEvent("burn_from", sdk.NewAttribute("amount", sdk.NewInt(100).String())),
		}
		verifyEventFunc(t, e, res.Events)
	}
}
