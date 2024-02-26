package keeper

import (
	"context"

	"github.com/Finschia/finschia-sdk/x/zkauth/types"
)

type msgServer struct {
	*Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper *Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) Execution(goCtx context.Context, msg *types.MsgExecution) (*types.MsgExecutionResponse, error) {
	return &types.MsgExecutionResponse{}, nil
}
