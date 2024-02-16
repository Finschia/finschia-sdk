package v3

import (
	gogotypes "github.com/cosmos/gogoproto/types"

	cmath "cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/Finschia/finschia-sdk/x/collection"
)

func MigrateStore(collectionStore, oldClassStore storetypes.KVStore, cdc codec.BinaryCodec) error {
	err := migrateClassStateAndRemoveLegacyState(collectionStore, oldClassStore, cdc)
	if err != nil {
		return err
	}

	err = initParams(collectionStore)
	if err != nil {
		return err
	}

	return nil
}

func migrateClassStateAndRemoveLegacyState(collectionStore, oldStore storetypes.KVStore, cdc codec.BinaryCodec) error {
	newClassStore := prefix.NewStore(collectionStore, classStorePrefix)
	iterator := storetypes.KVStorePrefixIterator(oldStore, idKeyPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		id := splitIDKey(iterator.Key())
		newClassStore.Set(idKey(id), []byte{})
		oldStore.Delete(idKey(id))

		detachNFT(collectionStore, cdc, id)
	}

	bz := oldStore.Get(nonceKey)
	if bz == nil {
		return sdkerrors.ErrNotFound.Wrap("next id must exist")
	}
	var nonce cmath.Uint
	if err := nonce.Unmarshal(bz); err != nil {
		return err
	}
	newClassStore.Set(nonceKey, bz)
	oldStore.Set(nonceKey, bz)

	return nil
}

func detachNFT(store storetypes.KVStore, cdc codec.BinaryCodec, id string) {
	iterator := storetypes.KVStorePrefixIterator(store, parentKeyPrefixByContractID(id))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		contractID, tokenID := splitParentKey(iterator.Key())
		if contractID != id {
			panic("invalid NFT attachment information")
		}

		var parentID gogotypes.StringValue
		cdc.MustUnmarshal(iterator.Value(), &parentID)
		addCoin(store, id, getRootOwner(store, cdc, contractID, tokenID), collection.NewCoin(tokenID, cmath.OneInt()))

		store.Delete(parentKey(contractID, tokenID))
		store.Delete(childKey(contractID, parentID.Value, tokenID))
		store.Delete(key)
	}
}

func initParams(store storetypes.KVStore) error {
	p := collection.Params{}
	bz, err := p.Marshal()
	if err != nil {
		return err
	}
	store.Set(paramsKey, bz)
	return nil
}
