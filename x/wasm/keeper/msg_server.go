package keeper

import (
	"context"

	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/wasm/types"
)

var _ types.MsgServer = msgServer{}

type msgServer struct {
	keeper types.ContractOpsKeeper
}

func NewMsgServerImpl(k types.ContractOpsKeeper) types.MsgServer {
	return &msgServer{keeper: k}
}

func (m msgServer) StoreCode(goCtx context.Context, msg *types.MsgStoreCode) (*types.MsgStoreCodeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	err := sdk.ValidateAccAddress(msg.Sender)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "sender")
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
	))

	codeID, err := m.keeper.Create(ctx, sdk.AccAddress(msg.Sender), msg.WASMByteCode, msg.InstantiatePermission)
	if err != nil {
		return nil, err
	}

	return &types.MsgStoreCodeResponse{
		CodeID: codeID,
	}, nil
}

func (m msgServer) InstantiateContract(goCtx context.Context, msg *types.MsgInstantiateContract) (*types.MsgInstantiateContractResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	err := sdk.ValidateAccAddress(msg.Sender)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "sender")
	}
	var adminAddr sdk.AccAddress
	if msg.Admin != "" {
		if err = sdk.ValidateAccAddress(msg.Admin); err != nil {
			return nil, sdkerrors.Wrap(err, "admin")
		}
		adminAddr = sdk.AccAddress(msg.Admin)
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
	))

	contractAddr, data, err := m.keeper.Instantiate(ctx, msg.CodeID, sdk.AccAddress(msg.Sender), adminAddr, msg.Msg, msg.Label, msg.Funds)
	if err != nil {
		return nil, err
	}

	return &types.MsgInstantiateContractResponse{
		Address: contractAddr.String(),
		Data:    data,
	}, nil
}

func (m msgServer) StoreCodeAndInstantiateContract(goCtx context.Context,
	msg *types.MsgStoreCodeAndInstantiateContract) (*types.MsgStoreCodeAndInstantiateContractResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	err := sdk.ValidateAccAddress(msg.Sender)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "sender")
	}
	codeID, err := m.keeper.Create(ctx, sdk.AccAddress(msg.Sender), msg.WASMByteCode, msg.InstantiatePermission)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
	))

	var adminAddr sdk.AccAddress
	if msg.Admin != "" {
		if err = sdk.ValidateAccAddress(msg.Admin); err != nil {
			return nil, sdkerrors.Wrap(err, "admin")
		}
		adminAddr = sdk.AccAddress(msg.Admin)
	}

	contractAddr, data, err := m.keeper.Instantiate(ctx, codeID, sdk.AccAddress(msg.Sender), adminAddr, msg.InitMsg,
		msg.Label, msg.Funds)
	if err != nil {
		return nil, err
	}

	return &types.MsgStoreCodeAndInstantiateContractResponse{
		CodeID:  codeID,
		Address: contractAddr.String(),
		Data:    data,
	}, nil
}

func (m msgServer) ExecuteContract(goCtx context.Context, msg *types.MsgExecuteContract) (*types.MsgExecuteContractResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	err := sdk.ValidateAccAddress(msg.Sender)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "sender")
	}
	err = sdk.ValidateAccAddress(msg.Contract)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "contract")
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
	))

	data, err := m.keeper.Execute(ctx, sdk.AccAddress(msg.Contract), sdk.AccAddress(msg.Sender), msg.Msg, msg.Funds)
	if err != nil {
		return nil, err
	}

	return &types.MsgExecuteContractResponse{
		Data: data,
	}, nil
}

func (m msgServer) MigrateContract(goCtx context.Context, msg *types.MsgMigrateContract) (*types.MsgMigrateContractResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	err := sdk.ValidateAccAddress(msg.Sender)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "sender")
	}
	err = sdk.ValidateAccAddress(msg.Contract)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "contract")
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
	))

	data, err := m.keeper.Migrate(ctx, sdk.AccAddress(msg.Contract), sdk.AccAddress(msg.Sender), msg.CodeID, msg.Msg)
	if err != nil {
		return nil, err
	}

	return &types.MsgMigrateContractResponse{
		Data: data,
	}, nil
}

func (m msgServer) UpdateAdmin(goCtx context.Context, msg *types.MsgUpdateAdmin) (*types.MsgUpdateAdminResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	err := sdk.ValidateAccAddress(msg.Sender)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "sender")
	}
	err = sdk.ValidateAccAddress(msg.Contract)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "contract")
	}
	err = sdk.ValidateAccAddress(msg.NewAdmin)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "new admin")
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
	))

	if err := m.keeper.UpdateContractAdmin(ctx, sdk.AccAddress(msg.Contract), sdk.AccAddress(msg.Sender),
		sdk.AccAddress(msg.NewAdmin)); err != nil {
		return nil, err
	}

	return &types.MsgUpdateAdminResponse{}, nil
}

func (m msgServer) ClearAdmin(goCtx context.Context, msg *types.MsgClearAdmin) (*types.MsgClearAdminResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	err := sdk.ValidateAccAddress(msg.Sender)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "sender")
	}
	err = sdk.ValidateAccAddress(msg.Contract)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "contract")
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
	))
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeUpdateAdmin,
		sdk.NewAttribute(types.AttributeKeyContractAddr, msg.Contract),
	))

	if err := m.keeper.ClearContractAdmin(ctx, sdk.AccAddress(msg.Contract), sdk.AccAddress(msg.Sender)); err != nil {
		return nil, err
	}

	return &types.MsgClearAdminResponse{}, nil
}

// UpdateContractStatus handles MsgUpdateContractStatus
// CONTRACT: msg.validateBasic() must be called before calling this
func (m msgServer) UpdateContractStatus(goCtx context.Context, msg *types.MsgUpdateContractStatus) (*types.MsgUpdateContractStatusResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	err := sdk.ValidateAccAddress(msg.Sender)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "sender")
	}
	err = sdk.ValidateAccAddress(msg.Contract)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "contract")
	}

	if err = m.keeper.UpdateContractStatus(ctx, sdk.AccAddress(msg.Contract), sdk.AccAddress(msg.Sender),
		msg.Status); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
		),
		sdk.NewEvent(
			types.EventTypeUpdateContractStatus,
			sdk.NewAttribute(types.AttributeKeyContractAddr, msg.Contract),
			sdk.NewAttribute(types.AttributeKeyContractStatus, msg.Status.String()),
		),
	})

	return &types.MsgUpdateContractStatusResponse{}, nil
}
