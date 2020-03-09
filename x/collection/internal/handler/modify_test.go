package handler

import (
	"testing"

	"github.com/line/link/x/collection/internal/types"
	"github.com/line/link/x/contract"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	linktype "github.com/line/link/types"
)

func TestHandleMsgModifyForCollection(t *testing.T) {
	ctx, h := cacheKeeper()
	const (
		modifiedName   = "modifiedName"
		modifiedImgURI = "modifiedImgURI"
		modifiedMeta   = "modifiedMeta"
	)

	var contractID string

	// Given MsgModify
	msg := types.NewMsgModify(addr1, defaultContractID, "", "", linktype.NewChanges(
		linktype.NewChange("name", modifiedName),
		linktype.NewChange("base_img_uri", modifiedImgURI),
		linktype.NewChange("meta", modifiedMeta),
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
		// Given created collection
		res := h(ctx, types.NewMsgCreateCollection(addr1, defaultName, defaultMeta, defaultImgURI))
		require.True(t, res.IsOK())
		contractID = GetMadeContractID(res.Events)

		// When handle MsgModify
		msg = types.NewMsgModify(addr1, contractID, "", "", linktype.NewChanges(
			linktype.NewChange("name", modifiedName),
			linktype.NewChange("base_img_uri", modifiedImgURI),
			linktype.NewChange("meta", modifiedMeta)))
		res = h(ctx, msg)

		// Then response is success
		require.True(t, res.Code.IsOK())
		contractID = GetMadeContractID(res.Events)

		// And events are returned
		expectedEvents := sdk.Events{
			sdk.NewEvent(types.EventTypeModifyCollection, sdk.NewAttribute(types.AttributeKeyContractID, contractID)),
			sdk.NewEvent(types.EventTypeModifyCollection, sdk.NewAttribute("name", modifiedName)),
			sdk.NewEvent(types.EventTypeModifyCollection, sdk.NewAttribute("base_img_uri", modifiedImgURI)),
			sdk.NewEvent(types.EventTypeModifyCollection, sdk.NewAttribute("meta", modifiedMeta)),
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
		modifiedMeta      = "modifiedMeta"
	)

	// created collection
	res := h(ctx, types.NewMsgCreateCollection(addr1, defaultName, defaultMeta, defaultImgURI))
	require.True(t, res.IsOK())
	contractID := GetMadeContractID(res.Events)

	// Given MsgModify
	msg := types.NewMsgModify(addr1, contractID, defaultTokenType, defaultTokenIndex, linktype.NewChanges(
		linktype.NewChange("name", modifiedTokenName),
		linktype.NewChange("meta", modifiedMeta),
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
		// Given token
		res = h(ctx, types.NewMsgIssueNFT(addr1, contractID, defaultName, defaultMeta))
		require.True(t, res.IsOK())
		res = h(ctx, types.NewMsgMintNFT(addr1, contractID, addr1, defaultName, defaultMeta, defaultTokenType))
		require.True(t, res.IsOK())

		// When handle MsgModify
		res = h(ctx, msg)

		// Then response is success
		require.True(t, res.Code.IsOK())
		// And events are returned
		expectedEvents := sdk.Events{
			sdk.NewEvent(types.EventTypeModifyToken, sdk.NewAttribute(types.AttributeKeyContractID, contractID)),
			sdk.NewEvent(types.EventTypeModifyToken, sdk.NewAttribute(types.AttributeKeyTokenID, defaultTokenID1)),
			sdk.NewEvent(types.EventTypeModifyToken, sdk.NewAttribute("name", modifiedTokenName)),
			sdk.NewEvent(types.EventTypeModifyToken, sdk.NewAttribute("meta", modifiedMeta)),
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
		modifiedMeta      = "modifiedMeta"
	)

	// created collection
	res := h(ctx, types.NewMsgCreateCollection(addr1, defaultName, defaultMeta, defaultImgURI))
	require.True(t, res.IsOK())
	contractID := GetMadeContractID(res.Events)

	// Given MsgModify
	msg := types.NewMsgModify(addr1, contractID, defaultTokenType, "", linktype.NewChanges(
		linktype.NewChange("name", modifiedTokenName),
		linktype.NewChange("meta", modifiedMeta),
	))

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
		res = h(ctx, types.NewMsgIssueNFT(addr1, contractID, defaultName, defaultMeta))
		require.True(t, res.IsOK())

		// When handle MsgModify
		res = h(ctx, msg)

		// Then response is success
		require.True(t, res.Code.IsOK())
		// And events are returned
		expectedEvents := sdk.Events{
			sdk.NewEvent(types.EventTypeModifyTokenType, sdk.NewAttribute(types.AttributeKeyContractID, contractID)),
			sdk.NewEvent(types.EventTypeModifyTokenType, sdk.NewAttribute(types.AttributeKeyTokenType, defaultTokenType)),
			sdk.NewEvent(types.EventTypeModifyTokenType, sdk.NewAttribute("name", modifiedTokenName)),
			sdk.NewEvent(types.EventTypeModifyTokenType, sdk.NewAttribute("meta", modifiedMeta)),
			sdk.NewEvent(sdk.EventTypeMessage, sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory)),
			sdk.NewEvent(sdk.EventTypeMessage, sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String())),
		}
		verifyEventFunc(t, expectedEvents, res.Events)
	}
}
