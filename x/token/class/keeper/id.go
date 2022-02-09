package keeper

import (
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"math"

	sdk "github.com/line/lbm-sdk/types"
)

// NewId returns a brand new ID.
func (k Keeper) NewId(ctx sdk.Context) string {
	for nonce := k.getNonce(ctx); nonce.LTE(sdk.NewUint(math.MaxUint64)); nonce = nonce.Incr() {
		encoded := nonceToId(nonce.Uint64())
		if !k.HasId(ctx, encoded) {
			k.addId(ctx, encoded)
			k.setNonce(ctx, nonce.Incr())
			return encoded
		}
	}
	panic("Class id space exhausted: uint64")
}

func nonceToId(id uint64) string {
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

func (k Keeper) getNonce(ctx sdk.Context) sdk.Uint {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(nonceKey)
	if bz == nil {
		panic("next id must exist")
	}
	var nonce sdk.Uint
	if err := nonce.Unmarshal(bz); err != nil {
		panic(err)
	}
	return nonce
}

func (k Keeper) setNonce(ctx sdk.Context, nonce sdk.Uint) {
	store := ctx.KVStore(k.storeKey)
	bz, err := nonce.Marshal()
	if err != nil {
		panic(err)
	}
	store.Set(nonceKey, bz)
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
