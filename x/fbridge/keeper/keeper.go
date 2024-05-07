package keeper

import (
	"encoding/binary"
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

	// authority can give a role to a specific address like guardian
	authority string
}

func NewKeeper(
	cdc codec.BinaryCodec,
	key, memKey sdk.StoreKey,
	authKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	targetDenom, authority string,
) Keeper {
	if addr := authKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(errors.New("fbridge module account has not been set"))
	}

	for _, addr := range types.AuthorityCandiates() {
		if authority == addr.String() {
			break
		}
		panic("x/bridge authority must be the gov or foundation module account")
	}

	return Keeper{
		storeKey:    key,
		memKey:      memKey,
		cdc:         cdc,
		authKeeper:  authKeeper,
		bankKeeper:  bankKeeper,
		targetDenom: targetDenom,
		authority:   authority,
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
		k.setRoleMetadata(noGasCtx, roleMetadata)

		nInactive := k.GetBridgeInactiveCounter(noGasCtx)
		k.setBridgeInactiveCounter(noGasCtx, nInactive)

		memStore := noGasCtx.KVStore(k.memKey)
		memStore.Set(types.KeyMemInitialized, []byte{1})
	}
}

// IsInitialized returns true if the keeper is properly initialized, and false otherwise.
func (k Keeper) IsInitialized(ctx sdk.Context) bool {
	memStore := ctx.KVStore(k.memKey)
	return memStore.Get(types.KeyMemInitialized) != nil
}

func (k Keeper) GetAuthority() string {
	return k.authority
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
		return types.RoleMetadata{}
	}
	k.cdc.MustUnmarshal(bz, &data)
	return data
}

func (k Keeper) GetBridgeStatus(ctx sdk.Context) types.BridgeStatus {
	roleMeta := k.GetRoleMetadata(ctx)
	bsMeta := k.GetBridgeStatusMetadata(ctx)
	if types.CheckTrustLevelThreshold(roleMeta.Guardian, bsMeta.Inactive, k.GetParams(ctx).GuardianTrustLevel) {
		return types.StatusInactive
	}

	return types.StatusActive
}

func (k Keeper) setBridgeInactiveCounter(ctx sdk.Context, nInactive uint64) {
	memStore := ctx.KVStore(k.memKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, nInactive)
	memStore.Set(types.KeyMemBridgeInactiveCounter, bz)
}

func (k Keeper) GetBridgeInactiveCounter(ctx sdk.Context) uint64 {
	memStore := ctx.KVStore(k.memKey)
	bz := memStore.Get(types.KeyMemBridgeInactiveCounter)
	if bz == nil {
		n := uint64(0)
		for _, bs := range k.GetBridgeSwitches(ctx) {
			if bs.Status == types.StatusInactive {
				n++
			}
		}
		return n
	}

	return binary.BigEndian.Uint64(bz)
}

func (k Keeper) GetBridgeStatusMetadata(ctx sdk.Context) types.BridgeStatusMetadata {
	bsMeta := types.BridgeStatusMetadata{}
	bsMeta.Inactive = k.GetBridgeInactiveCounter(ctx)
	bsMeta.Active = k.GetRoleMetadata(ctx).Guardian - bsMeta.Inactive
	return bsMeta
}
