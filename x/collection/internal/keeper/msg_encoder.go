package keeper

// import (
// 	"encoding/json"

// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
// 	"github.com/line/lbm-sdk/v2/x/collection/internal/types"
// 	"github.com/line/lbm-sdk/v2/x/wasm"
// )

// func NewMsgEncodeHandler(collectionKeeper Keeper) wasm.EncodeHandler {
// 	return func(jsonMsg json.RawMessage) ([]sdk.Msg, error) {
// 		var wasmCustomMsg types.WasmCustomMsg
// 		err := json.Unmarshal(jsonMsg, &wasmCustomMsg)
// 		if err != nil {
// 			return nil, err
// 		}

// 		switch types.MsgRoute(wasmCustomMsg.Route) {
// 		case types.RCreateCollection:
// 			return handleMsgCreateCollection(wasmCustomMsg.Data)
// 		case types.RIssueNFT:
// 			return handleMsgIssueNFT(wasmCustomMsg.Data)
// 		case types.RIssueFT:
// 			return handleMsgIssueFT(wasmCustomMsg.Data)
// 		case types.RMintNFT:
// 			return handleMsgMintNFT(wasmCustomMsg.Data)
// 		case types.RMintFT:
// 			return handleMsgMintFT(wasmCustomMsg.Data)
// 		case types.RBurnNFT:
// 			return handleMsgBurnNFT(wasmCustomMsg.Data)
// 		case types.RBurnNFTFrom:
// 			return handleMsgBurnNFTFrom(wasmCustomMsg.Data)
// 		case types.RBurnFT:
// 			return handleMsgBurnFT(wasmCustomMsg.Data)
// 		case types.RBurnFTFrom:
// 			return handleMsgBurnFTFrom(wasmCustomMsg.Data)
// 		case types.RTransferNFT:
// 			return handleMsgTransferNFT(wasmCustomMsg.Data)
// 		case types.RTransferNFTFrom:
// 			return handleMsgTransferNFTFrom(wasmCustomMsg.Data)
// 		case types.RTransferFT:
// 			return handleMsgTransferFT(wasmCustomMsg.Data)
// 		case types.RTransferFTFrom:
// 			return handleMsgTransferFTFrom(wasmCustomMsg.Data)
// 		case types.RModify:
// 			return handleMsgModify(wasmCustomMsg.Data)
// 		case types.RApprove:
// 			return handleMsgApprove(wasmCustomMsg.Data)
// 		case types.RDisapprove:
// 			return handleMsgDisapprove(wasmCustomMsg.Data)
// 		case types.RGrantPerm:
// 			return handleMsgGrantPerm(wasmCustomMsg.Data)
// 		case types.RRevokePerm:
// 			return handleMsgRevokePerm(wasmCustomMsg.Data)
// 		case types.RAttach:
// 			return handleMsgAttach(wasmCustomMsg.Data)
// 		case types.RDetach:
// 			return handleMsgDetach(wasmCustomMsg.Data)
// 		case types.RAttachFrom:
// 			return handleMsgAttachFrom(wasmCustomMsg.Data)
// 		case types.RDetachFrom:
// 			return handleMsgDetachFrom(wasmCustomMsg.Data)
// 		default:
// 			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized Msg route: %T", wasmCustomMsg.Route)
// 		}
// 	}
// }

// func handleMsgCreateCollection(msgData json.RawMessage) ([]sdk.Msg, error) {
// 	var msg types.MsgCreateCollection
// 	err := json.Unmarshal(msgData, &msg)
// 	if err != nil {
// 		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
// 	}
// 	return []sdk.Msg{msg}, nil
// }

// func handleMsgIssueNFT(msgData json.RawMessage) ([]sdk.Msg, error) {
// 	var msg types.MsgIssueNFT
// 	err := json.Unmarshal(msgData, &msg)
// 	if err != nil {
// 		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
// 	}

// 	return []sdk.Msg{msg}, nil
// }

// func handleMsgIssueFT(msgData json.RawMessage) ([]sdk.Msg, error) {
// 	var msg types.MsgIssueFT
// 	err := json.Unmarshal(msgData, &msg)
// 	if err != nil {
// 		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
// 	}
// 	if msg.Decimals.Int64() < 0 || msg.Decimals.Int64() > 18 {
// 		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid decimals. 0 <= decimals <= 18")
// 	}

// 	return []sdk.Msg{msg}, nil
// }

// func handleMsgMintNFT(msgData json.RawMessage) ([]sdk.Msg, error) {
// 	var msg types.MsgMintNFT
// 	err := json.Unmarshal(msgData, &msg)
// 	if err != nil {
// 		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
// 	}

// 	return []sdk.Msg{msg}, nil
// }

// func handleMsgMintFT(msgData json.RawMessage) ([]sdk.Msg, error) {
// 	var msg types.MsgMintFT
// 	err := json.Unmarshal(msgData, &msg)
// 	if err != nil {
// 		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
// 	}

// 	return []sdk.Msg{msg}, nil
// }

// func handleMsgBurnNFT(msgData json.RawMessage) ([]sdk.Msg, error) {
// 	var msg types.MsgBurnNFT
// 	err := json.Unmarshal(msgData, &msg)
// 	if err != nil {
// 		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
// 	}

