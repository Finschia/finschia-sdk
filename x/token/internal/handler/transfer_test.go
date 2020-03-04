package handler

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link/x/contract"

	"github.com/line/link/x/token/internal/types"
	"github.com/stretchr/testify/require"
)

func TestHandleMsgTransfer(t *testing.T) {
	ctx, h := cacheKeeper()

	t.Log("Prepare Token Issued")
	{
		k.NewContractID(ctx)
		token := types.NewToken(defaultContractID, defaultName, defaultSymbol, defaultImageURI, sdk.NewInt(defaultDecimals), true)
		err := k.IssueToken(ctx, token, sdk.NewInt(defaultAmount), addr1)
		require.NoError(t, err)
	}

	t.Log("Transfer Token")
	{
		msg := types.NewMsgTransfer(addr1, addr2, defaultContractID, sdk.NewInt(10))
		res := h(ctx, msg)
		require.True(t, res.Code.IsOK())
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
		require.EqualError(t, msg.ValidateBasic(), contract.ErrInvalidContractID(contract.ContractCodeSpace, defaultCoin).Error())
		res := h(ctx, msg)
		require.False(t, res.Code.IsOK())
	}
}
