package keeper

import (
	"errors"
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/Finschia/finschia-sdk/codec"
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/fbridge/types"
)

type Keeper struct {
	storeKey   sdk.StoreKey
	memKey     sdk.StoreKey
	cdc        codec.BinaryCodec
	authKeeper types.AccountKeeper
	bankKeeper types.BankKeeper

	// the target denom for the bridge
	targetDenom string
}

func NewKeeper(
	cdc codec.BinaryCodec,
	key, memKey sdk.StoreKey,
	authKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	targetDenom string,
) Keeper {
	if addr := authKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(errors.New("fbridge module account has not been set"))
	}

	return Keeper{
		storeKey:    key,
		memKey:      memKey,
		cdc:         cdc,
		authKeeper:  authKeeper,
		bankKeeper:  bankKeeper,
		targetDenom: targetDenom,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

func (k Keeper) InitMemStore(ctx sdk.Context) {
	memStore := ctx.KVStore(k.memKey)
	memStoreType := memStore.GetStoreType()
	if memStoreType != sdk.StoreTypeMemory {
		panic(fmt.Sprintf("invalid memory store type; got %s, expected: %s", memStoreType, sdk.StoreTypeMemory))
	}

	// create context with no block gas meter to ensure we do not consume gas during local initialization logic.
	noGasCtx := ctx.WithBlockGasMeter(sdk.NewInfiniteGasMeter())

	roleMetadata := types.RoleMetadata{}
	if !k.IsInitialized(noGasCtx) {
		memStore := noGasCtx.KVStore(k.memKey)

		for _, pair := range k.GetRolePairs(noGasCtx) {
			switch pair.Role {
			case types.RoleGuardian:
				roleMetadata.Guardian++
			case types.RoleOperator:
				roleMetadata.Operator++
			case types.RoleJudge:
				roleMetadata.Judge++
			}
		}
		memStore.Set(types.KeyMemRoleMetadata, k.cdc.MustMarshal(&roleMetadata))

		bsMeta := types.BridgeStatusMetadata{Inactive: 0, Active: 0}
		for _, bs := range k.GetBridgeSwitches(noGasCtx) {
			switch bs.Status {
			case types.StatusInactive:
				bsMeta.Inactive++
			case types.StatusActive:
				bsMeta.Active++
			default:
				panic("invalid bridge switch status")
			}
		}
		k.setBridgeStatusMetadata(noGasCtx, types.BridgeStatusMetadata{})

		memStore.Set(types.KeyMemInitialized, []byte{1})
	}
}

// IsInitialized returns true if the keeper is properly initialized, and false otherwise.
func (k Keeper) IsInitialized(ctx sdk.Context) bool {
	memStore := ctx.KVStore(k.memKey)
	return memStore.Get(types.KeyMemInitialized) != nil
}

func (k Keeper) IsBridgeHalted(ctx sdk.Context) bool {
	return k.GetBridgeStatus(ctx) == types.StatusInactive
}

func (k Keeper) setRoleMetadata(ctx sdk.Context, data types.RoleMetadata) {
	memStore := ctx.KVStore(k.memKey)
	bz := k.cdc.MustMarshal(&data)
	memStore.Set(types.KeyMemRoleMetadata, bz)
}

func (k Keeper) GetRoleMetadata(ctx sdk.Context) types.RoleMetadata {
	memStore := ctx.KVStore(k.memKey)
	data := types.RoleMetadata{}
	bz := memStore.Get(types.KeyMemRoleMetadata)
	if bz == nil {
		panic("role metadata must be set at genesis")
	}
	k.cdc.MustUnmarshal(bz, &data)
	return data
}

func (k Keeper) GetBridgeStatus(ctx sdk.Context) types.BridgeStatus {
	memStore := ctx.KVStore(k.memKey)
	bsMeta := types.BridgeStatusMetadata{}
	bz := memStore.Get(types.KeyMemBridgeStatus)
	k.cdc.MustUnmarshal(bz, &bsMeta)

	if types.CheckTrustLevelThreshold(bsMeta.Active+bsMeta.Inactive, bsMeta.Active, k.GetParams(ctx).GuardianTrustLevel) {
		return types.StatusActive
	}

	return types.StatusInactive
}

func (k Keeper) setBridgeStatusMetadata(ctx sdk.Context, status types.BridgeStatusMetadata) {
	memStore := ctx.KVStore(k.memKey)
	memStore.Set(types.KeyMemBridgeStatus, k.cdc.MustMarshal(&status))
}

func (k Keeper) GetBridgeStatusMetadata(ctx sdk.Context) types.BridgeStatusMetadata {
	memStore := ctx.KVStore(k.memKey)
	bsMeta := types.BridgeStatusMetadata{}
	bz := memStore.Get(types.KeyMemBridgeStatus)
	k.cdc.MustUnmarshal(bz, &bsMeta)
	return bsMeta
}
