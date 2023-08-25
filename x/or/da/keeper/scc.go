package keeper

import (
	"crypto/sha256"
	sdktypes "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/or/da/types"
)

func (k Keeper) SaveSCCBatch(ctx sdktypes.Context, rollupName string, batch *types.SCCBatch) error {
	return nil
}

func (k Keeper) DeleteSccBatch(ctx sdktypes.Context, rollupName string, batchIndex uint64) error {
	return nil
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
