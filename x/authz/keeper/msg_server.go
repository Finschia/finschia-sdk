package keeper

import (
	"context"

	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/authz"
)

var _ authz.MsgServer = Keeper{}

// GrantAuthorization implements the MsgServer.Grant method.
func (k Keeper) Grant(goCtx context.Context, msg *authz.MsgGrant) (*authz.MsgGrantResponse, error) {
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

	authorization := msg.GetAuthorization()
	if authorization == nil {
		return nil, sdkerrors.ErrUnpackAny.Wrap("Authorization is not present in the msg")
	}
	t := authorization.MsgTypeURL()
	if k.router.HandlerByTypeURL(t) == nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "%s doesn't exist.", t)
	}

	err = k.SaveGrant(ctx, grantee, granter, authorization, msg.Grant.Expiration)
	if err != nil {
		return nil, err
	}

	return &authz.MsgGrantResponse{}, nil
}

// RevokeAuthorization implements the MsgServer.Revoke method.
func (k Keeper) Revoke(goCtx context.Context, msg *authz.MsgRevoke) (*authz.MsgRevokeResponse, error) {
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

	err = k.DeleteGrant(ctx, grantee, granter, msg.MsgTypeUrl)
	if err != nil {
		return nil, err
	}

	return &authz.MsgRevokeResponse{}, nil
}

// Exec implements the MsgServer.Exec method.
func (k Keeper) Exec(goCtx context.Context, msg *authz.MsgExec) (*authz.MsgExecResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	err := sdk.ValidateAccAddress(msg.Grantee)
	if err != nil {
		return nil, err
	}
	grantee := sdk.AccAddress(msg.Grantee)
	msgs, err := msg.GetMessages()
	if err != nil {
		return nil, err
	}
	result, err := k.DispatchActions(ctx, grantee, msgs)
	if err != nil {
		return nil, err
	}
	return &authz.MsgExecResponse{Result: result}, nil
}
