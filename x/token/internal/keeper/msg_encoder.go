package keeper

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/line/link-modules/x/token/internal/types"
)

func NewMsgEncodeHandler(tokenKeeper Keeper) types.EncodeHandler {
	return func(jsonMsg json.RawMessage) ([]sdk.Msg, error) {
		var wasmCustomMsg types.WasmCustomMsg
		err := json.Unmarshal(jsonMsg, &wasmCustomMsg)
		if err != nil {
			return nil, err
		}
		switch types.MsgRoute(wasmCustomMsg.Route) {
		case types.RIssue:
			return handleMsgIssue(wasmCustomMsg.Data)
		case types.RTransfer:
			return handleMsgTransfer(wasmCustomMsg.Data)
		case types.RMint:
			return handleMsgMint(wasmCustomMsg.Data)
		case types.RBurn:
			return handleMsgBurn(wasmCustomMsg.Data)
		case types.RGrantPerm:
			return handleMsgGrantPerm(wasmCustomMsg.Data)
		case types.RRevokePerm:
			return handleMsgRevokePerm(wasmCustomMsg.Data)
		case types.RModify:
			return handleMsgModify(wasmCustomMsg.Data)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized Msg route: %T", wasmCustomMsg.Route)
		}
	}
}

func handleMsgIssue(msgData json.RawMessage) ([]sdk.Msg, error) {
	var wrapper types.IssueMsgWrapper
	err := json.Unmarshal(msgData, &wrapper)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	return []sdk.Msg{wrapper.MsgIssue}, nil
}

func handleMsgTransfer(msgData json.RawMessage) ([]sdk.Msg, error) {
	var wrapper types.TransferMsgWrapper
	err := json.Unmarshal(msgData, &wrapper)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	return []sdk.Msg{wrapper.MsgTransfer}, nil
}

func handleMsgMint(msgData json.RawMessage) ([]sdk.Msg, error) {
	var wrapper types.MintMsgWrapper
	err := json.Unmarshal(msgData, &wrapper)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	return []sdk.Msg{wrapper.MsgMint}, nil
}

func handleMsgBurn(msgData json.RawMessage) ([]sdk.Msg, error) {
	var wrapper types.BurnMsgWrapper
	err := json.Unmarshal(msgData, &wrapper)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	return []sdk.Msg{wrapper.MsgBurn}, nil
}

func handleMsgGrantPerm(msgData json.RawMessage) ([]sdk.Msg, error) {
	var wrapper types.GrantPermMsgWrapper
	err := json.Unmarshal(msgData, &wrapper)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	return []sdk.Msg{wrapper.MsgGrantPermission}, nil
}

func handleMsgRevokePerm(msgData json.RawMessage) ([]sdk.Msg, error) {
	var wrapper types.RevokePermMsgWrapper
	err := json.Unmarshal(msgData, &wrapper)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	return []sdk.Msg{wrapper.MsgRevokePermission}, nil
}

func handleMsgModify(msgData json.RawMessage) ([]sdk.Msg, error) {
	var wrapper types.ModifyMsgWrapper
	err := json.Unmarshal(msgData, &wrapper)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	return []sdk.Msg{wrapper.MsgModify}, nil
}
