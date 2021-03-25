package keeper

import (
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"regexp"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/line/link-modules/x/contract/internal/types"
)

/***

Actual generated contractID values from genesis

9be17165
678c146a
3336b76f
fee15a74
ca8bfd79
9636a07e
61e14383
2d8be688
f936898d
c4e12c92
...

*/

const (
	IDRegExprString = "[a-f0-9]{8}"
)

var (
	IDRegExpr = regexp.MustCompile(fmt.Sprintf("^%s$", IDRegExprString))
)

type ContractKeeper interface {
	NewContractID(ctx sdk.Context) string
	HasContractID(ctx sdk.Context, contractID string) bool
	DeleteContractID(ctx sdk.Context, contractID string)
}

func NewContractKeeper(cdc *codec.Codec, storeKey sdk.StoreKey) ContractKeeper {
	return BaseContractKeeper{
		storeKey: storeKey,
		cdc:      cdc,
	}
}

type BaseContractKeeper struct {
	storeKey sdk.StoreKey
	cdc      *codec.Codec
}

var _ ContractKeeper = (*BaseContractKeeper)(nil)

func VerifyContractID(contractID string) bool {
	return IDRegExpr.MatchString(contractID)
}

func (k BaseContractKeeper) NewContractID(ctx sdk.Context) string {
	store := ctx.KVStore(k.storeKey)

	nextCount := uint64(0)
	if store.Has(types.LastContractCountStoreKey()) {
		b := store.Get(types.LastContractCountStoreKey())
		k.cdc.MustUnmarshalBinaryBare(b, &nextCount)
		nextCount++
	}
	var id string
	hash := fnv.New32()
	for ok := false; !ok; {
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, nextCount)
		_, err := hash.Write(b)
		if err != nil {
			panic("hash should not fail")
		}
		id = fmt.Sprintf("%x", hash.Sum32())
		if len(id) < 8 {
			id = "00000000"[len(id):] + id
		}
		if store.Has(types.ContractIDStoreKey(id)) {
			nextCount++
		} else {
			ok = true
		}
	}

	store.Set(types.LastContractCountStoreKey(), k.cdc.MustMarshalBinaryBare(nextCount))
	store.Set(types.ContractIDStoreKey(id), k.cdc.MustMarshalBinaryBare(nextCount))
	return id
}

func (k BaseContractKeeper) HasContractID(ctx sdk.Context, contractID string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.ContractIDStoreKey(contractID))
}

func (k BaseContractKeeper) DeleteContractID(ctx sdk.Context, contractID string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.ContractIDStoreKey(contractID))
}
