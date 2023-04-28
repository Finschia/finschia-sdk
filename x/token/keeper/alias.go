package keeper

import (
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/token"
)

// iterate through the balances of a contract and perform the provided function
func (k Keeper) iterateContractBalances(ctx sdk.Context, contractID string, fn func(balance token.Balance) (stop bool)) {
	k.iterateBalancesImpl(ctx, balanceKeyPrefixByContractID(contractID), func(_ string, balance token.Balance) (stop bool) {
		return fn(balance)
	})
}

func (k Keeper) iterateBalancesImpl(ctx sdk.Context, prefix []byte, fn func(contractID string, balance token.Balance) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		contractID, addr := splitBalanceKey(iterator.Key())

		var amount sdk.Int
		if err := amount.Unmarshal(iterator.Value()); err != nil {
			panic(err)
		}
		balance := token.Balance{
			Address: addr.String(),
			Amount:  amount,
		}

		stop := fn(contractID, balance)
		if stop {
			break
		}
	}
}

// iterate through the classes and perform the provided function
func (k Keeper) iterateClasses(ctx sdk.Context, fn func(class token.Contract) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, classKeyPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var class token.Contract
		k.cdc.MustUnmarshal(iterator.Value(), &class)

		stop := fn(class)
		if stop {
			break
		}
	}
}

func (k Keeper) iterateContractGrants(ctx sdk.Context, contractID string, fn func(grant token.Grant) (stop bool)) {
	k.iterateGrantsImpl(ctx, grantKeyPrefixByContractID(contractID), func(_ string, grant token.Grant) (stop bool) {
		return fn(grant)
	})
}

func (k Keeper) iterateGrantsImpl(ctx sdk.Context, prefix []byte, fn func(contractID string, grant token.Grant) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		contractID, grantee, permission := splitGrantKey(iterator.Key())
		grant := token.Grant{
			Grantee:    grantee.String(),
			Permission: permission,
		}

		stop := fn(contractID, grant)
		if stop {
			break
		}
	}
}

func (k Keeper) iterateContractAuthorizations(ctx sdk.Context, contractID string, fn func(authorization token.Authorization) (stop bool)) {
	k.iterateAuthorizationsImpl(ctx, authorizationKeyPrefixByContractID(contractID), func(_ string, authorization token.Authorization) (stop bool) {
		return fn(authorization)
	})
}

func (k Keeper) iterateAuthorizationsImpl(ctx sdk.Context, prefix []byte, fn func(contractID string, authorization token.Authorization) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		contractID, operator, holder := splitAuthorizationKey(iterator.Key())
		authorization := token.Authorization{
			Holder:   holder.String(),
			Operator: operator.String(),
		}

		stop := fn(contractID, authorization)
		if stop {
			break
		}
	}
}

func (k Keeper) iterateStatistics(ctx sdk.Context, prefix []byte, fn func(contractID string, amount sdk.Int) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var amount sdk.Int
		if err := amount.Unmarshal(iterator.Value()); err != nil {
			panic(err)
		}

		contractID := splitStatisticsKey(iterator.Key(), prefix)

		stop := fn(contractID, amount)
		if stop {
			break
		}
	}
}

func (k Keeper) iterateSupplies(ctx sdk.Context, fn func(contractID string, amount sdk.Int) (stop bool)) {
	k.iterateStatistics(ctx, supplyKeyPrefix, fn)
}

func (k Keeper) iterateMinteds(ctx sdk.Context, fn func(contractID string, amount sdk.Int) (stop bool)) {
	k.iterateStatistics(ctx, mintKeyPrefix, fn)
}

func (k Keeper) iterateBurnts(ctx sdk.Context, fn func(contractID string, amount sdk.Int) (stop bool)) {
	k.iterateStatistics(ctx, burnKeyPrefix, fn)
}
