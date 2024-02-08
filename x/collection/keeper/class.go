package keeper

import (
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"math"

	cmath "cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	nonceKey    = []byte{0x01}
	idKeyPrefix = []byte{0x02}
)

func idKey(id string) []byte {
	key := make([]byte, len(idKeyPrefix)+len(id))
	copy(key, idKeyPrefix)
	copy(key[len(idKeyPrefix):], id)
	return key
}

func splitIDKey(key []byte) (id string) {
	return string(key[len(idKeyPrefix):])
}

// NewID returns a brand-new ID.
func (k Keeper) NewID(ctx sdk.Context) string {
	for nonce := k.getNonce(ctx); nonce.LTE(cmath.NewUint(math.MaxUint64)); nonce = nonce.Incr() {
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

func (k Keeper) getNonce(ctx sdk.Context) cmath.Uint {
	classStore := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), classStorePrefix)
	bz := classStore.Get(nonceKey)
	if bz == nil {
		panic("next id must exist")
	}
	var nonce cmath.Uint
	if err := nonce.Unmarshal(bz); err != nil {
		panic(err)
	}
	return nonce
}

func (k Keeper) setNonce(ctx sdk.Context, nonce cmath.Uint) {
	classStore := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), classStorePrefix)
	bz, err := nonce.Marshal()
	if err != nil {
		panic(err)
	}
	classStore.Set(nonceKey, bz)
}

func (k Keeper) addID(ctx sdk.Context, id string) {
	classStore := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), classStorePrefix)
	classStore.Set(idKey(id), []byte{})
}

func (k Keeper) HasID(ctx sdk.Context, id string) bool {
	classStore := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), classStorePrefix)
	return classStore.Has(idKey(id))
}

func (k Keeper) DeleteID(ctx sdk.Context, id string) {
	classStore := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), classStorePrefix)
	classStore.Delete(idKey(id))
}
