package v2

import (
	"fmt"

	"github.com/Finschia/finschia-sdk/codec"
	storetypes "github.com/Finschia/finschia-sdk/store/types"
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/collection"
)

// MigrateStore performs in-place store migrations from v1 to v2.
func MigrateStore(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.BinaryCodec) error {
	store := ctx.KVStore(storeKey)

	// fix ft statistics
	if err := fixFTStatistics(store, cdc); err != nil {
		return err
	}

	return nil
}

func fixFTStatistics(store storetypes.KVStore, cdc codec.BinaryCodec) error {
	iterator := sdk.KVStorePrefixIterator(store, contractKeyPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var contract collection.Contract
		if err := cdc.Unmarshal(iterator.Value(), &contract); err != nil {
			return err
		}

		if err := fixContractFTStatistics(store, contract.Id); err != nil {
			return err
		}
	}

	return nil
}

func fixContractFTStatistics(store storetypes.KVStore, contractID string) error {
	supplies, err := evalContractFTSupplies(store, contractID)
	if err != nil {
		return err
	}

	if err := updateContractFTStatistics(store, contractID, supplies); err != nil {
		return err
	}

	return nil
}

func evalContractFTSupplies(store storetypes.KVStore, contractID string) (map[string]sdk.Int, error) {
	prefix := balanceKeyPrefixByContractID(contractID)
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	supplies := map[string]sdk.Int{}
	for ; iterator.Valid(); iterator.Next() {
		_, _, tokenID := splitBalanceKey(iterator.Key())
		if err := collection.ValidateFTID(tokenID); err != nil {
			continue
		}

		var amount sdk.Int
		if err := amount.Unmarshal(iterator.Value()); err != nil {
			return nil, err
		}

		classID := collection.SplitTokenID(tokenID)
		if supply, ok := supplies[classID]; ok {
			supplies[classID] = supply.Add(amount)
		} else {
			supplies[classID] = amount
		}
	}

	return supplies, nil
}

func updateContractFTStatistics(store storetypes.KVStore, contractID string, supplies map[string]sdk.Int) error {
	bz := store.Get(NextClassIDKey(contractID))
	if bz == nil {
		return fmt.Errorf("no next class ids of contract %s", contractID)
	}

	var nextClassIDs collection.NextClassIDs
	if err := nextClassIDs.Unmarshal(bz); err != nil {
		return err
	}

	// In the old chains, classID of fungible tokens starts from zero
	// In the new chains, it starts from one, but it does not hurts because amount of zero is not set to the store.
	for intClassID := uint64(0); intClassID < nextClassIDs.Fungible.Uint64(); intClassID++ {
		classID := fmt.Sprintf("%08x", intClassID)

		// update supply
		supplyKey := StatisticKey(SupplyKeyPrefix, contractID, classID)
		supply, ok := supplies[classID]
		if ok {
			bz, err := supply.Marshal()
			if err != nil {
				return err
			}
			store.Set(supplyKey, bz)
		} else {
			supply = sdk.ZeroInt()
			store.Delete(supplyKey)
		}

		// get burnt
		burntKey := StatisticKey(BurntKeyPrefix, contractID, classID)
		burnt := sdk.ZeroInt()
		if bz := store.Get(burntKey); bz != nil {
			if err := burnt.Unmarshal(bz); err != nil {
				return err
			}
		}

		// update minted
		minted := supply.Add(burnt)
		mintedKey := StatisticKey(MintedKeyPrefix, contractID, classID)
		if !minted.IsZero() {
			bz, err := minted.Marshal()
			if err != nil {
				return err
			}
			store.Set(mintedKey, bz)
		} else {
			store.Delete(mintedKey)
		}
	}

	return nil
}
