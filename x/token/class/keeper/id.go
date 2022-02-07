package keeper

import (
	"encoding/binary"
	"fmt"
	"hash/fnv"

	sdk "github.com/line/lbm-sdk/types"
)

// NewId returns a brand new ID.
func (k Keeper) NewId(ctx sdk.Context) string {
	for nextId := k.getNextId(ctx); nextId.IsUint64(); nextId = nextId.Add(sdk.OneInt()) {
		encoded := encodeId(nextId.Uint64())
		if !k.HasId(ctx, encoded) {
			k.addId(ctx, encoded)
			k.setNextId(ctx, nextId.Add(sdk.OneInt()))
			return encoded
		}
	}
	panic("Class id space exhausted: uint64")
}

func encodeId(id uint64) string {
	hash := fnv.New32()
	bz := make([]byte, 8)
	binary.LittleEndian.PutUint64(bz, id)
	_, err := hash.Write(bz)
	if err != nil {
		panic("hash should not fail")
	}
	idStr := fmt.Sprintf("%x", hash.Sum32())
	if len(idStr) < 8 {
		idStr = "00000000"[len(idStr):] + idStr
	}
	return idStr
}

func (k Keeper) getNextId(ctx sdk.Context) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(nextIdKey)
	if bz == nil {
		panic("next id must exist")
	}
	var nextId sdk.Int
	if err := nextId.Unmarshal(bz); err != nil {
		panic(err)
	}
	return nextId
}

func (k Keeper) setNextId(ctx sdk.Context, id sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	bz, err := id.Marshal()
	if err != nil {
		panic(err)
	}
	store.Set(nextIdKey, bz)
}

func (k Keeper) addId(ctx sdk.Context, id string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(idKey(id), []byte{})
}

func (k Keeper) HasId(ctx sdk.Context, id string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(idKey(id))
}

func (k Keeper) DeleteId(ctx sdk.Context, id string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(idKey(id))
}
