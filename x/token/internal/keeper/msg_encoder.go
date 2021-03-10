package keeper

import (
	"encoding/json"

	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/token/internal/types"
	"github.com/line/lbm-sdk/x/wasm"
)

func NewMsgEncodeHandler(tokenKeeper Keeper) wasm.EncodeHandler {
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
		case types.RTransferFrom:
			return handleMsgTransferFrom(wasmCustomMsg.Data)
		case types.RMint:
			return handleMsgMint(wasmCustomMsg.Data)
		case types.RBurn:
			return handleMsgBurn(wasmCustomMsg.Data)
		case types.RBurnFrom:
			return handleMsgBurnFrom(wasmCustomMsg.Data)
		case types.RGrantPerm:
			return handleMsgGrantPerm(wasmCustomMsg.Data)
		case types.RRevokePerm:
			return handleMsgRevokePerm(wasmCustomMsg.Data)
		case types.RModify:
			return handleMsgModify(wasmCustomMsg.Data)
		case types.RApprove:
			return handleMsgApprove(wasmCustomMsg.Data)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized Msg route: %T", wasmCustomMsg.Route)
		}
	}
}

func handleMsgIssue(msgData json.RawMessage) ([]sdk.Msg, error) {
	var msg types.MsgIssue
	err := json.Unmarshal(msgData, &msg)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	return []sdk.Msg{msg}, nil
}

func handleMsgTransfer(msgData json.RawMessage) ([]sdk.Msg, error) {
	var msg types.MsgTransfer
	err := json.Unmarshal(msgData, &msg)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	return []sdk.Msg{msg}, nil
}

func handleMsgTransferFrom(msgData json.RawMessage) ([]sdk.Msg, error) {
	var msg types.MsgTransferFrom
	err := json.Unmarshal(msgData, &msg)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	return []sdk.Msg{msg}, nil
}

func handleMsgMint(msgData json.RawMessage) ([]sdk.Msg, error) {
	var msg types.MsgMint
	err := json.Unmarshal(msgData, &msg)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	return []sdk.Msg{msg}, nil
}

func handleMsgBurn(msgData json.RawMessage) ([]sdk.Msg, error) {
	var msg types.MsgBurn
	err := json.Unmarshal(msgData, &msg)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	return []sdk.Msg{msg}, nil
}

func handleMsgBurnFrom(msgData json.RawMessage) ([]sdk.Msg, error) {
	var msg types.MsgBurnFrom
	err := json.Unmarshal(msgData, &msg)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	return []sdk.Msg{msg}, nil
}

func handleMsgGrantPerm(msgData json.RawMessage) ([]sdk.Msg, error) {
	var msg types.MsgGrantPermission
	err := json.Unmarshal(msgData, &msg)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	return []sdk.Msg{msg}, nil
}

func handleMsgRevokePerm(msgData json.RawMessage) ([]sdk.Msg, error) {
	var msg types.MsgRevokePermission
	err := json.Unmarshal(msgData, &msg)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	return []sdk.Msg{msg}, nil
}

func handleMsgModify(msgData json.RawMessage) ([]sdk.Msg, error) {
	var msg types.MsgModify
	err := json.Unmarshal(msgData, &msg)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	return []sdk.Msg{msg}, nil
}

func handleMsgApprove(msgData json.RawMessage) ([]sdk.Msg, error) {
	var msg types.MsgApprove
	err := json.Unmarshal(msgData, &msg)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	return []sdk.Msg{msg}, nil
}
