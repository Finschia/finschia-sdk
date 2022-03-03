package keeper

import (
	"context"

	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/authz/types"
)

var _ types.MsgServer = Keeper{}

// GrantAuthorization implements the MsgServer.GrantAuthorization method.
func (k Keeper) GrantAuthorization(goCtx context.Context, msg *types.MsgGrantAuthorizationRequest) (*types.MsgGrantAuthorizationResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	err := sdk.ValidateAccAddress(msg.Grantee)
	if err != nil {
		return nil, err
	}
	grantee := sdk.AccAddress(msg.Grantee)

	err = sdk.ValidateAccAddress(msg.Granter)
	if err != nil {
		return nil, err
	}
	granter := sdk.AccAddress(msg.Granter)

	authorization := msg.GetGrantAuthorization()
	// If the granted service Msg doesn't exist, we throw an error.
	if k.router.Handler(authorization.MethodName()) == nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "%s doesn't exist.", authorization.MethodName())
	}

	err = k.Grant(ctx, grantee, granter, authorization, msg.Expiration)
	if err != nil {
		return nil, err
	}

	return &types.MsgGrantAuthorizationResponse{}, nil
}

// RevokeAuthorization implements the MsgServer.RevokeAuthorization method.
func (k Keeper) RevokeAuthorization(goCtx context.Context, msg *types.MsgRevokeAuthorizationRequest) (*types.MsgRevokeAuthorizationResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	err := sdk.ValidateAccAddress(msg.Grantee)
	if err != nil {
		return nil, err
	}
	grantee := sdk.AccAddress(msg.Grantee)

	err = sdk.ValidateAccAddress(msg.Granter)
	if err != nil {
		return nil, err
	}
	granter := sdk.AccAddress(msg.Granter)

	err = k.Revoke(ctx, grantee, granter, msg.MethodName)
	if err != nil {
		return nil, err
	}

	return &types.MsgRevokeAuthorizationResponse{}, nil
}

// ExecAuthorized implements the MsgServer.ExecAuthorized method.
func (k Keeper) ExecAuthorized(goCtx context.Context, msg *types.MsgExecAuthorizedRequest) (*types.MsgExecAuthorizedResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	err := sdk.ValidateAccAddress(msg.Grantee)
	if err != nil {
		return nil, err
	}
	grantee := sdk.AccAddress(msg.Grantee)
	msgs, err := msg.GetServiceMsgs()
	if err != nil {
		return nil, err
	}
	result, err := k.DispatchActions(ctx, grantee, msgs)
	if err != nil {
		return nil, err
	}
	return &types.MsgExecAuthorizedResponse{Result: result}, nil
}
