package v2

import (
	"fmt"

	"cosmossdk.io/core/store"
	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Finschia/finschia-sdk/x/collection"
)

// MigrateStore performs in-place store migrations from v1 to v2.
func MigrateStore(ctx sdk.Context, storeService store.KVStoreService, cdc codec.BinaryCodec) error {
	store := storeService.OpenKVStore(ctx)

	// fix ft statistics
	if err := fixFTStatistics(store, cdc); err != nil {
		return err
	}

	return nil
}

func fixFTStatistics(store store.KVStore, cdc codec.BinaryCodec) error {
	adapter := runtime.KVStoreAdapter(store)
	iterator := storetypes.KVStorePrefixIterator(adapter, contractKeyPrefix)
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

func fixContractFTStatistics(store store.KVStore, contractID string) error {
	supplies, err := evalContractFTSupplies(store, contractID)
	if err != nil {
		return err
	}

	if err := updateContractFTStatistics(store, contractID, supplies); err != nil {
		return err
	}

	return nil
}

func evalContractFTSupplies(store store.KVStore, contractID string) (map[string]math.Int, error) {
	prefix := balanceKeyPrefixByContractID(contractID)
	adapter := runtime.KVStoreAdapter(store)
	iterator := storetypes.KVStorePrefixIterator(adapter, prefix)
	defer iterator.Close()

	supplies := map[string]math.Int{}
	for ; iterator.Valid(); iterator.Next() {
		_, _, tokenID := splitBalanceKey(iterator.Key())
		if err := collection.ValidateFTID(tokenID); err != nil {
			continue
		}

		var amount math.Int
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

func updateContractFTStatistics(store store.KVStore, contractID string, supplies map[string]math.Int) error {
	bz, err := store.Get(NextClassIDKey(contractID))
	if err != nil {
		return err
	}
	if bz == nil {
		return fmt.Errorf("no next class ids of contract %s", contractID)
	}

	var nextClassIDs collection.NextClassIDs
	if err := nextClassIDs.Unmarshal(bz); err != nil {
		return err
	}

	// In the old chains, classID of fungible tokens starts from zero
	// In the new chains, it starts from one, but it does not hurt because amount of zero is not set to the store.
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
			if err = store.Set(supplyKey, bz); err != nil {
				return err
			}
		} else {
			supply = math.ZeroInt()
			if err := store.Delete(supplyKey); err != nil {
				return err
			}
		}

		// get burnt
		burntKey := StatisticKey(BurntKeyPrefix, contractID, classID)
		burnt := math.ZeroInt()
		bz, err := store.Get(burntKey)
		if err != nil {
			return err
		}
		if bz != nil {
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
			if err = store.Set(mintedKey, bz); err != nil {
				return err
			}
		} else {
			if err := store.Delete(mintedKey); err != nil {
				return err
			}
		}
	}

	return nil
}
