package handler

import (
	"testing"

	"github.com/line/link/x/contract"
	"github.com/line/link/x/token/internal/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	linktype "github.com/line/link/types"
)

func TestHandleMsgModify(t *testing.T) {
	ctx, h := cacheKeeper()

	contractID := contract.SampleContractID
	const (
		modifiedTokenName = "modifiedTokenName"
		modifiedTokenURI  = "modifiedTokenURI"
	)
	// Given MsgModify
	msg := types.NewMsgModify(addr1, contractID, linktype.NewChanges(
		linktype.NewChange("name", modifiedTokenName),
		linktype.NewChange("token_uri", modifiedTokenURI),
	))

	t.Log("Test with nonexistent token")
	{
		// When handle MsgModify
		res := h(ctx, msg)

		// Then response is error
		require.False(t, res.Code.IsOK())
		require.Equal(t, contract.ContractCodeSpace, res.Codespace)
		require.Equal(t, contract.CodeContractNotExist, res.Code)
		verifyEventFunc(t, nil, res.Events)
	}

	t.Log("Test modify token")
	{
		// Given issued token
		res := h(ctx, types.NewMsgIssue(addr1, addr1, defaultName, defaultContractID, defaultImageURI,
			sdk.NewInt(defaultAmount), sdk.NewInt(defaultDecimals), true))
		require.True(t, res.IsOK())
		contractID := GetMadeContractID(res.Events)

		msg := types.NewMsgModify(addr1, contractID, linktype.NewChanges(
			linktype.NewChange("name", modifiedTokenName),
			linktype.NewChange("token_uri", modifiedTokenURI),
		))

		// When handle MsgModify
		res = h(ctx, msg)

		// Then response is success
		require.True(t, res.Code.IsOK())
		// And events are returned
		expectedEvents := sdk.Events{
			sdk.NewEvent(types.EventTypeModifyToken, sdk.NewAttribute(types.AttributeKeyContractID, defaultContractID)),
			sdk.NewEvent(types.EventTypeModifyToken, sdk.NewAttribute("name", modifiedTokenName)),
			sdk.NewEvent(types.EventTypeModifyToken, sdk.NewAttribute("token_uri", modifiedTokenURI)),
			sdk.NewEvent(sdk.EventTypeMessage, sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory)),
			sdk.NewEvent(sdk.EventTypeMessage, sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String())),
		}
		verifyEventFunc(t, expectedEvents, res.Events)
	}
}
