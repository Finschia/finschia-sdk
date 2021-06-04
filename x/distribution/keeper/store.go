package keeper

import (
	gogotypes "github.com/gogo/protobuf/types"
	"github.com/line/lbm-sdk/v2/codec"
	types2 "github.com/line/lbm-sdk/v2/store/types"

	sdk "github.com/line/lbm-sdk/v2/types"
	"github.com/line/lbm-sdk/v2/x/distribution/types"
)

// get the delegator withdraw address, defaulting to the delegator address
func (k Keeper) GetDelegatorWithdrawAddr(ctx sdk.Context, delAddr sdk.AccAddress) sdk.AccAddress {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetDelegatorWithdrawAddrKey(delAddr), types2.GetBytesUnmarshalFunc())
	if b == nil {
		return delAddr
	}
	return b.([]byte)
}

// set the delegator withdraw address
func (k Keeper) SetDelegatorWithdrawAddr(ctx sdk.Context, delAddr, withdrawAddr sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetDelegatorWithdrawAddrKey(delAddr), withdrawAddr.Bytes(), types2.GetBytesMarshalFunc())
}

// delete a delegator withdraw addr
func (k Keeper) DeleteDelegatorWithdrawAddr(ctx sdk.Context, delAddr, withdrawAddr sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetDelegatorWithdrawAddrKey(delAddr))
}

// iterate over delegator withdraw addrs
func (k Keeper) IterateDelegatorWithdrawAddrs(ctx sdk.Context, handler func(del sdk.AccAddress, addr sdk.AccAddress) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.DelegatorWithdrawAddrPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		addr := sdk.AccAddress(iter.Value())
		del := types.GetDelegatorWithdrawInfoAddress(iter.Key())
		if handler(del, addr) {
			break
		}
	}
}

func GetFeePoolUnmarshalFunc(cdc codec.BinaryMarshaler) func (value []byte) interface{} {
	return func (value []byte) interface{} {
		val := types.FeePool{}
		cdc.MustUnmarshalBinaryBare(value, &val)
		return &val
	}
}

func GetFeePoolMarshalFunc(cdc codec.BinaryMarshaler) func (obj interface{}) []byte {
	return func (obj interface{}) []byte {
		return cdc.MustMarshalBinaryBare(obj.(*types.FeePool))
	}
}

// get the global fee pool distribution info
func (k Keeper) GetFeePool(ctx sdk.Context) (feePool types.FeePool) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.FeePoolKey, GetFeePoolUnmarshalFunc(k.cdc))
	if b == nil {
		panic("Stored fee pool should not have been nil")
	}
	return *b.(*types.FeePool)
}

// set the global fee pool distribution info
func (k Keeper) SetFeePool(ctx sdk.Context, feePool types.FeePool) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.FeePoolKey, &feePool, GetFeePoolMarshalFunc(k.cdc))
}

func GetByteValueUnmarshalFunc(cdc codec.BinaryMarshaler) func (value []byte) interface{} {
	return func (value []byte) interface{} {
		val := gogotypes.BytesValue{}
		cdc.MustUnmarshalBinaryBare(value, &val)
		return &val
	}
}

func GetByteValueMarshalFunc(cdc codec.BinaryMarshaler) func (obj interface{}) []byte {
	return func (obj interface{}) []byte {
		return cdc.MustMarshalBinaryBare(obj.(*gogotypes.BytesValue))
	}
}

// GetPreviousProposerConsAddr returns the proposer consensus address for the
// current block.
func (k Keeper) GetPreviousProposerConsAddr(ctx sdk.Context) sdk.ConsAddress {
	store := ctx.KVStore(k.storeKey)
	val := store.Get(types.ProposerKey, GetByteValueUnmarshalFunc(k.cdc))
	if val == nil {
		panic("previous proposer not set")
	}
	return (*val.(*gogotypes.BytesValue)).GetValue()
}

// set the proposer public key for this block
func (k Keeper) SetPreviousProposerConsAddr(ctx sdk.Context, consAddr sdk.ConsAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.ProposerKey, &gogotypes.BytesValue{Value: consAddr}, GetByteValueMarshalFunc(k.cdc))
}

func GetDelegatorStartingInfoUnmarshalFunc(cdc codec.BinaryMarshaler) func (value []byte) interface{} {
	return func (value []byte) interface{} {
		val := types.DelegatorStartingInfo{}
		cdc.MustUnmarshalBinaryBare(value, &val)
		return &val
	}
}

func GetDelegatorStartingInfoMarshalFunc(cdc codec.BinaryMarshaler) func (obj interface{}) []byte {
	return func (obj interface{}) []byte {
		return cdc.MustMarshalBinaryBare(obj.(*types.DelegatorStartingInfo))
	}
}

