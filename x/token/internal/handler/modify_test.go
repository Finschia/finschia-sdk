package handler

import (
	"testing"

	"github.com/line/link-modules/x/contract"
	"github.com/line/link-modules/x/token/internal/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestHandleMsgModify(t *testing.T) {
	ctx, h := cacheKeeper()

	contractID := contract.SampleContractID
	const (
		modifiedTokenName = "modifiedTokenName"
		modifiedImgURI    = "modifiedImgURI"
		modifiedMeta      = "modifiedMeta"
	)
	// Given MsgModify
	msg := types.NewMsgModify(addr1, contractID, types.NewChanges(
		types.NewChange("name", modifiedTokenName),
		types.NewChange("img_uri", modifiedImgURI),
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
		// Given issued token
		res, err := h(ctx, types.NewMsgIssue(addr1, addr1, defaultName, defaultContractID, defaultMeta, defaultImageURI,
			sdk.NewInt(defaultAmount), sdk.NewInt(defaultDecimals), true))
		require.NoError(t, err)
		contractID := GetMadeContractID(res.Events)

		msg := types.NewMsgModify(addr1, contractID, types.NewChanges(
			types.NewChange("name", modifiedTokenName),
			types.NewChange("img_uri", modifiedImgURI),
			types.NewChange("meta", modifiedMeta),
		))

		// When handle MsgModify
		res, err = h(ctx, msg)

		// Then response is success
		require.NoError(t, err)
		// And events are returned
		expectedEvents := sdk.Events{
			sdk.NewEvent(types.EventTypeModifyToken, sdk.NewAttribute(types.AttributeKeyContractID, defaultContractID)),
			sdk.NewEvent(types.EventTypeModifyToken, sdk.NewAttribute("name", modifiedTokenName)),
			sdk.NewEvent(types.EventTypeModifyToken, sdk.NewAttribute("img_uri", modifiedImgURI)),
			sdk.NewEvent(types.EventTypeModifyToken, sdk.NewAttribute("meta", modifiedMeta)),
			sdk.NewEvent(sdk.EventTypeMessage, sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory)),
			sdk.NewEvent(sdk.EventTypeMessage, sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String())),
		}
		verifyEventFunc(t, expectedEvents, res.Events)
	}
}
