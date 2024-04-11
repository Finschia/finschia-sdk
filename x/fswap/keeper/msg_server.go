package keeper

import (
	"context"

	"github.com/Finschia/finschia-sdk/x/fswap/types"
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

// Swap implements types.MsgServer.
func (m msgServer) Swap(context.Context, *types.MsgSwapRequest) (*types.MsgSwapResponse, error) {
	panic("unimplemented")
}

// SwapAll implements types.MsgServer.
func (m msgServer) SwapAll(context.Context, *types.MsgSwapAllRequest) (*types.MsgSwapAllResponse, error) {
	panic("unimplemented")
}