// 	return []sdk.Msg{msg}, nil
// }

// func handleMsgBurnNFTFrom(msgData json.RawMessage) ([]sdk.Msg, error) {
// 	var msg types.MsgBurnNFTFrom
// 	err := json.Unmarshal(msgData, &msg)
// 	if err != nil {
// 		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
// 	}

// 	return []sdk.Msg{msg}, nil
// }

// func handleMsgBurnFT(msgData json.RawMessage) ([]sdk.Msg, error) {
// 	var msg types.MsgBurnFT
// 	err := json.Unmarshal(msgData, &msg)
// 	if err != nil {
// 		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
// 	}

// 	return []sdk.Msg{msg}, nil
// }

// func handleMsgBurnFTFrom(msgData json.RawMessage) ([]sdk.Msg, error) {
// 	var msg types.MsgBurnFTFrom
// 	err := json.Unmarshal(msgData, &msg)
// 	if err != nil {
// 		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
// 	}

// 	return []sdk.Msg{msg}, nil
// }

// func handleMsgTransferNFT(msgData json.RawMessage) ([]sdk.Msg, error) {
// 	var msg types.MsgTransferNFT
// 	err := json.Unmarshal(msgData, &msg)
// 	if err != nil {
// 		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
// 	}

// 	return []sdk.Msg{msg}, nil
// }

// func handleMsgTransferNFTFrom(msgData json.RawMessage) ([]sdk.Msg, error) {
// 	var msg types.MsgTransferNFTFrom
// 	err := json.Unmarshal(msgData, &msg)
// 	if err != nil {
// 		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
// 	}

// 	return []sdk.Msg{msg}, nil
// }

// func handleMsgTransferFT(msgData json.RawMessage) ([]sdk.Msg, error) {
// 	var msg types.MsgTransferFT
// 	err := json.Unmarshal(msgData, &msg)
// 	if err != nil {
// 		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
// 	}

// 	return []sdk.Msg{msg}, nil
// }

// func handleMsgTransferFTFrom(msgData json.RawMessage) ([]sdk.Msg, error) {
// 	var msg types.MsgTransferFTFrom
// 	err := json.Unmarshal(msgData, &msg)
// 	if err != nil {
// 		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
// 	}

// 	return []sdk.Msg{msg}, nil
// }

// func handleMsgModify(msgData json.RawMessage) ([]sdk.Msg, error) {
// 	var msg types.MsgModify
// 	err := json.Unmarshal(msgData, &msg)
// 	if err != nil {
// 		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
// 	}

// 	return []sdk.Msg{msg}, nil
// }

// func handleMsgApprove(msgData json.RawMessage) ([]sdk.Msg, error) {
// 	var msg types.MsgApprove
// 	err := json.Unmarshal(msgData, &msg)
// 	if err != nil {
// 		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
// 	}

// 	return []sdk.Msg{msg}, nil
// }

// func handleMsgDisapprove(msgData json.RawMessage) ([]sdk.Msg, error) {
// 	var msg types.MsgDisapprove
// 	err := json.Unmarshal(msgData, &msg)
// 	if err != nil {
// 		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
// 	}

// 	return []sdk.Msg{msg}, nil
// }

// func handleMsgGrantPerm(msgData json.RawMessage) ([]sdk.Msg, error) {
// 	var msg types.MsgGrantPermission
// 	err := json.Unmarshal(msgData, &msg)
// 	if err != nil {
// 		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
// 	}

// 	return []sdk.Msg{msg}, nil
// }

// func handleMsgRevokePerm(msgData json.RawMessage) ([]sdk.Msg, error) {
// 	var msg types.MsgRevokePermission
// 	err := json.Unmarshal(msgData, &msg)
// 	if err != nil {
// 		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
// 	}

// 	return []sdk.Msg{msg}, nil
// }

// func handleMsgAttach(msgData json.RawMessage) ([]sdk.Msg, error) {
// 	var msg types.MsgAttach
// 	err := json.Unmarshal(msgData, &msg)
// 	if err != nil {
// 		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
// 	}

// 	return []sdk.Msg{msg}, nil
// }

// func handleMsgDetach(msgData json.RawMessage) ([]sdk.Msg, error) {
// 	var msg types.MsgDetach
// 	err := json.Unmarshal(msgData, &msg)
// 	if err != nil {
// 		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
// 	}

// 	return []sdk.Msg{msg}, nil
// }

// func handleMsgAttachFrom(msgData json.RawMessage) ([]sdk.Msg, error) {
// 	var msg types.MsgAttachFrom
// 	err := json.Unmarshal(msgData, &msg)
// 	if err != nil {
// 		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
// 	}

// 	return []sdk.Msg{msg}, nil
// }

// func handleMsgDetachFrom(msgData json.RawMessage) ([]sdk.Msg, error) {
// 	var msg types.MsgDetachFrom
// 	err := json.Unmarshal(msgData, &msg)
// 	if err != nil {
// 		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
// 	}

// 	return []sdk.Msg{msg}, nil
// }
