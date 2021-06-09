package keeper

import (
	"context"

	sdk "github.com/line/lfb-sdk/types"
	"github.com/line/lfb-sdk/x/auth/types"
)

type msgServer struct {
	AccountKeeper
}

// NewMsgServerImpl returns an implementation of the auth MsgServer interface for the provided Keeper.
func NewMsgServerImpl(keeper AccountKeeper) types.MsgServer {
	return &msgServer{AccountKeeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) Empty(goCtx context.Context, msg *types.MsgEmpty) (*types.MsgEmptyResponse, error) {
	sdk.UnwrapSDKContext(goCtx).EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.FromAddress),
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.EventEmpty),
		),
	)
	return &types.MsgEmptyResponse{}, nil
}
