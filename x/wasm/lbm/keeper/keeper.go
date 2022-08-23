package keeper

import (
	"github.com/line/lbm-sdk/codec"
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	authkeeper "github.com/line/lbm-sdk/x/auth/keeper"
	paramtypes "github.com/line/lbm-sdk/x/params/types"
	wasmkeeper "github.com/line/lbm-sdk/x/wasm/keeper"
	lbmwasmtypes "github.com/line/lbm-sdk/x/wasm/lbm/types"
	wasmtypes "github.com/line/lbm-sdk/x/wasm/types"
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
	bankKeeper wasmtypes.BankKeeper,
	stakingKeeper wasmtypes.StakingKeeper,
	distKeeper wasmtypes.DistributionKeeper,
	channelKeeper wasmtypes.ChannelKeeper,
	portKeeper wasmtypes.PortKeeper,
	capabilityKeeper wasmtypes.CapabilityKeeper,
	portSource wasmtypes.ICS20TransferPortSource,
	router wasmkeeper.MessageRouter,
	queryRouter wasmkeeper.GRPCQueryRouter,
	homeDir string,
	wasmConfig wasmtypes.WasmConfig,
	supportedFeatures string,
	customEncoders *wasmkeeper.MessageEncoders,
	customPlugins *wasmkeeper.QueryPlugins,
	opts ...wasmkeeper.Option,
) Keeper {
	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(lbmwasmtypes.ParamKeyTable())
	}

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

// activateContract delete the contract address from inactivateContract list if the contract is deactivated.
func (k Keeper) activateContract(ctx sdk.Context, contractAddress sdk.AccAddress) error {
	if !k.IsInactiveContract(ctx, contractAddress) {
		return sdkerrors.Wrapf(wasmtypes.ErrNotFound, "no inactivate contract %s", contractAddress.String())
	}
	if !k.HasContractInfo(ctx, contractAddress) {
		return sdkerrors.Wrapf(wasmtypes.ErrInvalid, "no contract %s", contractAddress)
	}

	k.deleteInactiveContract(ctx, contractAddress)

	return nil
}

// deactivateContract add the contract address to inactivateContract list.
func (k Keeper) deactivateContract(ctx sdk.Context, contractAddress sdk.AccAddress) error {
	if k.IsInactiveContract(ctx, contractAddress) {
		return sdkerrors.Wrapf(wasmtypes.ErrAccountExists, "already inactivate contract %s", contractAddress.String())
	}
	if !k.HasContractInfo(ctx, contractAddress) {
		return sdkerrors.Wrapf(wasmtypes.ErrInvalid, "no contract %s", contractAddress)
	}

	k.addInactiveContract(ctx, contractAddress)

	return nil
}

func Querier(k *Keeper) *GrpcQuerier {
	return NewGrpcQuerier(k.cdc, k.storeKey, k)
}
