package keeper

import (
	"context"

	sdktypes "github.com/Finschia/finschia-rdk/types"
	"github.com/Finschia/finschia-rdk/x/or/da/types"
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

	if _, err := k.rollupKeeper.GetRollupInfo(ctx, msg.RollupName); err != nil {
		return nil, err
	}

	batch, err := k.DecompressCCBatch(ctx, msg.Batch)
	if err != nil {
		return nil, err
	}
	if len(batch.Frames) == 0 {
		return nil, types.ErrInvalidCCBatch.Wrapf("empty batch")
	}

	if err := k.SaveCCBatch(ctx, msg.RollupName, batch); err != nil {
		return nil, err
	}

	return &types.MsgAppendCCBatchResponse{}, nil
}

func (k msgServer) Enqueue(goCtx context.Context, msg *types.MsgEnqueue) (*types.MsgEnqueueResponse, error) {
	ctx := sdktypes.UnwrapSDKContext(goCtx)

	rollupInfo, err := k.rollupKeeper.GetRollupInfo(ctx, msg.RollupName)
	if err != nil {
		return nil, err
	}

	if msg.Txraw == nil {
		return nil, types.ErrInvalidQueueTx.Wrapf("empty tx")
	} else if uint64(len(msg.Txraw)) > k.MaxQueueTxSize(ctx) {
		return nil, types.ErrInvalidQueueTx.Wrapf("tx data size exceeds maximum for rollup tx")
	}

	if msg.GasLimit < k.MinQueueTxGas(ctx) {
		return nil, types.ErrInvalidQueueTx.Wrapf("gas limit too low to enqueue tx")
	}

	if err := k.SaveQueueTx(ctx, msg.RollupName, msg.Txraw, msg.GasLimit, rollupInfo.L1ToL2GasRatio); err != nil {
		return nil, err
	}

	return &types.MsgEnqueueResponse{}, nil
}

func (k msgServer) AppendSCCBatch(goCtx context.Context, msg *types.MsgAppendSCCBatch) (*types.MsgAppendSCCBatchResponse, error) {
	panic("implement me")
}

func (k msgServer) RemoveSCCBatch(goCtx context.Context, msg *types.MsgRemoveSCCBatch) (*types.MsgRemoveSCCBatchResponse, error) {
	panic("implement me")
}
