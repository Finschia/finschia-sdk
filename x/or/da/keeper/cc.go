package keeper

import (
	"bytes"
	"compress/zlib"
	"crypto/sha256"
	"io"

	sdktypes "github.com/Finschia/finschia-sdk/types"
	sdkerror "github.com/Finschia/finschia-sdk/types/errors"
	"github.com/Finschia/finschia-sdk/x/or/da/types"
)

func (k Keeper) SaveCCBatch(ctx sdktypes.Context, rollupName string, batch *types.CCBatch) error {
	var ccState *types.CCState
	ccState, _ = k.GetCCState(ctx, rollupName)
	if ccState == nil {
		if batch.ShouldStartAtFrame != 0 {
			return types.ErrInvalidCCBatch.Wrapf("batch should start at frame %d but cc state not found", batch.ShouldStartAtFrame)
		}
		ccState = &types.CCState{
			Base:           1,
			NextQueueIndex: 1,
		}
	}

	var totalFrames uint64
	prevCCRef, err := k.GetCCRef(ctx, rollupName, ccState.Height)
	if prevCCRef == nil {
		totalFrames = uint64(len(batch.Frames))
	} else if err == nil {
		if prevCCRef.TotalFrames != (batch.ShouldStartAtFrame + 1) {
			return types.ErrInvalidCCBatch.Wrapf("batch should start at frame %d but this batch start at frame %d", prevCCRef.TotalFrames, batch.ShouldStartAtFrame)
		}

		totalFrames = prevCCRef.TotalFrames + uint64(len(batch.Frames))
	}

	batchHash := sha256.Sum256(k.cdc.MustMarshal(batch))
	ref := types.NewCCRef(ctx.TxBytes(), uint32(ctx.MsgIndex()), uint32(len(batch.Frames)), totalFrames, batchHash[:])

	for i, frame := range batch.Frames {
		if len(frame.Elements) == 0 {
			return types.ErrInvalidCCBatch.Wrapf("frame %d has empty elements", i)
		}

		if frame.Header.GetParentHash() == nil {
			return types.ErrInvalidCCBatch.Wrapf("frame %d has nil parent hash", i)
		}
		if frame.Header.Timestamp.Before(ccState.GetTimestamp()) {
			return types.ErrInvalidCCBatch.Wrapf("frame %d is outdated: %d < %d", i, frame.Header.Timestamp, ccState.GetTimestamp())
		}
		if frame.Header.L1Height < ccState.GetL1Height() {
			return types.ErrInvalidCCBatch.Wrapf("frame %d has invalid l1 height %d", i, frame.Header.L1Height)
		}
		if frame.Header.GetL2Height() != ccState.ProcessedL2Block+1 {
			return types.ErrInvalidCCBatch.Wrapf("frame %d has invalid l2 height %d, expected %d", i, frame.Header.GetL2Height(), ccState.ProcessedL2Block+1)
		}

		for j, elem := range frame.Elements {
			if elem.Txraw == nil {
				if elem.QueueIndex < 1 || elem.QueueIndex != ccState.NextQueueIndex {
					return types.ErrInvalidCCBatch.Wrapf("frame %d element %d queue index %d, expected %d", i, j, elem.QueueIndex, ccState.NextQueueIndex)
				}
				// TODO: check the existence of queue tx
				ccState.NextQueueIndex++
			}
		}

		ccState.ProcessedL2Block++
		if i == (len(batch.Frames) - 1) {
			ccState.Timestamp = frame.Header.Timestamp
			ccState.L1Height = frame.Header.L1Height
		}
	}

	ccState.Height++
	k.setCCRef(ctx, rollupName, ccState.Height, ref)
	k.setCCState(ctx, rollupName, ccState)

	if err := ctx.EventManager().EmitTypedEvent(&types.EventUpdateCCBatch{
		RollupName:       rollupName,
		BatchIndex:       ccState.Height,
		NextQueueIndex:   ccState.NextQueueIndex,
		TotalFrames:      ref.TotalFrames,
		BatchSize:        ref.BatchSize,
		BatchHash:        ref.BatchRoot,
		ProcessedL2Block: ccState.ProcessedL2Block,
	}); err != nil {
		return err
	}

	return nil
}

func (k Keeper) DecompressCCBatch(ctx sdktypes.Context, origin types.CompressedCCBatch) (*types.CCBatch, error) {
	p, err := k.GetParams(ctx)
	if err != nil {
		return nil, err
	}
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
		buf := make([]byte, p.CCBatchMaxBytes)
		n, err := r.Read(buf)
		out := buf[:n]
		if err != nil && err != io.EOF {
			return nil, err
		}
		batch := new(types.CCBatch)
		k.cdc.MustUnmarshal(out, batch)

		return batch, nil

	case types.OptionZSTD:
		return nil, types.ErrInvalidCompressedData.Wrapf("compression %s not supported", origin.Compression)
	case types.OptionEmpty:
		return nil, types.ErrInvalidCompressedData
	default:
		return nil, sdkerror.ErrInvalidRequest.Wrapf("no compression option provided")
	}
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
