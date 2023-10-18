package v2_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	simappparams "github.com/Finschia/finschia-sdk/simapp/params"
	"github.com/Finschia/finschia-sdk/testutil"
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/collection"

	"github.com/Finschia/finschia-sdk/x/collection/keeper/migrations/v2"
)

func TestMigrateStore(t *testing.T) {
	collectionKey := sdk.NewKVStoreKey(collection.StoreKey)
	newKey := sdk.NewTransientStoreKey("transient_test")
	encCfg := simappparams.MakeTestEncodingConfig()
	ctx := testutil.DefaultContext(collectionKey, newKey)

	// set state
	store := ctx.KVStore(collectionKey)

	contractID := "deadbeef"
	store.Set(v2.ContractKey(contractID), encCfg.Marshaler.MustMarshal(&collection.Contract{Id: contractID}))
	nextClassIDs := collection.DefaultNextClassIDs(contractID)
	classID := fmt.Sprintf("%08x", nextClassIDs.Fungible.Uint64())
	nextClassIDs.Fungible = nextClassIDs.Fungible.Incr()
	store.Set(v2.NextClassIDKey(contractID), encCfg.Marshaler.MustMarshal(&nextClassIDs))

	tokenID := collection.NewFTID(classID)
	oneIntBz, err := sdk.OneInt().Marshal()
	require.NoError(t, err)
	addresses := []sdk.AccAddress{
		sdk.AccAddress("fennec"),
		sdk.AccAddress("penguin"),
		sdk.AccAddress("cheetah"),
	}
	for _, addr := range addresses {
		store.Set(v2.BalanceKey(contractID, addr, tokenID), oneIntBz)
	}
	store.Set(v2.StatisticKey(v2.SupplyKeyPrefix, contractID, classID), oneIntBz)
	store.Set(v2.StatisticKey(v2.MintedKeyPrefix, contractID, classID), oneIntBz)
	store.Set(v2.StatisticKey(v2.BurntKeyPrefix, contractID, classID), oneIntBz)

	for name, tc := range map[string]struct {
		malleate func(ctx sdk.Context)
		valid    bool
		supply   int
		minted   int
	}{
		"valid": {
			valid:  true,
			supply: len(addresses),
			minted: len(addresses) + 1,
		},
		"valid (nil supply)": {
			malleate: func(ctx sdk.Context) {
				store := ctx.KVStore(collectionKey)
				store.Delete(v2.StatisticKey(v2.SupplyKeyPrefix, contractID, classID))
			},
			valid:  true,
			supply: len(addresses),
			minted: len(addresses) + 1,
		},
		"valid (nil minted)": {
			malleate: func(ctx sdk.Context) {
				store := ctx.KVStore(collectionKey)
				store.Delete(v2.StatisticKey(v2.MintedKeyPrefix, contractID, classID))
			},
			valid:  true,
			supply: len(addresses),
			minted: len(addresses) + 1,
		},
		"valid (nil burnt)": {
			malleate: func(ctx sdk.Context) {
				store := ctx.KVStore(collectionKey)
				store.Delete(v2.StatisticKey(v2.BurntKeyPrefix, contractID, classID))
			},
			valid:  true,
			supply: len(addresses),
			minted: len(addresses),
		},
		"contract unmarshal failed": {
			malleate: func(ctx sdk.Context) {
				store := ctx.KVStore(collectionKey)
				store.Set(v2.ContractKey(contractID), encCfg.Marshaler.MustMarshal(&collection.GenesisState{}))
			},
		},
		"balance unmarshal failed": {
			malleate: func(ctx sdk.Context) {
				store := ctx.KVStore(collectionKey)
				store.Set(v2.BalanceKey(contractID, sdk.AccAddress("hyena"), tokenID), encCfg.Marshaler.MustMarshal(&collection.GenesisState{}))
			},
		},
		"no next class id": {
			malleate: func(ctx sdk.Context) {
				store := ctx.KVStore(collectionKey)
				store.Delete(v2.NextClassIDKey(contractID))
			},
		},
		"next class id unmarshal failed": {
			malleate: func(ctx sdk.Context) {
				store := ctx.KVStore(collectionKey)
				store.Set(v2.NextClassIDKey(contractID), []byte("invalid"))
			},
		},
		"burnt unmarshal failed": {
			malleate: func(ctx sdk.Context) {
				store := ctx.KVStore(collectionKey)
				store.Set(v2.StatisticKey(v2.BurntKeyPrefix, contractID, classID), encCfg.Marshaler.MustMarshal(&collection.GenesisState{}))
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			ctx, _ := ctx.CacheContext()
			if tc.malleate != nil {
				tc.malleate(ctx)
			}

			// migrate
			err := v2.MigrateStore(ctx, collectionKey, encCfg.Marshaler)
			if !tc.valid {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			store := ctx.KVStore(collectionKey)

			// supply
			supplyKey := v2.StatisticKey(v2.SupplyKeyPrefix, contractID, classID)
			supply := sdk.ZeroInt()
			if bz := store.Get(supplyKey); bz != nil {
				err := supply.Unmarshal(bz)
				require.NoError(t, err)
			}
			require.Equal(t, int64(tc.supply), supply.Int64())

			// minted
			mintedKey := v2.StatisticKey(v2.MintedKeyPrefix, contractID, classID)
			minted := sdk.ZeroInt()
			if bz := store.Get(mintedKey); bz != nil {
				err := minted.Unmarshal(bz)
				require.NoError(t, err)
			}
			require.Equal(t, int64(tc.minted), minted.Int64())
		})
	}
}