// get the starting info associated with a delegator
func (k Keeper) GetDelegatorStartingInfo(ctx sdk.Context, val sdk.ValAddress, del sdk.AccAddress) (period types.DelegatorStartingInfo) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetDelegatorStartingInfoKey(val, del), GetDelegatorStartingInfoUnmarshalFunc(k.cdc))
	return *b.(*types.DelegatorStartingInfo)
}

// set the starting info associated with a delegator
func (k Keeper) SetDelegatorStartingInfo(ctx sdk.Context, val sdk.ValAddress, del sdk.AccAddress, period types.DelegatorStartingInfo) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetDelegatorStartingInfoKey(val, del), &period, GetDelegatorStartingInfoMarshalFunc(k.cdc))
}

// check existence of the starting info associated with a delegator
func (k Keeper) HasDelegatorStartingInfo(ctx sdk.Context, val sdk.ValAddress, del sdk.AccAddress) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetDelegatorStartingInfoKey(val, del))
}

// delete the starting info associated with a delegator
func (k Keeper) DeleteDelegatorStartingInfo(ctx sdk.Context, val sdk.ValAddress, del sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetDelegatorStartingInfoKey(val, del))
}

// iterate over delegator starting infos
func (k Keeper) IterateDelegatorStartingInfos(ctx sdk.Context, handler func(val sdk.ValAddress, del sdk.AccAddress, info types.DelegatorStartingInfo) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.DelegatorStartingInfoPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		info := iter.ValueObject(GetDelegatorStartingInfoUnmarshalFunc(k.cdc))
		val, del := types.GetDelegatorStartingInfoAddresses(iter.Key())
		if handler(val, del, *info.(*types.DelegatorStartingInfo)) {
			break
		}
	}
}

func GetValidatorHistoricalRewardsUnmarshalFunc(cdc codec.BinaryMarshaler) func (value []byte) interface{} {
	return func (value []byte) interface{} {
		val := types.ValidatorHistoricalRewards{}
		cdc.MustUnmarshalBinaryBare(value, &val)
		return &val
	}
}

func GetValidatorHistoricalRewardsMarshalFunc(cdc codec.BinaryMarshaler) func (obj interface{}) []byte {
	return func (obj interface{}) []byte {
		return cdc.MustMarshalBinaryBare(obj.(*types.ValidatorHistoricalRewards))
	}
}

// get historical rewards for a particular period
func (k Keeper) GetValidatorHistoricalRewards(ctx sdk.Context, val sdk.ValAddress, period uint64) (rewards types.ValidatorHistoricalRewards) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetValidatorHistoricalRewardsKey(val, period), GetValidatorHistoricalRewardsUnmarshalFunc(k.cdc))
	return *b.(*types.ValidatorHistoricalRewards)
}

// set historical rewards for a particular period
func (k Keeper) SetValidatorHistoricalRewards(ctx sdk.Context, val sdk.ValAddress, period uint64, rewards types.ValidatorHistoricalRewards) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetValidatorHistoricalRewardsKey(val, period), &rewards, GetValidatorHistoricalRewardsMarshalFunc(k.cdc))
}

// iterate over historical rewards
func (k Keeper) IterateValidatorHistoricalRewards(ctx sdk.Context, handler func(val sdk.ValAddress, period uint64, rewards types.ValidatorHistoricalRewards) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.ValidatorHistoricalRewardsPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		rewards := iter.ValueObject(GetValidatorHistoricalRewardsUnmarshalFunc(k.cdc))
		addr, period := types.GetValidatorHistoricalRewardsAddressPeriod(iter.Key())
		if handler(addr, period, *rewards.(*types.ValidatorHistoricalRewards)) {
			break
		}
	}
}

// delete a historical reward
func (k Keeper) DeleteValidatorHistoricalReward(ctx sdk.Context, val sdk.ValAddress, period uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetValidatorHistoricalRewardsKey(val, period))
}

// delete historical rewards for a validator
func (k Keeper) DeleteValidatorHistoricalRewards(ctx sdk.Context, val sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetValidatorHistoricalRewardsPrefix(val))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		store.Delete(iter.Key())
	}
}

// delete all historical rewards
func (k Keeper) DeleteAllValidatorHistoricalRewards(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.ValidatorHistoricalRewardsPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		store.Delete(iter.Key())
	}
}

// historical reference count (used for testcases)
func (k Keeper) GetValidatorHistoricalReferenceCount(ctx sdk.Context) (count uint64) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.ValidatorHistoricalRewardsPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		rewards := iter.ValueObject(GetValidatorHistoricalRewardsUnmarshalFunc(k.cdc))
		count += uint64((*rewards.(*types.ValidatorHistoricalRewards)).ReferenceCount)
	}
	return
}

