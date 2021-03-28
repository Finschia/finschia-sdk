package handler

import (
	"testing"

	"github.com/line/lbm-sdk/v2/x/collection/internal/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
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
	msg := types.NewMsgModify(addr1, defaultContractID, "", "", types.NewChanges(
		types.NewChange("name", modifiedName),
		types.NewChange("base_img_uri", modifiedImgURI),
		types.NewChange("meta", modifiedMeta),
	))

	t.Log("Test with nonexistent token")
	{
		// When handle MsgModify
		_, err := h(ctx, msg)

		// Then response is error
		require.Error(t, err)
	}

	t.Log("Test modify token")
	{
		// Given created collection
		res, err := h(ctx, types.NewMsgCreateCollection(addr1, defaultName, defaultMeta, defaultImgURI))
		require.NoError(t, err)
		contractID = GetMadeContractID(res.Events)

		// When handle MsgModify
		msg = types.NewMsgModify(addr1, contractID, "", "", types.NewChanges(
			types.NewChange("name", modifiedName),
			types.NewChange("base_img_uri", modifiedImgURI),
			types.NewChange("meta", modifiedMeta)))
		res, err = h(ctx, msg)

		// Then response is success
		require.NoError(t, err)
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
	res, err := h(ctx, types.NewMsgCreateCollection(addr1, defaultName, defaultMeta, defaultImgURI))
	require.NoError(t, err)
	contractID := GetMadeContractID(res.Events)

	// Given MsgModify
	msg := types.NewMsgModify(addr1, contractID, defaultTokenType, defaultTokenIndex, types.NewChanges(
		types.NewChange("name", modifiedTokenName),
		types.NewChange("meta", modifiedMeta),
	))

	t.Log("Test with nonexistent token")
	{
		// When handle MsgModify
		_, err := h(ctx, msg)

		// Then response is error
		require.Error(t, err)
	}

	t.Log("Test modify token")
	{
		// Given token
		_, err = h(ctx, types.NewMsgIssueNFT(addr1, contractID, defaultName, defaultMeta))
		require.NoError(t, err)
		param := types.NewMintNFTParam(defaultName, defaultMeta, defaultTokenType)
		_, err = h(ctx, types.NewMsgMintNFT(addr1, contractID, addr1, param))
		require.NoError(t, err)

		// When handle MsgModify
		res, err = h(ctx, msg)

		// Then response is success
		require.NoError(t, err)
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
	res, err := h(ctx, types.NewMsgCreateCollection(addr1, defaultName, defaultMeta, defaultImgURI))
	require.NoError(t, err)
	contractID := GetMadeContractID(res.Events)

	// Given MsgModify
	msg := types.NewMsgModify(addr1, contractID, defaultTokenType, "", types.NewChanges(
		types.NewChange("name", modifiedTokenName),
		types.NewChange("meta", modifiedMeta),
	))

	t.Log("Test with nonexistent token type")
	{
		// When handle MsgModify
		_, err := h(ctx, msg)

		// Then response is error
		require.Error(t, err)
	}

	t.Log("Test modify token type")
	{
		// Given token type
		_, err = h(ctx, types.NewMsgIssueNFT(addr1, contractID, defaultName, defaultMeta))
		require.NoError(t, err)

		// When handle MsgModify
		res, err = h(ctx, msg)

		// Then response is success
		require.NoError(t, err)
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
