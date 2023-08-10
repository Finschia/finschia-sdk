package keeper

import (
	"bytes"
	"compress/zlib"
	"crypto/sha256"
	"encoding/binary"
	"io"
	"time"

	sdktypes "github.com/Finschia/finschia-sdk/types"
	sdkerror "github.com/Finschia/finschia-sdk/types/errors"
	"github.com/Finschia/finschia-sdk/x/or/da/types"
)

func (k Keeper) SaveQueueTx(ctx sdktypes.Context, rollupName string, tx []byte, gasLimit, L1ToL2GasRatio uint64) error {
	// Transactions submitted to the queue lack a method for paying gas fees to the Sequencer.
	// For transaction with a high L2 gas limit, we burn some extra gas on L1.
	gasToConsume := (gasLimit - k.EnqueueL2GasPrepaid(ctx, tx)) / L1ToL2GasRatio
	ctx.GasMeter().ConsumeGas(gasToConsume, "enqueue tx")

	var queueState *types.QueueTxState
	queueState, _ = k.GetQueueTxState(ctx, rollupName)
	if queueState == nil {
		queueState = &types.QueueTxState{
			ProcessedQueueIndex: 0,
			NextQueueIndex:      1,
		}
	}

	qtx := &types.L1ToL2Queue{
		Timestamp: ctx.BlockTime(),
		L1Height:  ctx.BlockHeight(),
		Txraw:     tx,
		Status:    types.QUEUE_TX_PENDING,
	}

	k.setQueueTx(ctx, rollupName, queueState.NextQueueIndex, qtx)
	queueState.NextQueueIndex++

	k.setQueueTxState(ctx, rollupName, queueState)

	if err := ctx.EventManager().EmitTypedEvent(&types.EventSaveQueueTx{
		RollupName:       rollupName,
		NextQueueIndex:   queueState.NextQueueIndex,
		ExtraConsumedGas: gasToConsume,
		L2GasLimit:       gasLimit,
	}); err != nil {
		return err
	}

	return nil
}

