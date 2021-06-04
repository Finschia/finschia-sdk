package keeper

import (
	"container/list"
	"fmt"

	"github.com/line/ostracon/libs/log"

	"github.com/line/lbm-sdk/v2/codec"
	sdk "github.com/line/lbm-sdk/v2/types"
	paramtypes "github.com/line/lbm-sdk/v2/x/params/types"
	"github.com/line/lbm-sdk/v2/x/staking/types"
)

// Implements ValidatorSet interface
var _ types.ValidatorSet = Keeper{}

// Implements DelegationSet interface
var _ types.DelegationSet = Keeper{}

// keeper of the staking store
type Keeper struct {
	storeKey           sdk.StoreKey
	cdc                codec.BinaryMarshaler
	authKeeper         types.AccountKeeper
	bankKeeper         types.BankKeeper
	hooks              types.StakingHooks
	paramstore         *paramtypes.Subspace
	validatorCacheList *list.List
}

// NewKeeper creates a new staking Keeper instance
func NewKeeper(
	cdc codec.BinaryMarshaler, key sdk.StoreKey, ak types.AccountKeeper, bk types.BankKeeper,
	ps *paramtypes.Subspace,
) Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	// ensure bonded and not bonded module accounts are set
	if addr := ak.GetModuleAddress(types.BondedPoolName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.BondedPoolName))
	}

	if addr := ak.GetModuleAddress(types.NotBondedPoolName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.NotBondedPoolName))
	}

	return Keeper{
		storeKey:           key,
		cdc:                cdc,
		authKeeper:         ak,
		bankKeeper:         bk,
		paramstore:         ps,
		hooks:              nil,
		validatorCacheList: list.New(),
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

// Set the validator hooks
func (k *Keeper) SetHooks(sh types.StakingHooks) *Keeper {
	if k.hooks != nil {
		panic("cannot set validator hooks twice")
	}

	k.hooks = sh

	return k
}

func GetIntProtoUnmarshalFunc(cdc codec.BinaryMarshaler) func (value []byte) interface{} {
	return func (value []byte) interface{} {
		val := sdk.IntProto{}
		cdc.MustUnmarshalBinaryBare(value, &val)
		return &val
	}
}

func GetIntProtoMarshalFunc(cdc codec.BinaryMarshaler) func (obj interface{}) []byte {
	return func (obj interface{}) []byte {
		return cdc.MustMarshalBinaryBare(obj.(*sdk.IntProto))
	}
}

// Load the last total validator power.
func (k Keeper) GetLastTotalPower(ctx sdk.Context) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	val := store.Get(types.LastTotalPowerKey, GetIntProtoUnmarshalFunc(k.cdc))

	if val == nil {
		return sdk.ZeroInt()
	}

	return (*val.(*sdk.IntProto)).Int
}

// Set the last total validator power.
func (k Keeper) SetLastTotalPower(ctx sdk.Context, power sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.LastTotalPowerKey, &sdk.IntProto{Int: power}, GetIntProtoMarshalFunc(k.cdc))}
