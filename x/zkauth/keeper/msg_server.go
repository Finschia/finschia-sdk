package keeper

import (
	"context"

	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/zkauth/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) Execution(goCtx context.Context, msg *types.MsgExecution) (*types.MsgExecutionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	msgs, err := msg.GetMessages()
	if err != nil {
		return nil, err
	}

	// zksigner is assumed to be one
	zksigner := msg.GetSigners()[0]

	results, err := k.DispatchMsgs(ctx, msgs, zksigner)
	if err != nil {
		return nil, err
	}

	return &types.MsgExecutionResponse{Results: results}, nil
}
