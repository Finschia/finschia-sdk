package keeper

import (
	"context"

	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/feegrant/types"

	sdkerrors "github.com/line/lbm-sdk/types/errors"
	"github.com/line/lbm-sdk/x/feegrant"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the feegrant MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(k Keeper) feegrant.MsgServer {
	return &msgServer{
		Keeper: k,
	}
}

var _ feegrant.MsgServer = msgServer{}

// GrantAllowance grants an allowance from the granter's funds to be used by the grantee.
func (k msgServer) GrantAllowance(goCtx context.Context, msg *feegrant.MsgGrantAllowance) (*feegrant.MsgGrantAllowanceResponse, error) {
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

	// Checking for duplicate entry
	if f, _ := k.Keeper.GetAllowance(ctx, granter, grantee); f != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "fee allowance already exists")
	}

	allowance, err := msg.GetFeeAllowanceI()
	if err != nil {
		return nil, err
	}

	err = k.Keeper.GrantAllowance(ctx, granter, grantee, allowance)
	if err != nil {
		return nil, err
	}

	return &feegrant.MsgGrantAllowanceResponse{}, nil
}

// RevokeAllowance revokes a fee allowance between a granter and grantee.
func (k msgServer) RevokeAllowance(goCtx context.Context, msg *feegrant.MsgRevokeAllowance) (*feegrant.MsgRevokeAllowanceResponse, error) {
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

	err = k.Keeper.revokeAllowance(ctx, granter, grantee)
	if err != nil {
		return nil, err
	}

	return &feegrant.MsgRevokeAllowanceResponse{}, nil
}