func (k Keeper) SaveCCBatch(ctx sdktypes.Context, rollupName string, batch *types.CCBatch) error {
	var ccState *types.CCState
	ccState, _ = k.GetCCState(ctx, rollupName)
	if ccState == nil {
		if batch.ShouldStartAtFrame != 0 {
			return types.ErrInvalidCCBatch.Wrapf("cc state not found: this batch should start at frame 0, but start at %d", batch.ShouldStartAtFrame)
		}
		ccState = &types.CCState{
			Base:             1,
			Height:           0,
			ProcessedL2Block: 0,
			Timestamp:        time.Unix(0, 0).UTC(),
			L1Height:         0,
		}
	}

	var queueState *types.QueueTxState
	queueState, _ = k.GetQueueTxState(ctx, rollupName)
	if queueState == nil {
		queueState = &types.QueueTxState{
			ProcessedQueueIndex: 0,
			NextQueueIndex:      1,
		}
	}

	var totalFrames uint64
	if ccState.Height == 0 {
		totalFrames = uint64(len(batch.Frames))
	} else {
		prevCCRef, err := k.GetCCRef(ctx, rollupName, ccState.Height)
		if err != nil {
			return err
		}

		if prevCCRef.TotalFrames != batch.ShouldStartAtFrame {
			return types.ErrInvalidCCBatch.Wrapf("batch should start at frame %d but this batch start at frame %d", prevCCRef.TotalFrames, batch.ShouldStartAtFrame)
		}
		totalFrames = prevCCRef.TotalFrames + uint64(len(batch.Frames))
	}

	batchHash := sha256.Sum256(k.cdc.MustMarshal(batch))
	txHash := sha256.Sum256(ctx.TxBytes())
	ref := types.NewCCRef(txHash[:], uint32(ctx.MsgIndex()), uint32(len(batch.Frames)), totalFrames, batchHash[:])

	// start to process batch frames
	ccState.Height++

	for i, frame := range batch.Frames {
		if len(frame.Elements) == 0 {
			return types.ErrInvalidCCBatch.Wrapf("frame %d has empty elements", i)
		}

		if frame.Header.GetL2Height() != ccState.ProcessedL2Block+1 {
			return types.ErrInvalidCCBatch.Wrapf("frame %d has invalid l2 height %d, expected %d", i, frame.Header.GetL2Height(), ccState.ProcessedL2Block+1)
		}
		if frame.Header.GetL2Height() != 1 && frame.Header.GetParentHash() == nil {
			return types.ErrInvalidCCBatch.Wrapf("frame %d has nil parent hash", i)
		}
		if !ccState.Timestamp.Equal(time.Unix(0, 0).UTC()) && frame.Header.Timestamp.Before(ccState.GetTimestamp()) {
			return types.ErrInvalidCCBatch.Wrapf("frame %d is outdated: %s (frame) < %s (CCState)", i, frame.Header.GetTimestamp(), ccState.GetTimestamp())
		}
		if frame.Header.GetL1Height() < ccState.GetL1Height() {
			return types.ErrInvalidCCBatch.Wrapf("frame %d is outdated: %d (frame) < %d (CCState)", i, frame.Header.GetL1Height(), ccState.GetL1Height())
		}

		for j, elem := range frame.Elements {
			if elem.Txraw == nil {
				if elem.QueueIndex < 1 || elem.QueueIndex != queueState.ProcessedQueueIndex+1 {
					return types.ErrInvalidCCBatch.Wrapf("queue index of frame %d element %d is %d, expected %d", i, j, elem.QueueIndex, queueState.ProcessedQueueIndex+1)
				}

				qtx, err := k.GetQueueTx(ctx, rollupName, elem.QueueIndex)
				if err != nil {
					return err
				}
				qtx.Status = types.QUEUE_TX_SUBMITTED
				k.setQueueTx(ctx, rollupName, elem.QueueIndex, qtx)

				queueState.ProcessedQueueIndex++
			}
		}

		ccState.ProcessedL2Block++
		ccState.Timestamp = frame.Header.Timestamp
		ccState.L1Height = frame.Header.L1Height
		k.setL2HeightBatchMap(ctx, rollupName, ccState.ProcessedL2Block, ccState.Height)
	}

	k.setCCRef(ctx, rollupName, ccState.Height, ref)
	k.setCCState(ctx, rollupName, ccState)
	k.setQueueTxState(ctx, rollupName, queueState)

	if err := ctx.EventManager().EmitTypedEvent(&types.EventAppendCCBatch{
		RollupName:          rollupName,
		BatchIndex:          ccState.Height,
		ProcessedQueueIndex: queueState.ProcessedQueueIndex,
		TotalFrames:         ref.TotalFrames,
		BatchSize:           ref.BatchSize,
		BatchHash:           ref.BatchRoot,
		ProcessedL2Block:    ccState.ProcessedL2Block,
	}); err != nil {
		return err
	}

	return nil
}

func (k Keeper) UpdateQueueTxsStatus(ctx sdktypes.Context) error {
	rollupList := k.rollupKeeper.GetAllRollup(ctx)
	if rollupList == nil {
		return nil
	}

	for _, rollup := range rollupList {
		state, err := k.GetQueueTxState(ctx, rollup.RollupName)
		if err != nil {
			continue
		}

		if err := k.processQueueTxs(ctx, state, rollup.RollupName); err != nil {
			return err
		}
	}

	return nil
}

func (k Keeper) processQueueTxs(ctx sdktypes.Context, state *types.QueueTxState, name string) error {
	qi := state.ProcessedQueueIndex + 1
	for ; qi < state.NextQueueIndex; qi++ {
		qtx, err := k.GetQueueTx(ctx, name, qi)
		if err != nil {
			return types.ErrQueueTxStateNotFound.Wrapf("rollup %s queue tx index %d", name, qi)
		}

		if (uint64(qtx.L1Height) + k.QueueTxExpirationWindow(ctx)) > uint64(ctx.BlockHeight()) {
			break
		}

		if qtx.Status == types.QUEUE_TX_PENDING {
			qtx.Status = types.QUEUE_TX_EXPIRED
			k.setQueueTx(ctx, name, qi, qtx)
			// TODO: slash registered sequencers
		}
	}
	state.ProcessedQueueIndex = qi - 1
	k.setQueueTxState(ctx, name, state)

	return nil
}

