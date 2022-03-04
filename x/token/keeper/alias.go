package keeper

import (
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/token"
)

// iterate through the balances and perform the provided function
func (k Keeper) iterateBalances(ctx sdk.Context, fn func(addr sdk.AccAddress, token token.FT) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, balanceKeyPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var amount sdk.Int
		if err := amount.Unmarshal(iterator.Value()); err != nil {
			panic(err)
		}

		addr, classID := splitBalanceKey(iterator.Key())
		token := token.FT{
			ClassId: classID,
			Amount:  amount,
		}

		stop := fn(addr, token)
		if stop {
			break
		}
	}
}

// iterate through the classes and perform the provided function
func (k Keeper) iterateClasses(ctx sdk.Context, fn func(class token.Token) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, classKeyPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var class token.Token
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &class)

		stop := fn(class)
		if stop {
			break
		}
	}
}

func (k Keeper) iterateGrants(ctx sdk.Context, fn func(grant token.Grant) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, grantKeyPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		grantee, classID, action := splitGrantKey(iterator.Key())
		grant := token.Grant{
			Grantee: grantee.String(),
			ClassId: classID,
			Action:  action,
		}

		stop := fn(grant)
		if stop {
			break
		}
	}
}

func (k Keeper) iterateApproves(ctx sdk.Context, fn func(approve token.Approve) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, approveKeyPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		classID, proxy, approver := splitApproveKey(iterator.Key())
		approve := token.Approve{
			Approver: approver.String(),
			Proxy:    proxy.String(),
			ClassId:  classID,
		}

		stop := fn(approve)
		if stop {
			break
		}
	}
}

func (k Keeper) iterateStatistics(ctx sdk.Context, fn func(amount token.FT) (stop bool), keyPrefix []byte) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, keyPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var amt sdk.Int
		if err := amt.Unmarshal(iterator.Value()); err != nil {
			panic(err)
		}

		classID := splitStatisticsKey(iterator.Key(), keyPrefix)
		amount := token.FT{
			ClassId: classID,
			Amount:  amt,
		}

		stop := fn(amount)
		if stop {
			break
		}
	}
}

func (k Keeper) iterateSupplies(ctx sdk.Context, fn func(amount token.FT) (stop bool)) {
	k.iterateStatistics(ctx, fn, supplyKeyPrefix)
}

func (k Keeper) iterateMints(ctx sdk.Context, fn func(amount token.FT) (stop bool)) {
	k.iterateStatistics(ctx, fn, mintKeyPrefix)
}

func (k Keeper) iterateBurns(ctx sdk.Context, fn func(amount token.FT) (stop bool)) {
	k.iterateStatistics(ctx, fn, burnKeyPrefix)
}
