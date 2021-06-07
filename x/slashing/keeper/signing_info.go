package keeper

import (
	"time"

	gogotypes "github.com/gogo/protobuf/types"
	"github.com/line/lbm-sdk/v2/codec"

	sdk "github.com/line/lbm-sdk/v2/types"
	"github.com/line/lbm-sdk/v2/x/slashing/types"
)

func GetValidatorSigningInfoUnmarshalFunc(cdc codec.BinaryMarshaler) func (value []byte) interface{} {
	return func (value []byte) interface{} {
		val := types.ValidatorSigningInfo{}
		cdc.MustUnmarshalBinaryBare(value, &val)
		return &val
	}
}

func GetValidatorSigningInfoMarshalFunc(cdc codec.BinaryMarshaler) func (obj interface{}) []byte {
	return func (obj interface{}) []byte {
		return cdc.MustMarshalBinaryBare(obj.(*types.ValidatorSigningInfo))
	}
}

// GetValidatorSigningInfo retruns the ValidatorSigningInfo for a specific validator
// ConsAddress
func (k Keeper) GetValidatorSigningInfo(ctx sdk.Context, address sdk.ConsAddress) (info types.ValidatorSigningInfo, found bool) {
	store := ctx.KVStore(k.storeKey)
	val := store.Get(types.ValidatorSigningInfoKey(address), GetValidatorSigningInfoUnmarshalFunc(k.cdc))
	if val == nil {
		found = false
		return
	}
	info = *val.(*types.ValidatorSigningInfo)
	found = true
	return
}

// HasValidatorSigningInfo returns if a given validator has signing information
// persited.
func (k Keeper) HasValidatorSigningInfo(ctx sdk.Context, consAddr sdk.ConsAddress) bool {
	_, ok := k.GetValidatorSigningInfo(ctx, consAddr)
	return ok
}

// SetValidatorSigningInfo sets the validator signing info to a consensus address key
func (k Keeper) SetValidatorSigningInfo(ctx sdk.Context, address sdk.ConsAddress, info types.ValidatorSigningInfo) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.ValidatorSigningInfoKey(address), &info, GetValidatorSigningInfoMarshalFunc(k.cdc))
}

// IterateValidatorSigningInfos iterates over the stored ValidatorSigningInfo
func (k Keeper) IterateValidatorSigningInfos(ctx sdk.Context,
	handler func(address sdk.ConsAddress, info types.ValidatorSigningInfo) (stop bool)) {

	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.ValidatorSigningInfoKeyPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		address := types.ValidatorSigningInfoAddress(iter.Key())
		info := iter.ValueObject(GetValidatorSigningInfoUnmarshalFunc(k.cdc))
		if handler(address, *info.(*types.ValidatorSigningInfo)) {
			break
		}
	}
}

func GetBoolValueUnmarshalFunc(cdc codec.BinaryMarshaler) func (value []byte) interface{} {
	return func (value []byte) interface{} {
		var val gogotypes.BoolValue
		cdc.MustUnmarshalBinaryBare(value, &val)
		return &val
	}
}

func GetBoolValueMarshalFunc(cdc codec.BinaryMarshaler) func (obj interface{}) []byte {
	return func (obj interface{}) []byte {
		return cdc.MustMarshalBinaryBare(obj.(*gogotypes.BoolValue))
	}
}

// GetValidatorMissedBlockBitArray gets the bit for the missed blocks array
func (k Keeper) GetValidatorMissedBlockBitArray(ctx sdk.Context, address sdk.ConsAddress, index int64) bool {
	store := ctx.KVStore(k.storeKey)
	val := store.Get(types.ValidatorMissedBlockBitArrayKey(address, index), GetBoolValueUnmarshalFunc(k.cdc))
	if val == nil {
		// lazy: treat empty key as not missed
		return false
	}
	return (*val.(*gogotypes.BoolValue)).Value
}

// IterateValidatorMissedBlockBitArray iterates over the signed blocks window
// and performs a callback function
func (k Keeper) IterateValidatorMissedBlockBitArray(ctx sdk.Context,
	address sdk.ConsAddress, handler func(index int64, missed bool) (stop bool)) {

	store := ctx.KVStore(k.storeKey)
	index := int64(0)
	// Array may be sparse
	for ; index < k.SignedBlocksWindow(ctx); index++ {
		val := store.Get(types.ValidatorMissedBlockBitArrayKey(address, index), GetBoolValueUnmarshalFunc(k.cdc))
		if val == nil {
			continue
		}

		if handler(index, (*val.(*gogotypes.BoolValue)).Value) {
			break
		}
	}
}

// GetValidatorMissedBlocks returns array of missed blocks for given validator Cons address
func (k Keeper) GetValidatorMissedBlocks(ctx sdk.Context, address sdk.ConsAddress) []types.MissedBlock {
	missedBlocks := []types.MissedBlock{}
	k.IterateValidatorMissedBlockBitArray(ctx, address, func(index int64, missed bool) (stop bool) {
		missedBlocks = append(missedBlocks, types.NewMissedBlock(index, missed))
		return false
	})

	return missedBlocks
}

// JailUntil attempts to set a validator's JailedUntil attribute in its signing
// info. It will panic if the signing info does not exist for the validator.
func (k Keeper) JailUntil(ctx sdk.Context, consAddr sdk.ConsAddress, jailTime time.Time) {
	signInfo, ok := k.GetValidatorSigningInfo(ctx, consAddr)
	if !ok {
		panic("cannot jail validator that does not have any signing information")
	}

	signInfo.JailedUntil = jailTime
	k.SetValidatorSigningInfo(ctx, consAddr, signInfo)
}

// Tombstone attempts to tombstone a validator. It will panic if signing info for
// the given validator does not exist.
func (k Keeper) Tombstone(ctx sdk.Context, consAddr sdk.ConsAddress) {
	signInfo, ok := k.GetValidatorSigningInfo(ctx, consAddr)
	if !ok {
		panic("cannot tombstone validator that does not have any signing information")
	}

	if signInfo.Tombstoned {
		panic("cannot tombstone validator that is already tombstoned")
	}

	signInfo.Tombstoned = true
	k.SetValidatorSigningInfo(ctx, consAddr, signInfo)
}

// IsTombstoned returns if a given validator by consensus address is tombstoned.
func (k Keeper) IsTombstoned(ctx sdk.Context, consAddr sdk.ConsAddress) bool {
	signInfo, ok := k.GetValidatorSigningInfo(ctx, consAddr)
	if !ok {
		return false
	}

	return signInfo.Tombstoned
}

// SetValidatorMissedBlockBitArray sets the bit that checks if the validator has
// missed a block in the current window
func (k Keeper) SetValidatorMissedBlockBitArray(ctx sdk.Context, address sdk.ConsAddress, index int64, missed bool) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.ValidatorMissedBlockBitArrayKey(address, index), &gogotypes.BoolValue{Value: missed},
		GetBoolValueMarshalFunc(k.cdc))
}

// clearValidatorMissedBlockBitArray deletes every instance of ValidatorMissedBlockBitArray in the store
func (k Keeper) clearValidatorMissedBlockBitArray(ctx sdk.Context, address sdk.ConsAddress) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.ValidatorMissedBlockBitArrayPrefixKey(address))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		store.Delete(iter.Key())
	}
}
