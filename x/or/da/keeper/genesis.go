package keeper

import (
	"encoding/binary"

	sdktypes "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/or/da/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func (k Keeper) InitGenesis(ctx sdktypes.Context, genState types.GenesisState) {
	if err := k.SetParams(ctx, genState.Params); err != nil {
		panic(err)
	}

	for _, cc := range genState.CCList {
		cc := cc
		name := cc.RollupName
		k.setCCState(ctx, name, &cc.CCState)
		k.setQueueTxState(ctx, name, &cc.QueueTxState)

		for i, ref := range cc.History {
			ref := ref
			k.setCCRef(ctx, name, uint64(i+1), &ref)
		}

		for i, tx := range cc.QueueList {
			tx := tx
			k.setQueueTx(ctx, name, uint64(i+1), &tx)
		}

		for _, bMap := range cc.L2BatchMap {
			k.setL2HeightBatchMap(ctx, name, bMap.L2Height, bMap.BatchIdx)
		}
	}
}

// ExportGenesis returns the module's exported genesis
func (k Keeper) ExportGenesis(ctx sdktypes.Context) *types.GenesisState {
	return &types.GenesisState{
		Params:  k.GetParams(ctx),
		CCList:  k.GetAllCCs(ctx),
		SCCList: k.GetAllSCCs(ctx),
	}
}

func (k Keeper) GetAllCCs(ctx sdktypes.Context) []types.CC {
	m := make(map[string]types.CC)
	ccStates := k.GetAllCCStates(ctx)
	ccRefs := k.GetAllCCRefs(ctx)
	qStates := k.GetAllQueueTxStates(ctx)
	qTxs := k.GetAllQueueTxs(ctx)
	bMaps := k.GetAllL2BatchMaps(ctx)

	ccs := make([]types.CC, len(m))
	for k, v := range ccStates {
		cc := types.CC{}
		cc.CCState = v
		cc.History = ccRefs[k]
		cc.QueueTxState = qStates[k]
		cc.QueueList = qTxs[k]
		cc.L2BatchMap = bMaps[k]
		cc.RollupName = k

		ccs = append(ccs, cc)
	}

	return ccs
}

func (k Keeper) GetAllSCCs(ctx sdktypes.Context) []types.SCC {
	m := make(map[string]types.SCC)
	sccStates := k.GetAllSCCStates(ctx)
	sccRefs := k.GetAllSCCRefs(ctx)

	sccs := make([]types.SCC, len(m))
	for k, v := range sccStates {
		scc := types.SCC{}
		scc.State = v
		scc.History = sccRefs[k]
		scc.RollupName = k

		sccs = append(sccs, scc)
	}

	return sccs
}

func (k Keeper) GetAllCCRefs(ctx sdktypes.Context) map[string][]types.CCRef {
	m := make(map[string][]types.CCRef)
	k.iterateAllCCRefs(ctx, []byte{types.CCBatchIndexPrefix}, func(rollupName string, ccRef types.CCRef) bool {
		m[rollupName] = append(m[rollupName], ccRef)
		return false
	})
	return m
}

func (k Keeper) iterateAllCCRefs(ctx sdktypes.Context, prefix []byte, cb func(rollupName string, ccRef types.CCRef) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdktypes.KVStorePrefixIterator(store, prefix)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var ccRef types.CCRef
		k.cdc.MustUnmarshal(iter.Value(), &ccRef)
		_, name, _ := types.SplitPrefixIndexKey(iter.Key())

		if cb(name, ccRef) {
			break
		}
	}
}

func (k Keeper) GetAllCCStates(ctx sdktypes.Context) map[string]types.CCState {
	m := make(map[string]types.CCState)
	k.iterateAllCCStates(ctx, []byte{types.CCStateStoreKey}, func(rollupName string, ccState types.CCState) bool {
		_, ok := m[rollupName]
		if ok {
			panic("duplicate cc state")
		}
		m[rollupName] = ccState
		return false
	})
	return m
}

func (k Keeper) iterateAllCCStates(ctx sdktypes.Context, prefix []byte, cb func(rollupName string, ccs types.CCState) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdktypes.KVStorePrefixIterator(store, prefix)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var ccs types.CCState
		k.cdc.MustUnmarshal(iter.Value(), &ccs)
		_, name := types.SplitPrefixKey(iter.Key())

		if cb(name, ccs) {
			break
		}
	}
}

