package keeper

import (
	"context"
	"fmt"

	sdk "github.com/line/lfb-sdk/types"
	sdkerrors "github.com/line/lfb-sdk/types/errors"
	"github.com/line/lfb-sdk/x/wasm/internal/types"
)

var _ types.MsgServer = msgServer{}

type msgServer struct {
	keeper *Keeper
}

func NewMsgServerImpl(k *Keeper) types.MsgServer {
	return &msgServer{keeper: k}
}

func (m msgServer) StoreCode(goCtx context.Context, msg *types.MsgStoreCode) (*types.MsgStoreCodeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "sender")
	}
	codeID, err := m.keeper.Create(ctx, senderAddr, msg.WASMByteCode, msg.Source, msg.Builder, msg.InstantiatePermission)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(types.AttributeKeySigner, msg.Sender),
		sdk.NewAttribute(types.AttributeKeyCodeID, fmt.Sprintf("%d", codeID)),
	))

	return &types.MsgStoreCodeResponse{
		CodeID: codeID,
	}, nil
}

func (m msgServer) InstantiateContract(goCtx context.Context, msg *types.MsgInstantiateContract) (*types.MsgInstantiateContractResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "sender")
	}
	var adminAddr sdk.AccAddress
	if msg.Admin != "" {
		if adminAddr, err = sdk.AccAddressFromBech32(msg.Admin); err != nil {
			return nil, sdkerrors.Wrap(err, "admin")
		}
	}

	contractAddr, data, err := m.keeper.Instantiate(ctx, msg.CodeID, senderAddr, adminAddr, msg.InitMsg, msg.Label, msg.Funds)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(types.AttributeKeySigner, msg.Sender),
		sdk.NewAttribute(types.AttributeKeyCodeID, fmt.Sprintf("%d", msg.CodeID)),
		sdk.NewAttribute(types.AttributeKeyContract, contractAddr.String()),
	))

	return &types.MsgInstantiateContractResponse{
		Address: contractAddr.String(),
		Data:    data,
	}, nil
}

func (m msgServer) StoreCodeAndInstantiateContract(goCtx context.Context, msg *types.MsgStoreCodeAndInstantiateContract) (*types.MsgStoreCodeAndInstantiateContractResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "sender")
	}
	codeID, err := m.keeper.Create(ctx, senderAddr, msg.WASMByteCode, msg.Source, msg.Builder, msg.InstantiatePermission)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(types.AttributeKeySigner, msg.Sender),
		sdk.NewAttribute(types.AttributeKeyCodeID, fmt.Sprintf("%d", codeID)),
	))

	var adminAddr sdk.AccAddress
	if msg.Admin != "" {
		if adminAddr, err = sdk.AccAddressFromBech32(msg.Admin); err != nil {
			return nil, sdkerrors.Wrap(err, "admin")
		}
	}

	contractAddr, data, err := m.keeper.Instantiate(ctx, codeID, senderAddr, adminAddr, msg.InitMsg, msg.Label, msg.Funds)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(types.AttributeKeySigner, msg.Sender),
		sdk.NewAttribute(types.AttributeKeyCodeID, fmt.Sprintf("%d", codeID)),
		sdk.NewAttribute(types.AttributeKeyContract, contractAddr.String()),
	))

	return &types.MsgStoreCodeAndInstantiateContractResponse{
		CodeID:  codeID,
		Address: contractAddr.String(),
		Data:    data,
	}, nil
}

func (m msgServer) ExecuteContract(goCtx context.Context, msg *types.MsgExecuteContract) (*types.MsgExecuteContractResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "sender")
	}
	contractAddr, err := sdk.AccAddressFromBech32(msg.Contract)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "contract")
	}

	res, err := m.keeper.Execute(ctx, contractAddr, senderAddr, msg.Msg, msg.Funds)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(types.AttributeKeySigner, msg.Sender),
		sdk.NewAttribute(types.AttributeKeyContract, msg.Contract),
	))

	return &types.MsgExecuteContractResponse{
		Data: res.Data,
	}, nil

}

func (m msgServer) MigrateContract(goCtx context.Context, msg *types.MsgMigrateContract) (*types.MsgMigrateContractResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "sender")
	}
	contractAddr, err := sdk.AccAddressFromBech32(msg.Contract)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "contract")
	}

	res, err := m.keeper.Migrate(ctx, contractAddr, senderAddr, msg.CodeID, msg.MigrateMsg)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(types.AttributeKeySigner, msg.Sender),
		sdk.NewAttribute(types.AttributeKeyContract, msg.Contract),
	))

	return &types.MsgMigrateContractResponse{
		Data: res.Data,
	}, nil
}

func (m msgServer) UpdateAdmin(goCtx context.Context, msg *types.MsgUpdateAdmin) (*types.MsgUpdateAdminResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "sender")
	}
	contractAddr, err := sdk.AccAddressFromBech32(msg.Contract)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "contract")
	}
	newAdminAddr, err := sdk.AccAddressFromBech32(msg.NewAdmin)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "new admin")
	}

	if err := m.keeper.UpdateContractAdmin(ctx, contractAddr, senderAddr, newAdminAddr); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(types.AttributeKeySigner, msg.Sender),
		sdk.NewAttribute(types.AttributeKeyContract, msg.Contract),
	))

	return &types.MsgUpdateAdminResponse{}, nil
}

func (m msgServer) ClearAdmin(goCtx context.Context, msg *types.MsgClearAdmin) (*types.MsgClearAdminResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "sender")
	}
	contractAddr, err := sdk.AccAddressFromBech32(msg.Contract)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "contract")
	}

	if err := m.keeper.ClearContractAdmin(ctx, contractAddr, senderAddr); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(types.AttributeKeySigner, msg.Sender),
		sdk.NewAttribute(types.AttributeKeyContract, msg.Contract),
	))

	return &types.MsgClearAdminResponse{}, nil
}

// UpdateContractStatus handles MsgUpdateContractStatus
// CONTRACT: msg.validateBasic() must be called before calling this
func (m msgServer) UpdateContractStatus(goCtx context.Context, msg *types.MsgUpdateContractStatus) (*types.MsgUpdateContractStatusResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	senderAddr, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "sender")
	}
	contractAddr, err := sdk.AccAddressFromBech32(msg.Contract)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "contract")
	}

	if err = m.keeper.UpdateContractStatus(ctx, contractAddr, senderAddr, msg.Status); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeySigner, msg.Sender),
			sdk.NewAttribute(types.AttributeKeyContract, msg.Contract),
		),
		sdk.NewEvent(
			types.EventTypeUpdateContractStatus,
			sdk.NewAttribute(types.AttributeKeyContract, msg.Contract),
			sdk.NewAttribute(types.AttributeKeyContractStatus, msg.Status.String()),
		),
	})

	return &types.MsgUpdateContractStatusResponse{}, nil
}
