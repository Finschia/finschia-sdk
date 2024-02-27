package v3

import (
	"fmt"

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
		id := splitIDKey(iterator.Key()) // contract ID
		newClassStore.Set(idKey(id), []byte{})
		oldStore.Delete(idKey(id))

		detachNFTs(collectionStore, cdc, id)
		removeFTs(collectionStore, cdc, id)
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

func detachNFTs(store storetypes.KVStore, cdc codec.BinaryCodec, id string) {
	iterator := storetypes.KVStorePrefixIterator(store, parentKeyPrefixByContractID(id))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		contractID, tokenID := splitParentKey(iterator.Key())
		if contractID != id {
			panic("invalid NFT attachment information")
		}

		var parentID gogotypes.StringValue
		cdc.MustUnmarshal(iterator.Value(), &parentID)
		addCoin(store, id, getRootOwner(store, cdc, contractID, parentID.Value), collection.NewCoin(tokenID, cmath.OneInt()))

		store.Delete(parentKey(contractID, tokenID))
		store.Delete(childKey(contractID, parentID.Value, tokenID))
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

func removeFTs(store storetypes.KVStore, cdc codec.BinaryCodec, contractID string) {
	iterator := storetypes.KVStorePrefixIterator(store, classKeyPrefixByContractID(contractID))
	defer iterator.Close()

	ftMap := make(map[string]bool)
	for ; iterator.Valid(); iterator.Next() {
		bz := store.Get(iterator.Key())
		var class collection.TokenClass
		err := cdc.UnmarshalInterface(bz, &class)
		if err != nil {
			panic(err)
		}

		key := legacyTokenTypeKey(contractID, class.GetId())
		store.Delete(key)

		if ftClass, ok := class.(*collection.FTClass); ok {
			ftID := newFTID(ftClass.Id)
			if !validateFTID(ftID) {
				panic(fmt.Sprintf("v3 migration: invalid FT ID from TokenClass <%s>", ftID))
			}

			ftMap[ftID] = true
			store.Delete(legacyTokenKey(contractID, ftID)) // remove LegacyToken
			store.Delete(iterator.Key())                   // remove TokenClass
		}
	}

	removeFTBalances(store, contractID, ftMap)
}

func removeFTBalances(store storetypes.KVStore, contractID string, ftMap map[string]bool) {
	iterator := storetypes.KVStorePrefixIterator(store, balanceKeyPrefixByContractID(contractID))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		id, address, tokenID := splitBalanceKey(iterator.Key())
		if id != contractID {
			panic(fmt.Sprintf("v3 migration: inconsistent ContractID, got: %s, expected: %s", id, contractID))
		}

		if ftMap[tokenID] {
			classID := tokenID[:lengthClassID]
			var amount cmath.Int
			if err := amount.Unmarshal(iterator.Value()); err != nil {
				panic(err)
			}

			// remove FT statistics
			store.Delete(statisticKey(supplyKeyPrefix, contractID, classID))
			store.Delete(statisticKey(mintedKeyPrefix, contractID, classID))
			store.Delete(statisticKey(burntKeyPrefix, contractID, classID))

			removeCoin(store, contractID, address, tokenID)
		}
	}
}
