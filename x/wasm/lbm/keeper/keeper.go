package keeper

import (
	"github.com/line/lbm-sdk/codec"
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	authkeeper "github.com/line/lbm-sdk/x/auth/keeper"
	paramtypes "github.com/line/lbm-sdk/x/params/types"
	wasmkeeper "github.com/line/lbm-sdk/x/wasm/keeper"
	"github.com/line/lbm-sdk/x/wasm/types"
)

type Keeper struct {
	wasmkeeper.Keeper

	cdc      codec.Codec
	storeKey sdk.StoreKey
}

func NewKeeper(
	cdc codec.Codec,
	storeKey sdk.StoreKey,
	paramSpace paramtypes.Subspace,
	accountKeeper authkeeper.AccountKeeper,
	bankKeeper types.BankKeeper,
	stakingKeeper types.StakingKeeper,
	distKeeper types.DistributionKeeper,
	channelKeeper types.ChannelKeeper,
	portKeeper types.PortKeeper,
	capabilityKeeper types.CapabilityKeeper,
	portSource types.ICS20TransferPortSource,
	router wasmkeeper.MessageRouter,
	queryRouter wasmkeeper.GRPCQueryRouter,
	homeDir string,
	wasmConfig types.WasmConfig,
	supportedFeatures string,
	customEncoders *wasmkeeper.MessageEncoders,
	customPlugins *wasmkeeper.QueryPlugins,
	opts ...wasmkeeper.Option,
) Keeper {
	wasmKeeper := wasmkeeper.NewKeeper(cdc, storeKey, paramSpace, accountKeeper, bankKeeper, stakingKeeper, distKeeper,
		channelKeeper, portKeeper, capabilityKeeper, portSource, router, queryRouter, homeDir, wasmConfig,
		supportedFeatures, customEncoders, customPlugins, opts...)

	return Keeper{
		Keeper:   wasmKeeper,
		cdc:      cdc,
		storeKey: storeKey,
	}
}

func (k Keeper) IsInactiveContract(ctx sdk.Context, contractAddress sdk.AccAddress) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(getInactiveContractKey(contractAddress))
}

func (k Keeper) IterateInactiveContracts(ctx sdk.Context, fn func(contractAddress sdk.AccAddress) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	prefix := inactiveContractPrefix
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		contractAddress := sdk.AccAddress(iterator.Value())
		if stop := fn(contractAddress); stop {
			break
		}
	}
}

func (k Keeper) addInactiveContract(ctx sdk.Context, contractAddress sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	key := getInactiveContractKey(contractAddress)

	store.Set(key, contractAddress)
}

func (k Keeper) deleteInactiveContract(ctx sdk.Context, contractAddress sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	key := getInactiveContractKey(contractAddress)
	store.Delete(key)
}

func (k Keeper) activateContract(ctx sdk.Context, contractAddress sdk.AccAddress) error {
	if k.IsInactiveContract(ctx, contractAddress) {
		return sdkerrors.Wrap(types.ErrAccountExists, "inactivate contract")
	}

	k.addInactiveContract(ctx, contractAddress)

	return nil
}

func (k Keeper) deactivateContract(ctx sdk.Context, contractAddress sdk.AccAddress) error {
	if !k.IsInactiveContract(ctx, contractAddress) {
		return sdkerrors.Wrap(types.ErrNotFound, "deactivate contract")
	}

	k.deleteInactiveContract(ctx, contractAddress)

	return nil
}

func Querier(k *Keeper) *GrpcQuerier {
	return NewGrpcQuerier(k.cdc, k.storeKey, k)
}
