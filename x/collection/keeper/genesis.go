package keeper

import (
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/collection"
)

// InitGenesis new collection genesis
func (k Keeper) InitGenesis(ctx sdk.Context, data *collection.GenesisState) {
	for _, contractBalances := range data.Balances {
		for _, balance := range contractBalances.Balances {
			for _, coin := range balance.Amount {
				k.setBalance(ctx, contractBalances.ContractId, sdk.AccAddress(balance.Address), coin.TokenId, coin.Amount)
			}
		}
	}

	// for _, class := range data.Classes {
	// 	if err := k.setClass(ctx, class); err != nil {
	// 		panic(err)
	// 	}
	// }

	// for _, contractGrants := range data.Grants {
	// 	for _, grant := range contractGrants.Grants {
	// 		permission := collection.Permission(collection.Permission_value[grant.Permission])
	// 		k.setGrant(ctx, contractGrants.ContractId, sdk.AccAddress(grant.Grantee), permission)
	// 	}
	// }

	// for _, contractAuthorizations := range data.Authorizations {
	// 	for _, authorization := range contractAuthorizations.Authorizations {
	// 		k.setAuthorization(ctx, contractAuthorizations.ContractId, sdk.AccAddress(authorization.Approver), sdk.AccAddress(authorization.Proxy))
	// 	}
	// }

	// // TODO: remove it (derive it using mints and burns)
	// for _, amount := range data.Supplies {
	// 	if err := k.setSupply(ctx, amount.ContractId, amount.Amount); err != nil {
	// 		panic(err)
	// 	}
	// }

	// for _, amount := range data.Mints {
	// 	if err := k.setMinted(ctx, amount.ContractId, amount.Amount); err != nil {
	// 		panic(err)
	// 	}
	// }

	// for _, amount := range data.Burns {
	// 	if err := k.setBurnt(ctx, amount.ContractId, amount.Amount); err != nil {
	// 		panic(err)
	// 	}
	// }
}

// ExportGenesis returns a GenesisState for a given context.
func (k Keeper) ExportGenesis(ctx sdk.Context) *collection.GenesisState {
	// var classes []collection.TokenClass
	// k.iterateClasses(ctx, func(class collection.TokenClass) (stop bool) {
	// 	classes = append(classes, class)
	// 	return false
	// })

	// var balances []collection.ContractBalances
	// for _, class := range classes {
	// 	id := class.ContractId
	// 	contractBalances := collection.ContractBalances{
	// 		ContractId: id,
	// 	}

	// 	k.iterateContractBalances(ctx, id, func(balance collection.Balance) (stop bool) {
	// 		contractBalances.Balances = append(contractBalances.Balances, balance)
	// 		return false
	// 	})
	// 	if len(contractBalances.Balances) != 0 {
	// 		balances = append(balances, contractBalances)
	// 	}
	// }

	// var grants []collection.ContractGrants
	// for _, class := range classes {
	// 	id := class.ContractId
	// 	contractGrants := collection.ContractGrants{
	// 		ContractId: id,
	// 	}

	// 	k.iterateContractGrants(ctx, id, func(grant collection.Grant) (stop bool) {
	// 		contractGrants.Grants = append(contractGrants.Grants, grant)
	// 		return false
	// 	})
	// 	if len(contractGrants.Grants) != 0 {
	// 		grants = append(grants, contractGrants)
	// 	}
	// }

	// var authorizations []collection.ContractAuthorizations
	// for _, class := range classes {
	// 	id := class.ContractId
	// 	contractAuthorizations := collection.ContractAuthorizations{
	// 		ContractId: id,
	// 	}

	// 	k.iterateContractAuthorizations(ctx, id, func(authorization collection.Authorization) (stop bool) {
	// 		contractAuthorizations.Authorizations = append(contractAuthorizations.Authorizations, authorization)
	// 		return false
	// 	})
	// 	if len(contractAuthorizations.Authorizations) != 0 {
	// 		authorizations = append(authorizations, contractAuthorizations)
	// 	}
	// }

	return &collection.GenesisState{
		// Balances:       balances,
		// Classes:        classes,
		// Grants:         grants,
		// Authorizations: authorizations,
		// Supplies:       supplies,
		// Mints:          mints,
		// Burns:          burns,
	}
}
