package keeper

import (
	"github.com/Finschia/finschia-sdk/store/prefix"
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/or/rollup/types"
)

func (k Keeper) SetSequencersByRollup(ctx sdk.Context, sequencersByRollup types.SequencersByRollup) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.SequencersByRollupKeyPrefix)
	b := k.cdc.MustMarshal(&sequencersByRollup)
	store.Set(types.SequencersByRollupKey(
		sequencersByRollup.RollupName,
	), b)
}

func (k Keeper) GetSequencersByRollupName(ctx sdk.Context, rollupName string) (val types.SequencersByRollup, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.SequencersByRollupKeyPrefix)

	b := store.Get(types.SequencersByRollupKey(
		rollupName,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) SetSequencer(ctx sdk.Context, sequencer types.Sequencer) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.SequencerKeyPrefix)
	b := k.cdc.MustMarshal(&sequencer)
	store.Set(types.SequencerKey(
		sequencer.SequencerAddress,
	), b)
}

func (k Keeper) GetSequencer(ctx sdk.Context, sequencerAddress string) (val types.Sequencer, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.SequencerKeyPrefix)

	b := store.Get(types.SequencerKey(
		sequencerAddress,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) SetDeposit(ctx sdk.Context, depoist types.Deposit) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DepositKeyPrefix)

	b := k.cdc.MustMarshal(&depoist)
	store.Set(types.DepositKey(
		depoist.RollupName,
		depoist.SequencerAddress,
	), b)
}

func (k Keeper) GetDeposit(ctx sdk.Context, rollupName, sequencerAddress string) (val types.Deposit, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DepositKeyPrefix)

	b := store.Get(types.DepositKey(
		rollupName,
		sequencerAddress,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}
