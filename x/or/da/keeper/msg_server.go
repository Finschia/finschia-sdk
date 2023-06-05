package keeper

import (
	"context"
	"crypto/sha256"

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
	if err := k.validateGovAuthority(msg.Authority); err != nil {
		return nil, err
	}

	if err := msg.Params.Validate(); err != nil {
		return nil, err
	}

	ctx := sdktypes.UnwrapSDKContext(goCtx)
	if err := k.SetParams(ctx, msg.Params); err != nil {
		return nil, err
	}

	if err := ctx.EventManager().EmitTypedEvent(&types.EventUpdateParams{Params: msg.Params}); err != nil {
		return nil, err
	}

	return &types.MsgUpdateParamsResponse{}, nil
}

func (k msgServer) AppendCCBatch(goCtx context.Context, msg *types.MsgAppendCCBatch) (*types.MsgAppendCCBatchResponse, error) {
	if err := k.validateSequencerAuthority(msg.FromAddress); err != nil {
		return nil, err
	}

	ctx := sdktypes.UnwrapSDKContext(goCtx)

	batch, err := k.DecompressCCBatch(msg.Batch)
	if err != nil {
		return nil, err
	}
	_ = batch
	_ = sha256.Sum256(ctx.TxBytes())
	k.appendSequencerBatch()
	ctx.HeaderHash()
	return &types.MsgAppendCCBatchResponse{}, nil
}

func (k msgServer) Enqueue(goCtx context.Context, msg *types.MsgEnqueue) (*types.MsgEnqueueResponse, error) {
	panic("implement me")
}

func (k msgServer) AppendSCCBatch(goCtx context.Context, msg *types.MsgAppendSCCBatch) (*types.MsgAppendSCCBatchResponse, error) {
	panic("implement me")
}

func (k msgServer) RemoveSCCBatch(goCtx context.Context, msg *types.MsgRemoveSCCBatch) (*types.MsgRemoveSCCBatchResponse, error) {
	panic("implement me")
}