func (k Keeper) GetAllQueueTxStates(ctx sdktypes.Context) map[string]types.QueueTxState {
	m := make(map[string]types.QueueTxState)
	k.iterateAllQueueTxState(ctx, []byte{types.QueueTxStateStoreKey}, func(rollupName string, txState types.QueueTxState) bool {
		_, ok := m[rollupName]
		if ok {
			panic("duplicate queue tx state")
		}
		m[rollupName] = txState
		return false
	})
	return m
}

func (k Keeper) iterateAllQueueTxState(ctx sdktypes.Context, prefix []byte, cb func(rollupName string, txState types.QueueTxState) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdktypes.KVStorePrefixIterator(store, prefix)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var txState types.QueueTxState
		k.cdc.MustUnmarshal(iter.Value(), &txState)
		_, name := types.SplitPrefixKey(iter.Key())

		if cb(name, txState) {
			break
		}
	}
}

func (k Keeper) GetAllQueueTxs(ctx sdktypes.Context) map[string][]types.L1ToL2Queue {
	m := make(map[string][]types.L1ToL2Queue)
	k.iterateAllQueueTxs(ctx, []byte{types.CCQueueTxPrefix}, func(rollupName string, tx types.L1ToL2Queue) bool {
		m[rollupName] = append(m[rollupName], tx)
		return false
	})
	return m
}

func (k Keeper) iterateAllQueueTxs(ctx sdktypes.Context, prefix []byte, cb func(rollupName string, tx types.L1ToL2Queue) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdktypes.KVStorePrefixIterator(store, prefix)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var tx types.L1ToL2Queue
		k.cdc.MustUnmarshal(iter.Value(), &tx)
		_, name, _ := types.SplitPrefixIndexKey(iter.Key())

		if cb(name, tx) {
			break
		}
	}
}

func (k Keeper) GetAllL2BatchMaps(ctx sdktypes.Context) map[string][]types.L2BatchMap {
	m := make(map[string][]types.L2BatchMap)
	k.iterateAllL2BatchMaps(ctx, []byte{types.CCL2HeightToBatchPrefix}, func(rollupName string, batchMap types.L2BatchMap) bool {
		m[rollupName] = append(m[rollupName], batchMap)
		return false
	})
	return m
}

func (k Keeper) iterateAllL2BatchMaps(ctx sdktypes.Context, prefix []byte, cb func(rollupName string, batchMap types.L2BatchMap) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdktypes.KVStorePrefixIterator(store, prefix)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var batchMap types.L2BatchMap
		_, name, h := types.SplitPrefixIndexKey(iter.Key())
		batchMap.BatchIdx = binary.BigEndian.Uint64(iter.Value())
		batchMap.L2Height = h

		if cb(name, batchMap) {
			break
		}
	}
}

func (k Keeper) GetAllSCCRefs(ctx sdktypes.Context) map[string][]types.SCCRef {
	m := make(map[string][]types.SCCRef)
	k.iterateAllSCCRefs(ctx, []byte{types.SCCBatchIndexPrefix}, func(rollupName string, sccRef types.SCCRef) bool {
		m[rollupName] = append(m[rollupName], sccRef)
		return false
	})
	return m
}

func (k Keeper) iterateAllSCCRefs(ctx sdktypes.Context, prefix []byte, cb func(rollupName string, sccRef types.SCCRef) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdktypes.KVStorePrefixIterator(store, prefix)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var sccRef types.SCCRef
		k.cdc.MustUnmarshal(iter.Value(), &sccRef)
		_, name, _ := types.SplitPrefixIndexKey(iter.Key())

		if cb(name, sccRef) {
			break
		}
	}
}

func (k Keeper) GetAllSCCStates(ctx sdktypes.Context) map[string]types.SCCState {
	m := make(map[string]types.SCCState)
	k.iterateAllSCCStates(ctx, []byte{types.SCCStateStoreKey}, func(rollupName string, sccState types.SCCState) bool {
		_, ok := m[rollupName]
		if ok {
			panic("duplicate scc state")
		}
		m[rollupName] = sccState
		return false
	})
	return m
}

func (k Keeper) iterateAllSCCStates(ctx sdktypes.Context, prefix []byte, cb func(rollupName string, sccState types.SCCState) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdktypes.KVStorePrefixIterator(store, prefix)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var sccState types.SCCState
		k.cdc.MustUnmarshal(iter.Value(), &sccState)
		_, name := types.SplitPrefixKey(iter.Key())

		if cb(name, sccState) {
			break
		}
	}
}
