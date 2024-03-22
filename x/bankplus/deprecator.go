package bankplus

import (
	"context"

	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// inactiveAddrsKeyPrefix Keys for bankplus store but this prefix must not be overlap with bank key prefix.
var inactiveAddrsKeyPrefix = []byte{0xa0}

// inactiveAddrKey key of a specific inactiveAddr from store
func inactiveAddrKey(addr sdk.AccAddress) []byte {
	return append(inactiveAddrsKeyPrefix, addr.Bytes()...)
}

// DeprecateBankPlus performs remove logic for bankplus v1.
// This will remove all the state(inactive addresses)
// This supposed to be called in simapp.
//
// Example) simapp/upgrades.go
//
//	func (app SimApp) RegisterUpgradeHandlers() {
//		app.UpgradeKeeper.SetUpgradeHandler(
//			UpgradeName,
//			func(ctx context.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
//				app.deprecateBankPlusFromSimapp(ctx)
//				return app.ModuleManager.RunMigrations(ctx, app.Configurator(), fromVM)
//			},
//		)
//		...
//
//	func (app SimApp) deprecateBankPlusFromSimapp(ctx context.Context) {
//		for _, key := range app.kvStoreKeys() {
//			if key.Name() == banktypes.StoreKey {
//				err := internal.DeprecateBankPlus(ctx, key)
//				if err != nil {
//					panic(fmt.Errorf("failed to deprecate x/bankplus: %w", err))
//				}
//			}
//		}
//	}
func DeprecateBankPlus(ctx context.Context, bankKey *storetypes.KVStoreKey) error {
	kss := runtime.NewKVStoreService(bankKey)
	ks := kss.OpenKVStore(ctx)
	adapter := runtime.KVStoreAdapter(ks)
	iter := storetypes.KVStorePrefixIterator(adapter, inactiveAddrsKeyPrefix)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		err := ks.Delete(iter.Key())
		if err != nil {
			return err
		}
	}
	return nil
}
