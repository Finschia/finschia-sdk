package keeper

import (
	"bytes"
	"crypto/sha256"

	sdktypes "github.com/Finschia/finschia-sdk/types"
	sdkerrors "github.com/Finschia/finschia-sdk/types/errors"
	"github.com/Finschia/finschia-sdk/x/or/da/types"
	rutypes "github.com/Finschia/finschia-sdk/x/or/rollup/types"
)

// SaveSC
func (k Keeper) SaveSCCBatch(ctx sdktypes.Context, from string, rollupName string, batch *types.SCCBatch) error {
	var sccState *types.SCCState
	sccState, _ = k.GetSCCState(ctx, rollupName)
	if sccState == nil {
		if batch.ShouldStartAtFrame != 0 {
			return types.ErrInvalidSCCBatch.Wrapf("scc state not found: this batch should start at frame 0, but start at %d", batch.ShouldStartAtFrame)
		}

		sccState = &types.SCCState{
			Base:   1,
			Height: 0,
		}
	} else {
		seqInfos, found := k.rollupKeeper.GetSequencersByRollupName(ctx, rollupName)
		if !found {
			return sdkerrors.ErrNotFound.Wrapf("sequencers not found for rollup %s", rollupName)
		}
		limit := sccState.LastSequencerSubmit.Add(k.SequencerPublishWindow(ctx))
		if limit.After(ctx.BlockTime().UTC()) && !containSequencer(seqInfos.Sequencers, from) {
			return sdkerrors.ErrUnauthorized.Wrapf("non-sequencer can only submit after %s", limit)
		}
	}

	sccState.Height++
	sccState.LastSequencerSubmit = ctx.BlockTime().UTC()

	ccRef, err := k.GetCCRef(ctx, rollupName, sccState.Height)
	if err != nil {
		panic(err)
	}

	if len(batch.IntermediateStateRoots) < 2 {
		return types.ErrInvalidSCCBatch.Wrapf("Intermediate State Roots length must be at least 2(pre-state and post-state), got %d", len(batch.IntermediateStateRoots))
	} else if len(batch.IntermediateStateRoots)-1 != int(ccRef.BatchSize) {
		return types.ErrInvalidSCCBatch.Wrapf("Intermediate State Roots length must be equal to target cc batch size, expected %d got %d", ccRef.BatchSize, len(batch.IntermediateStateRoots))
	}

	var totalFrames uint64
	if sccState.Height-1 == 0 {
		totalFrames = uint64(ccRef.BatchSize)
	} else {
		prefSCCRef, err := k.GetSCCRef(ctx, rollupName, sccState.Height-1)
		if err != nil {
			return err
		}

		if prefSCCRef.TotalFrames != batch.ShouldStartAtFrame {
			return types.ErrInvalidSCCBatch.Wrapf("batch should start at frame %d but this batch start at frame %d", prefSCCRef.TotalFrames, batch.ShouldStartAtFrame)
		}

		if bytes.Equal(prefSCCRef.IntermediateStateRoots[len(prefSCCRef.IntermediateStateRoots)-1], batch.IntermediateStateRoots[0]) {
			return types.ErrInvalidSCCBatch.Wrapf("SCC batch pre-state root not match with previous SCC batch post-state root")
		}

		totalFrames = prefSCCRef.TotalFrames + uint64(ccRef.BatchSize)
	}

	rootHash := computeMerkleRoot(batch.IntermediateStateRoots)
	ref := types.NewSCCRef(totalFrames, ccRef.BatchSize, rootHash, batch.IntermediateStateRoots, ctx.BlockTime().UTC())

	k.setSCCRef(ctx, rollupName, sccState.Height, ref)
	k.setSCCState(ctx, rollupName, sccState)

	if err := ctx.EventManager().EmitTypedEvent(&types.EventAppendSCCBatch{
		RollupName:          rollupName,
		BatchIndex:          sccState.Height,
		TotalFrames:         ref.TotalFrames,
		BatchSize:           ref.BatchSize,
		BatchRoot:           ref.BatchRoot,
		LastSequencerSubmit: ctx.BlockTime().UTC(),
	}); err != nil {
		return err
	}

	return nil
}

// DeleteSCCBatch deletes a scc batch from specific batch index.
// This is because each state root is continuous,
// so if the state root at a certain point in time is wrong, all subsequent values are wrong.
// Caution: This function must only be executed by Rollup Settlement module.
func (k Keeper) DeleteSCCBatch(ctx sdktypes.Context, rollupName string, batchIndex uint64) error {
	sccState, err := k.GetSCCState(ctx, rollupName)
	if err != nil {
		return err
	}

	for i := batchIndex; i <= sccState.Height; i++ {
		if err := k.deleteSCCRef(ctx, rollupName, i); err != nil {
			panic("SCC State does not match the SCC data")
		}
	}

	sccState.Height = batchIndex - 1
	k.setSCCState(ctx, rollupName, sccState)

	if err := ctx.EventManager().EmitTypedEvent(&types.EventDeleteSCCBatch{
		RollupName:    rollupName,
		NewBatchIndex: sccState.Height,
	}); err != nil {
		return err
	}

	return nil
}

func (k Keeper) GetSCCRef(ctx sdktypes.Context, rollupName string, idx uint64) (*types.SCCRef, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetSCCBatchIndexKey(rollupName, idx))
	if bz == nil {
		return nil, types.ErrSCCRefNotFound
	}
	ref := new(types.SCCRef)
	k.cdc.MustUnmarshal(bz, ref)
	return ref, nil
}

func (k Keeper) setSCCRef(ctx sdktypes.Context, rollupName string, idx uint64, ref *types.SCCRef) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(ref)
	store.Set(types.GetSCCBatchIndexKey(rollupName, idx), bz)
}

func (k Keeper) deleteSCCRef(ctx sdktypes.Context, rollupName string, idx uint64) error {
	store := ctx.KVStore(k.storeKey)
	if store.Has(types.GetSCCBatchIndexKey(rollupName, idx)) {
		store.Delete(types.GetSCCBatchIndexKey(rollupName, idx))
		return nil
	}
	return types.ErrSCCRefNotFound
}

func (k Keeper) GetSCCState(ctx sdktypes.Context, rollupName string) (*types.SCCState, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetSCCStateStoreKey(rollupName))
	if bz == nil {
		return nil, types.ErrSCCStateNotFound
	}
	state := new(types.SCCState)
	k.cdc.MustUnmarshal(bz, state)
	return state, nil
}

func (k Keeper) setSCCState(ctx sdktypes.Context, rollupName string, state *types.SCCState) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(state)
	store.Set(types.GetSCCStateStoreKey(rollupName), bz)
}

func hash(data []byte) []byte {
	h := sha256.Sum256(data)
	return h[:]
}

func computeMerkleRoot(hashes [][]byte) []byte {
	if len(hashes) == 0 {
		return nil
	}
	if len(hashes) == 1 {
		return hashes[0]
	}
	if len(hashes)%2 != 0 {
		hashes = append(hashes, hashes[len(hashes)-1])
	}
	var nextLevel [][]byte
	for i := 0; i < len(hashes); i += 2 {
		nextLevel = append(nextLevel, hash(append(hashes[i], hashes[i+1]...)))
	}
	return computeMerkleRoot(nextLevel)
}

func containSequencer(sequencers []rutypes.Sequencer, address string) bool {
	for _, seq := range sequencers {
		if seq.SequencerAddress == address {
			return true
		}
	}
	return false
}
