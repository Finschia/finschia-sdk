package keeper

import (
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/token"
)

// iterate through the balances of a contract and perform the provided function
func (k Keeper) iterateContractBalances(ctx sdk.Context, classID string, fn func(balance token.Balance) (stop bool)) {
	k.iterateBalancesImpl(ctx, balanceKeyPrefixByContractID(classID), func(_ string, balance token.Balance) (stop bool) {
		return fn(balance)
	})
}

func (k Keeper) iterateBalancesImpl(ctx sdk.Context, prefix []byte, fn func(classID string, balance token.Balance) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		classID, addr := splitBalanceKey(iterator.Key())

		var amount sdk.Int
		if err := amount.Unmarshal(iterator.Value()); err != nil {
			panic(err)
		}
		balance := token.Balance{
			Address: addr.String(),
			Amount:  amount,
		}

		stop := fn(classID, balance)
		if stop {
			break
		}
	}
}

// iterate through the classes and perform the provided function
func (k Keeper) iterateClasses(ctx sdk.Context, fn func(class token.TokenClass) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, classKeyPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var class token.TokenClass
		k.cdc.MustUnmarshal(iterator.Value(), &class)

		stop := fn(class)
		if stop {
			break
		}
	}
}

func (k Keeper) iterateGrants(ctx sdk.Context, fn func(classID string, grant token.Grant) (stop bool)) {
	k.iterateGrantsImpl(ctx, grantKeyPrefix, fn)
}

func (k Keeper) iterateContractGrants(ctx sdk.Context, classID string, fn func(grant token.Grant) (stop bool)) {
	k.iterateGrantsImpl(ctx, grantKeyPrefixByContractID(classID), func(_ string, grant token.Grant) (stop bool) {
		return fn(grant)
	})
}

func (k Keeper) iterateGrantsImpl(ctx sdk.Context, prefix []byte, fn func(classID string, grant token.Grant) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		classID, grantee, permission := splitGrantKey(iterator.Key())
		grant := token.Grant{
			Grantee:    grantee.String(),
			Permission: token.Permission_name[int32(permission)],
		}

		stop := fn(classID, grant)
		if stop {
			break
		}
	}
}

func (k Keeper) iterateAuthorizations(ctx sdk.Context, fn func(classID string, authorization token.Authorization) (stop bool)) {
	k.iterateAuthorizationsImpl(ctx, authorizationKeyPrefix, fn)
}

func (k Keeper) iterateContractAuthorizations(ctx sdk.Context, classID string, fn func(authorization token.Authorization) (stop bool)) {
	k.iterateAuthorizationsImpl(ctx, authorizationKeyPrefixByContractID(classID), func(_ string, authorization token.Authorization) (stop bool) {
		return fn(authorization)
	})
}

func (k Keeper) iterateAuthorizationsImpl(ctx sdk.Context, prefix []byte, fn func(classID string, authorization token.Authorization) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		classID, proxy, approver := splitAuthorizationKey(iterator.Key())
		authorization := token.Authorization{
			Approver: approver.String(),
			Proxy:    proxy.String(),
		}

		stop := fn(classID, authorization)
		if stop {
			break
		}
	}
}

func (k Keeper) iterateStatistics(ctx sdk.Context, prefix []byte, fn func(classID string, amount sdk.Int) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var amount sdk.Int
		if err := amount.Unmarshal(iterator.Value()); err != nil {
			panic(err)
		}

		classID := splitStatisticsKey(iterator.Key(), prefix)

		stop := fn(classID, amount)
		if stop {
			break
		}
	}
}

func (k Keeper) iterateSupplies(ctx sdk.Context, fn func(classID string, amount sdk.Int) (stop bool)) {
	k.iterateStatistics(ctx, supplyKeyPrefix, fn)
}

func (k Keeper) iterateMinteds(ctx sdk.Context, fn func(classID string, amount sdk.Int) (stop bool)) {
	k.iterateStatistics(ctx, mintKeyPrefix, fn)
}

func (k Keeper) iterateBurnts(ctx sdk.Context, fn func(classID string, amount sdk.Int) (stop bool)) {
	k.iterateStatistics(ctx, burnKeyPrefix, fn)
}
