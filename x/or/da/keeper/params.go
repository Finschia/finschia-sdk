package keeper

import (
	"github.com/Finschia/finschia-sdk/x/or/da/types"
	"time"

	sdk "github.com/Finschia/finschia-sdk/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte{types.ParamsKey})
	var params types.Params
	if bz == nil {
		return params
	}
	k.cdc.MustUnmarshal(bz, &params)
	return params
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) error {
	store := ctx.KVStore(k.storeKey)
	bz, err := k.cdc.Marshal(&params)
	if err != nil {
		return err
	}

	store.Set([]byte{types.ParamsKey}, bz)
	return nil
}

func (k Keeper) CCBatchMaxBytes(ctx sdk.Context) uint64 {
	params := k.GetParams(ctx)

	return params.CCBatchMaxBytes
}

func (k Keeper) MaxQueueTxSize(ctx sdk.Context) uint64 {
	params := k.GetParams(ctx)
	return params.MaxQueueTxSize
}

func (k Keeper) MinQueueTxGas(ctx sdk.Context) uint64 {
	params := k.GetParams(ctx)
	return params.MinQueueTxGas
}

func (k Keeper) QueueTxExpirationWindow(ctx sdk.Context) uint64 {
	params := k.GetParams(ctx)
	return params.QueueTxExpirationWindow
}

func (k Keeper) SCCBatchMaxBytes(ctx sdk.Context) uint64 {
	params := k.GetParams(ctx)
	return params.SCCBatchMaxBytes
}

func (k Keeper) FraudProofWindow(ctx sdk.Context) time.Duration {
	params := k.GetParams(ctx)
	return time.Second * time.Duration(params.FraudProofWindow)
}

func (k Keeper) SequencerPublishWindow(ctx sdk.Context) time.Duration {
	params := k.GetParams(ctx)
	return time.Second * time.Duration(params.SequencerPublishWindow)
}