func GetValidatorCurrentRewardsUnmarshalFunc(cdc codec.BinaryMarshaler) func (value []byte) interface{} {
	return func (value []byte) interface{} {
		val := types.ValidatorCurrentRewards{}
		cdc.MustUnmarshalBinaryBare(value, &val)
		return &val
	}
}

func GetValidatorCurrentRewardsMarshalFunc(cdc codec.BinaryMarshaler) func (obj interface{}) []byte {
	return func (obj interface{}) []byte {
		return cdc.MustMarshalBinaryBare(obj.(*types.ValidatorCurrentRewards))
	}
}

// get current rewards for a validator
func (k Keeper) GetValidatorCurrentRewards(ctx sdk.Context, val sdk.ValAddress) (rewards types.ValidatorCurrentRewards) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetValidatorCurrentRewardsKey(val), GetValidatorCurrentRewardsUnmarshalFunc(k.cdc))
	return *b.(*types.ValidatorCurrentRewards)
}

// set current rewards for a validator
func (k Keeper) SetValidatorCurrentRewards(ctx sdk.Context, val sdk.ValAddress, rewards types.ValidatorCurrentRewards) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetValidatorCurrentRewardsKey(val), &rewards, GetValidatorCurrentRewardsMarshalFunc(k.cdc))
}

// delete current rewards for a validator
func (k Keeper) DeleteValidatorCurrentRewards(ctx sdk.Context, val sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetValidatorCurrentRewardsKey(val))
}

// iterate over current rewards
func (k Keeper) IterateValidatorCurrentRewards(ctx sdk.Context, handler func(val sdk.ValAddress, rewards types.ValidatorCurrentRewards) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.ValidatorCurrentRewardsPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		val := iter.ValueObject(GetValidatorCurrentRewardsUnmarshalFunc(k.cdc))
		addr := types.GetValidatorCurrentRewardsAddress(iter.Key())
		if handler(addr,  *val.(*types.ValidatorCurrentRewards)) {
			break
		}
	}
}

func GetValidatorAccumulatedCommissionUnmarshalFunc(cdc codec.BinaryMarshaler) func (value []byte) interface{} {
	return func (value []byte) interface{} {
		val := types.ValidatorAccumulatedCommission{}
		cdc.MustUnmarshalBinaryBare(value, &val)
		return &val
	}
}

func GetValidatorAccumulatedCommissionMarshalFunc(cdc codec.BinaryMarshaler) func (obj interface{}) []byte {
	return func (obj interface{}) []byte {
		return cdc.MustMarshalBinaryBare(obj.(*types.ValidatorAccumulatedCommission))
	}
}

// get accumulated commission for a validator
func (k Keeper) GetValidatorAccumulatedCommission(ctx sdk.Context, val sdk.ValAddress) (commission types.ValidatorAccumulatedCommission) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetValidatorAccumulatedCommissionKey(val), GetValidatorAccumulatedCommissionUnmarshalFunc(k.cdc))
	if b == nil {
		return types.ValidatorAccumulatedCommission{}
	}
	return *b.(*types.ValidatorAccumulatedCommission)
}

// set accumulated commission for a validator
func (k Keeper) SetValidatorAccumulatedCommission(ctx sdk.Context, val sdk.ValAddress, commission types.ValidatorAccumulatedCommission) {
	store := ctx.KVStore(k.storeKey)
	if commission.Commission.IsZero() {
		store.Set(types.GetValidatorAccumulatedCommissionKey(val), &types.ValidatorAccumulatedCommission{},
			GetValidatorAccumulatedCommissionMarshalFunc(k.cdc))
	} else {
		store.Set(types.GetValidatorAccumulatedCommissionKey(val), &commission,
			GetValidatorAccumulatedCommissionMarshalFunc(k.cdc))
	}
}

// delete accumulated commission for a validator
func (k Keeper) DeleteValidatorAccumulatedCommission(ctx sdk.Context, val sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetValidatorAccumulatedCommissionKey(val))
}

// iterate over accumulated commissions
func (k Keeper) IterateValidatorAccumulatedCommissions(ctx sdk.Context, handler func(val sdk.ValAddress, commission types.ValidatorAccumulatedCommission) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.ValidatorAccumulatedCommissionPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		val := iter.ValueObject(GetValidatorAccumulatedCommissionUnmarshalFunc(k.cdc))
		addr := types.GetValidatorAccumulatedCommissionAddress(iter.Key())
		if handler(addr, *val.(*types.ValidatorAccumulatedCommission)) {
			break
		}
	}
}

func GetValidatorOutstandingRewardsUnmarshalFunc(cdc codec.BinaryMarshaler) func (value []byte) interface{} {
	return func (value []byte) interface{} {
		val := types.ValidatorOutstandingRewards{}
		cdc.MustUnmarshalBinaryBare(value, &val)
		return &val
	}
}

