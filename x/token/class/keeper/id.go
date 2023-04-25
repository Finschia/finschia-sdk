package keeper

import (
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"math"

	sdk "github.com/Finschia/finschia-sdk/types"
)

// NewID returns a brand new ID.
func (k Keeper) NewID(ctx sdk.Context) string {
	for nonce := k.getNonce(ctx); nonce.LTE(sdk.NewUint(math.MaxUint64)); nonce = nonce.Incr() {
		encoded := nonceToID(nonce.Uint64())
		if !k.HasID(ctx, encoded) {
			k.addID(ctx, encoded)
			k.setNonce(ctx, nonce.Incr())
			return encoded
		}
	}
	panic("contract id space exhausted: uint64")
}

func nonceToID(nonce uint64) string {
	hash := fnv.New32()
	bz := make([]byte, 8)
	binary.LittleEndian.PutUint64(bz, nonce)
	_, err := hash.Write(bz)
	if err != nil {
		panic("hash should not fail")
	}
	id := fmt.Sprintf("%x", hash.Sum32())
	if len(id) < 8 {
		id = "00000000"[len(id):] + id
	}
	return id
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

func (k Keeper) addID(ctx sdk.Context, id string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(idKey(id), []byte{})
}

func (k Keeper) HasID(ctx sdk.Context, id string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(idKey(id))
}

func (k Keeper) DeleteID(ctx sdk.Context, id string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(idKey(id))
}