func (k Keeper) DecompressCCBatch(ctx sdktypes.Context, origin types.CompressedCCBatch) (*types.CCBatch, error) {
	p := k.GetParams(ctx)
	if uint64(len(origin.Data)) > p.CCBatchMaxBytes {
		return nil, types.ErrInvalidCompressedData.Wrapf("compressed data size %d exceeds max batch size %d", len(origin.Data), p.CCBatchMaxBytes)
	}

	switch origin.Compression {
	case types.OptionZLIB:
		b := bytes.NewReader(origin.Data)
		r, err := zlib.NewReader(b)
		defer r.Close()
		if err != nil {
			return nil, types.ErrInvalidCompressedData.Wrap(err.Error())
		}

		out := make([]byte, 0)
		for {
			buf := make([]byte, p.CCBatchMaxBytes)
			n, err := r.Read(buf)
			out = append(out, buf[:n]...)
			if err != nil {
				if err == io.EOF {
					break
				}
				return nil, err
			}
		}
		batch := new(types.CCBatch)
		k.cdc.MustUnmarshal(out, batch)

		return batch, nil

	case types.OptionZSTD:
		return nil, types.ErrInvalidCompressedData.Wrapf("compression %s not supported", origin.Compression)
	case types.OptionEmpty:
		return nil, types.ErrInvalidCompressedData.Wrapf("batch data must be compressed")
	default:
		return nil, sdkerror.ErrInvalidRequest.Wrapf("no compression option provided")
	}
}

func (k Keeper) EnqueueL2GasPrepaid(ctx sdktypes.Context, tx []byte) uint64 {
	return k.accountKeeper.GetParams(ctx).TxSizeCostPerByte * uint64(len(tx)+30)
}

func (k Keeper) GetCCState(ctx sdktypes.Context, rollupName string) (*types.CCState, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetCCStateStoreKey(rollupName))
	if bz == nil {
		return nil, types.ErrCCStateNotFound
	}
	state := new(types.CCState)
	k.cdc.MustUnmarshal(bz, state)
	return state, nil
}

func (k Keeper) setCCState(ctx sdktypes.Context, rollupName string, state *types.CCState) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(state)
	store.Set(types.GetCCStateStoreKey(rollupName), bz)
}

func (k Keeper) GetCCRef(ctx sdktypes.Context, rollupName string, idx uint64) (*types.CCRef, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetCCBatchIndexKey(rollupName, idx))
	if bz == nil {
		return nil, types.ErrCCRefNotFound
	}
	ref := new(types.CCRef)
	k.cdc.MustUnmarshal(bz, ref)
	return ref, nil
}

func (k Keeper) setCCRef(ctx sdktypes.Context, rollupName string, idx uint64, ref *types.CCRef) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(ref)
	store.Set(types.GetCCBatchIndexKey(rollupName, idx), bz)
}

func (k Keeper) GetQueueTxState(ctx sdktypes.Context, rollupName string) (*types.QueueTxState, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetQueueTxStateStoreKey(rollupName))
	if bz == nil {
		return nil, types.ErrQueueTxStateNotFound
	}
	state := new(types.QueueTxState)
	k.cdc.MustUnmarshal(bz, state)
	return state, nil
}

func (k Keeper) setQueueTxState(ctx sdktypes.Context, rollupName string, state *types.QueueTxState) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(state)
	store.Set(types.GetQueueTxStateStoreKey(rollupName), bz)
}

func (k Keeper) GetQueueTx(ctx sdktypes.Context, rollupName string, idx uint64) (*types.L1ToL2Queue, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetCCQueueTxKey(rollupName, idx))
	if bz == nil {
		return nil, types.ErrQueueTxNotFound
	}
	tx := new(types.L1ToL2Queue)
	k.cdc.MustUnmarshal(bz, tx)
	return tx, nil
}

func (k Keeper) setQueueTx(ctx sdktypes.Context, rollupName string, idx uint64, elem *types.L1ToL2Queue) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(elem)
	store.Set(types.GetCCQueueTxKey(rollupName, idx), bz)
}

func (k Keeper) GetL2HeightBatchMap(ctx sdktypes.Context, rollupName string, l2height uint64) (uint64, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetCCL2HeightToBatchKey(rollupName, l2height))
	if bz == nil {
		return 0, types.ErrL2HeightBatchMapNotFound
	}
	batchIdx := binary.BigEndian.Uint64(bz)
	return batchIdx, nil
}

func (k Keeper) setL2HeightBatchMap(ctx sdktypes.Context, rollupName string, l2height, batchIdx uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, batchIdx)
	store.Set(types.GetCCL2HeightToBatchKey(rollupName, l2height), bz)
}