func GetValidatorOutstandingRewardsMarshalFunc(cdc codec.BinaryMarshaler) func (obj interface{}) []byte {
	return func (obj interface{}) []byte {
		return cdc.MustMarshalBinaryBare(obj.(*types.ValidatorOutstandingRewards))
	}
}

// get validator outstanding rewards
func (k Keeper) GetValidatorOutstandingRewards(ctx sdk.Context, val sdk.ValAddress) types.ValidatorOutstandingRewards {
	store := ctx.KVStore(k.storeKey)
	v := store.Get(types.GetValidatorOutstandingRewardsKey(val), GetValidatorOutstandingRewardsUnmarshalFunc(k.cdc))
	return *v.(*types.ValidatorOutstandingRewards)
}

// set validator outstanding rewards
func (k Keeper) SetValidatorOutstandingRewards(ctx sdk.Context, val sdk.ValAddress, rewards types.ValidatorOutstandingRewards) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetValidatorOutstandingRewardsKey(val), &rewards, GetValidatorOutstandingRewardsMarshalFunc(k.cdc))
}

// delete validator outstanding rewards
func (k Keeper) DeleteValidatorOutstandingRewards(ctx sdk.Context, val sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetValidatorOutstandingRewardsKey(val))
}

// iterate validator outstanding rewards
func (k Keeper) IterateValidatorOutstandingRewards(ctx sdk.Context, handler func(val sdk.ValAddress, rewards types.ValidatorOutstandingRewards) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.ValidatorOutstandingRewardsPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		rewards := iter.ValueObject(GetValidatorOutstandingRewardsUnmarshalFunc(k.cdc))
		addr := types.GetValidatorOutstandingRewardsAddress(iter.Key())
		if handler(addr, *rewards.(*types.ValidatorOutstandingRewards)) {
			break
		}
	}
}

func GetValidatorSlashEventUnmarshalFunc(cdc codec.BinaryMarshaler) func (value []byte) interface{} {
	return func (value []byte) interface{} {
		val := types.ValidatorSlashEvent{}
		cdc.MustUnmarshalBinaryBare(value, &val)
		return &val
	}
}

func GetValidatorSlashEventMarshalFunc(cdc codec.BinaryMarshaler) func (obj interface{}) []byte {
	return func (obj interface{}) []byte {
		return cdc.MustMarshalBinaryBare(obj.(*types.ValidatorSlashEvent))
	}
}

// get slash event for height
func (k Keeper) GetValidatorSlashEvent(ctx sdk.Context, val sdk.ValAddress, height, period uint64) (event types.ValidatorSlashEvent, found bool) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetValidatorSlashEventKey(val, height, period), GetValidatorSlashEventUnmarshalFunc(k.cdc))
	if b == nil {
		return types.ValidatorSlashEvent{}, false
	}
	return *b.(*types.ValidatorSlashEvent), true
}

// set slash event for height
func (k Keeper) SetValidatorSlashEvent(ctx sdk.Context, val sdk.ValAddress, height, period uint64, event types.ValidatorSlashEvent) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetValidatorSlashEventKey(val, height, period), &event, GetValidatorSlashEventMarshalFunc(k.cdc))
}

// iterate over slash events between heights, inclusive
func (k Keeper) IterateValidatorSlashEventsBetween(ctx sdk.Context, val sdk.ValAddress, startingHeight uint64, endingHeight uint64,
	handler func(height uint64, event types.ValidatorSlashEvent) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := store.Iterator(
		types.GetValidatorSlashEventKeyPrefix(val, startingHeight),
		types.GetValidatorSlashEventKeyPrefix(val, endingHeight+1),
	)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		event := iter.ValueObject(GetValidatorSlashEventUnmarshalFunc(k.cdc))
		_, height := types.GetValidatorSlashEventAddressHeight(iter.Key())
		if handler(height, *event.(*types.ValidatorSlashEvent)) {
			break
		}
	}
}

// iterate over all slash events
func (k Keeper) IterateValidatorSlashEvents(ctx sdk.Context, handler func(val sdk.ValAddress, height uint64, event types.ValidatorSlashEvent) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.ValidatorSlashEventPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		event := iter.ValueObject(GetValidatorSlashEventUnmarshalFunc(k.cdc))
		val, height := types.GetValidatorSlashEventAddressHeight(iter.Key())
		if handler(val, height, *event.(*types.ValidatorSlashEvent)) {
			break
		}
	}
}

// delete slash events for a particular validator
func (k Keeper) DeleteValidatorSlashEvents(ctx sdk.Context, val sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetValidatorSlashEventPrefix(val))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		store.Delete(iter.Key())
	}
}

// delete all slash events
func (k Keeper) DeleteAllValidatorSlashEvents(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.ValidatorSlashEventPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		store.Delete(iter.Key())
	}
}
