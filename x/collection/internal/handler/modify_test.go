package handler

import (
	"testing"

	"github.com/line/link/x/collection/internal/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	linktype "github.com/line/link/types"
)

func TestHandleMsgModifyForCollection(t *testing.T) {
	ctx, h := cacheKeeper()
	const (
		modifiedName   = "modifiedName"
		modifiedImgURI = "modifiedImgURI"
	)

	// Given MsgModify
	msg := types.NewMsgModify(addr1, defaultSymbol, "", "", linktype.NewChanges(
		linktype.NewChange("name", modifiedName),
		linktype.NewChange("base_img_uri", modifiedImgURI),
	))

	t.Log("Test with nonexistent token")
	{
		// When handle MsgModify
		res := h(ctx, msg)

		// Then response is error
		require.False(t, res.Code.IsOK())
		require.Equal(t, types.DefaultCodespace, res.Codespace)
		require.Equal(t, types.CodeCollectionNotExist, res.Code)
		verifyEventFunc(t, nil, res.Events)
	}

	t.Log("Test modify token")
	{
		// Given created collection
		res := h(ctx, types.NewMsgCreateCollection(addr1, defaultName, defaultSymbol, defaultImgURI))
		require.True(t, res.IsOK())

		// When handle MsgModify
		res = h(ctx, msg)

		// Then response is success
		require.True(t, res.Code.IsOK())
		// And events are returned
		expectedEvents := sdk.Events{
			sdk.NewEvent(types.EventTypeModifyCollection, sdk.NewAttribute(types.AttributeKeyModifiedField, "name")),
			sdk.NewEvent(types.EventTypeModifyCollection, sdk.NewAttribute(types.AttributeKeyName, modifiedName)),
			sdk.NewEvent(types.EventTypeModifyCollection, sdk.NewAttribute(types.AttributeKeySymbol, defaultSymbol)),
			sdk.NewEvent(types.EventTypeModifyCollection, sdk.NewAttribute(types.AttributeKeyOwner, addr1.String())),
			sdk.NewEvent(types.EventTypeModifyCollection, sdk.NewAttribute(types.AttributeKeyBaseImgURI, defaultImgURI)),
			sdk.NewEvent(types.EventTypeModifyCollection, sdk.NewAttribute(types.AttributeKeyModifiedField, "base_img_uri")),
			sdk.NewEvent(types.EventTypeModifyCollection, sdk.NewAttribute(types.AttributeKeyName, modifiedName)),
			sdk.NewEvent(types.EventTypeModifyCollection, sdk.NewAttribute(types.AttributeKeySymbol, defaultSymbol)),
			sdk.NewEvent(types.EventTypeModifyCollection, sdk.NewAttribute(types.AttributeKeyOwner, addr1.String())),
			sdk.NewEvent(types.EventTypeModifyCollection, sdk.NewAttribute(types.AttributeKeyBaseImgURI, modifiedImgURI)),
			sdk.NewEvent(sdk.EventTypeMessage, sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory)),
			sdk.NewEvent(sdk.EventTypeMessage, sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String())),
		}
		verifyEventFunc(t, expectedEvents, res.Events)
	}
}

func TestHandleMsgModifyForToken(t *testing.T) {
	ctx, h := cacheKeeper()
	const (
		modifiedTokenName = "modifiedTokenName"
	)
	// Given MsgModify
	msg := types.NewMsgModify(addr1, defaultSymbol, defaultTokenType, defaultTokenIndex, linktype.NewChanges(
		linktype.NewChange("name", modifiedTokenName),
	))

	// And created collection
	res := h(ctx, types.NewMsgCreateCollection(addr1, defaultName, defaultSymbol, defaultImgURI))
	require.True(t, res.IsOK())

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
		// Given token
		res = h(ctx, types.NewMsgIssueNFT(addr1, defaultSymbol, defaultName))
		require.True(t, res.IsOK())
		res = h(ctx, types.NewMsgMintNFT(addr1, addr1, defaultName, defaultSymbol, defaultTokenType))
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
			sdk.NewEvent(types.EventTypeModifyToken, sdk.NewAttribute(types.AttributeKeyTokenID, defaultTokenID1)),
			sdk.NewEvent(types.EventTypeModifyToken, sdk.NewAttribute(types.AttributeKeyOwner, addr1.String())),
			sdk.NewEvent(sdk.EventTypeMessage, sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory)),
			sdk.NewEvent(sdk.EventTypeMessage, sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String())),
		}
		verifyEventFunc(t, expectedEvents, res.Events)
	}
}

func TestHandleMsgModifyForTokenType(t *testing.T) {
	ctx, h := cacheKeeper()
	const (
		modifiedTokenName = "modifiedTokenName"
	)

	// Given MsgModify
	msg := types.NewMsgModify(addr1, defaultSymbol, defaultTokenType, "", linktype.NewChanges(
		linktype.NewChange("name", modifiedTokenName),
	))
	// And created collection
	res := h(ctx, types.NewMsgCreateCollection(addr1, defaultName, defaultSymbol, defaultImgURI))
	require.True(t, res.IsOK())

	t.Log("Test with nonexistent token type")
	{
		// When handle MsgModify
		res := h(ctx, msg)

		// Then response is error
		require.False(t, res.Code.IsOK())
		require.Equal(t, types.DefaultCodespace, res.Codespace)
		require.Equal(t, types.CodeCollectionTokenTypeNotExist, res.Code)
		verifyEventFunc(t, nil, res.Events)
	}

	t.Log("Test modify token type")
	{
		// Given token type
		res = h(ctx, types.NewMsgIssueNFT(addr1, defaultSymbol, defaultName))
		require.True(t, res.IsOK())

		// When handle MsgModify
		res = h(ctx, msg)

		// Then response is success
		require.True(t, res.Code.IsOK())
		// And events are returned
		expectedEvents := sdk.Events{
			sdk.NewEvent(types.EventTypeModifyTokenType, sdk.NewAttribute(types.AttributeKeyModifiedField, "name")),
			sdk.NewEvent(types.EventTypeModifyTokenType, sdk.NewAttribute(types.AttributeKeyName, modifiedTokenName)),
			sdk.NewEvent(types.EventTypeModifyTokenType, sdk.NewAttribute(types.AttributeKeySymbol, defaultSymbol)),
			sdk.NewEvent(types.EventTypeModifyTokenType, sdk.NewAttribute(types.AttributeKeyTokenType, defaultTokenType)),
			sdk.NewEvent(types.EventTypeModifyTokenType, sdk.NewAttribute(types.AttributeKeyOwner, addr1.String())),
			sdk.NewEvent(sdk.EventTypeMessage, sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory)),
			sdk.NewEvent(sdk.EventTypeMessage, sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String())),
		}
		verifyEventFunc(t, expectedEvents, res.Events)
	}
}
