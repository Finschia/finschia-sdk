package handler

import (
	"testing"

	"github.com/line/link/x/token/internal/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	linktype "github.com/line/link/types"
)

func TestHandleMsgModify(t *testing.T) {
	ctx, h := cacheKeeper()
	const (
		modifiedTokenName = "modifiedTokenName"
		modifiedTokenURI  = "modifiedTokenURI"
	)
	// Given MsgModify
	msg := types.NewMsgModify(addr1, defaultSymbol, linktype.NewChanges(
		linktype.NewChange("name", modifiedTokenName),
		linktype.NewChange("token_uri", modifiedTokenURI),
	))

	t.Log("Test with nonexistent token")
	{
		// When handle MsgModify
		res := h(ctx, msg)

		// Then response is error
		require.False(t, res.Code.IsOK())
		require.Equal(t, types.DefaultCodespace, res.Codespace)
		require.Equal(t, types.CodeTokenNotExist, res.Code)
		verifyEventFunc(t, nil, res.Events)
	}

	t.Log("Test modify token")
	{
		// Given issued token
		res := h(ctx, types.NewMsgIssue(addr1, defaultName, defaultSymbol, defaultTokenURI,
			sdk.NewInt(defaultAmount), sdk.NewInt(defaultDecimals), true))
		require.True(t, res.IsOK())

		// When handle MsgModify
		res = h(ctx, msg)

		// Then response is success
		require.True(t, res.Code.IsOK())
		// And events are returned
		expectedEvents := sdk.Events{
			sdk.NewEvent(types.EventTypeModifyToken, sdk.NewAttribute(types.AttributeKeyModifiedField, "name")),
			sdk.NewEvent(types.EventTypeModifyToken, sdk.NewAttribute(types.AttributeKeyName, modifiedTokenName)),
			sdk.NewEvent(types.EventTypeModifyToken, sdk.NewAttribute(types.AttributeKeySymbol, defaultSymbol)),
			sdk.NewEvent(types.EventTypeModifyToken, sdk.NewAttribute(types.AttributeKeyOwner, addr1.String())),
			sdk.NewEvent(types.EventTypeModifyToken, sdk.NewAttribute(types.AttributeKeyTokenURI, defaultTokenURI)),
			sdk.NewEvent(types.EventTypeModifyToken, sdk.NewAttribute(types.AttributeKeyModifiedField, "token_uri")),
			sdk.NewEvent(types.EventTypeModifyToken, sdk.NewAttribute(types.AttributeKeyName, modifiedTokenName)),
			sdk.NewEvent(types.EventTypeModifyToken, sdk.NewAttribute(types.AttributeKeySymbol, defaultSymbol)),
			sdk.NewEvent(types.EventTypeModifyToken, sdk.NewAttribute(types.AttributeKeyOwner, addr1.String())),
			sdk.NewEvent(types.EventTypeModifyToken, sdk.NewAttribute(types.AttributeKeyTokenURI, modifiedTokenURI)),
			sdk.NewEvent(sdk.EventTypeMessage, sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory)),
			sdk.NewEvent(sdk.EventTypeMessage, sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String())),
		}
		verifyEventFunc(t, expectedEvents, res.Events)
	}
}
