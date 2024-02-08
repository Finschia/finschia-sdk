package v3

import (
	cmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/Finschia/finschia-sdk/x/collection"
)

// TODO: need to implement logics for migrating ClassState and generate new genesis state for Collection
func MigrateStore(store, classStore storetypes.KVStore, cdc codec.BinaryCodec) error {
	_, err := getClassState(classStore)
	if err != nil {
		return err
	}

	_ = store
	_ = cdc

	return nil
}

func getClassState(store storetypes.KVStore) (collection.ClassState, error) {
	cs := collection.ClassState{}

	ids := make([]string, 0)
	iterateIDs(store, func(id string) (stop bool) {
		ids = append(ids, id)
		return false
	})
	cs.Ids = ids

	bz := store.Get(nonceKey)
	if bz == nil {
		return cs, sdkerrors.ErrNotFound.Wrap("next id must exist")
	}
	var nonce cmath.Uint
	if err := nonce.Unmarshal(bz); err != nil {
		return cs, err
	}
	cs.Nonce = nonce

	return cs, nil
}

func iterateIDs(store storetypes.KVStore, fn func(id string) (stop bool)) {
	iterator := storetypes.KVStorePrefixIterator(store, idKeyPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		id := splitIDKey(iterator.Key())

		stop := fn(id)
		if stop {
			break
		}
	}
}
