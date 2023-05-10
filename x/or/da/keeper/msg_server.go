package keeper

import (
	"context"

	sdktypes "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/or/da/types"
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

func (k msgServer) UpdateParams(goCtx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	ctx := sdktypes.UnwrapSDKContext(goCtx)
	if err := k.validateAuthority(msg.Authority); err != nil {
		return nil, err
	}

	if err := msg.Params.Validate(); err != nil {
		return nil, err
	}

	if err := k.SetParams(ctx, msg.Params); err != nil {
		return nil, err
	}

	return &types.MsgUpdateParamsResponse{}, nil
}

func (k msgServer) AppendCTCBatch(ctx context.Context, batch *types.MsgAppendCTCBatch) (*types.MsgAppendCTCBatchResponse, error) {
	k.appendSequencerBatch()
	panic("implement me")
}
