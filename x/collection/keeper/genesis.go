package keeper

import (
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/collection"
)

// InitGenesis new collection genesis
func (k Keeper) InitGenesis(ctx sdk.Context, data *collection.GenesisState) {
	for _, contract := range data.Contracts {
		k.setContract(ctx, contract)
	}

	for _, contractClasses := range data.Classes {
		contractID := contractClasses.ContractId

		for _, any := range contractClasses.Classes {
			class := collection.TokenClassFromAny(any)
			k.setTokenClass(ctx, contractID, class)
		}
	}

	for _, contractBalances := range data.Balances {
		contractID := contractBalances.ContractId

		for _, balance := range contractBalances.Balances {
			for _, coin := range balance.Amount {
				k.setBalance(ctx, contractID, sdk.AccAddress(balance.Address), coin.TokenId, coin.Amount)

				if err := collection.ValidateNFTID(coin.TokenId); err == nil {
					k.setOwner(ctx, contractID, coin.TokenId, sdk.AccAddress(balance.Address))
				}
			}
		}
	}

	for _, contractParents := range data.Parents {
		contractID := contractParents.ContractId

		for _, relation := range contractParents.Relations {
			tokenID := relation.Self
			parentID := relation.Other
			k.setParent(ctx, contractID, tokenID, parentID)
			k.setChild(ctx, contractID, parentID, tokenID)
		}
	}

	for _, contractAuthorizations := range data.Authorizations {
		for _, authorization := range contractAuthorizations.Authorizations {
			k.setAuthorization(ctx, contractAuthorizations.ContractId, sdk.AccAddress(authorization.Holder), sdk.AccAddress(authorization.Operator))
		}
	}

	for _, contractGrants := range data.Grants {
		for _, grant := range contractGrants.Grants {
			permission := collection.Permission(collection.Permission_value[grant.Permission])
			k.setGrant(ctx, contractGrants.ContractId, sdk.AccAddress(grant.Grantee), permission)
		}
	}

	for _, contractBurnts := range data.Burnts {
		contractID := contractBurnts.ContractId
		for _, burnt := range contractBurnts.Statistics {
			k.setBurnt(ctx, contractID, burnt.ClassId, burnt.Amount)
		}
	}

	for _, contractSupplies := range data.Supplies {
		contractID := contractSupplies.ContractId
		for _, supply := range contractSupplies.Statistics {
			k.setSupply(ctx, contractID, supply.ClassId, supply.Amount)

			// calculate the amount of minted tokens
			burnt := k.GetBurnt(ctx, contractID, supply.ClassId)
			minted := supply.Amount.Add(burnt)
			k.setMinted(ctx, contractID, supply.ClassId, minted)
		}
	}
}

// ExportGenesis returns a GenesisState for a given context.
func (k Keeper) ExportGenesis(ctx sdk.Context) *collection.GenesisState {
	contracts := k.getContracts(ctx)

	return &collection.GenesisState{
		Contracts:      contracts,
		Classes:        k.getClasses(ctx, contracts),
		Balances:       k.getBalances(ctx, contracts),
		Parents:        k.getParents(ctx, contracts),
		Grants:         k.getGrants(ctx, contracts),
		Authorizations: k.getAuthorizations(ctx, contracts),
		Supplies:       k.getSupplies(ctx, contracts),
		Burnts:         k.getBurnts(ctx, contracts),
	}
}

func (k Keeper) getContracts(ctx sdk.Context) []collection.Contract {
	var contracts []collection.Contract
	k.iterateContracts(ctx, func(contract collection.Contract) (stop bool) {
		contracts = append(contracts, contract)
		return false
	})

	return contracts
}

func (k Keeper) getClasses(ctx sdk.Context, contracts []collection.Contract) []collection.ContractClasses {
	var classes []collection.ContractClasses
	for _, contract := range contracts {
		contractID := contract.ContractId
		contractClasses := collection.ContractClasses{
			ContractId: contractID,
		}

		k.iterateContractClasses(ctx, contractID, func(class collection.TokenClass) (stop bool) {
			any := collection.TokenClassToAny(class)
			contractClasses.Classes = append(contractClasses.Classes, any)
			return false
		})
		if len(contractClasses.Classes) != 0 {
			classes = append(classes, contractClasses)
		}
	}

	return classes
}

func (k Keeper) getBalances(ctx sdk.Context, contracts []collection.Contract) []collection.ContractBalances {
	var balances []collection.ContractBalances
	for _, contract := range contracts {
		contractID := contract.ContractId
		contractBalances := collection.ContractBalances{
			ContractId: contractID,
		}

		contractBalances.Balances = k.getContractBalances(ctx, contractID)
		if len(contractBalances.Balances) != 0 {
			balances = append(balances, contractBalances)
		}
	}

	return balances
}

func (k Keeper) getContractBalances(ctx sdk.Context, contractID string) []collection.Balance {
	var balances []collection.Balance
	addressToBalanceIndex := make(map[sdk.AccAddress]int)

	k.iterateContractBalances(ctx, contractID, func(address sdk.AccAddress, balance collection.Coin) (stop bool) {
		index, ok := addressToBalanceIndex[address]
		if ok {
			balances[index].Amount = append(balances[index].Amount, balance)
			return false
		}

		accountBalance := collection.Balance{
			Address: address.String(),
			Amount:  collection.Coins{},
		}
		balances = append(balances, accountBalance)
		addressToBalanceIndex[address] = len(balances) - 1
		return false
	})

	return balances
}

func (k Keeper) getParents(ctx sdk.Context, contracts []collection.Contract) []collection.ContractTokenRelations {
	var parents []collection.ContractTokenRelations
	for _, contract := range contracts {
		contractID := contract.ContractId
		contractParents := collection.ContractTokenRelations{
			ContractId: contractID,
		}

		k.iterateContractParents(ctx, contractID, func(tokenID, parentID string) (stop bool) {
			relation := collection.TokenRelation{
				Self:  tokenID,
				Other: parentID,
			}
			contractParents.Relations = append(contractParents.Relations, relation)
			return false
		})
		if len(contractParents.Relations) != 0 {
			parents = append(parents, contractParents)
		}
	}

	return parents
}

func (k Keeper) getAuthorizations(ctx sdk.Context, contracts []collection.Contract) []collection.ContractAuthorizations {
	var authorizations []collection.ContractAuthorizations
	for _, contract := range contracts {
		contractID := contract.ContractId
		contractAuthorizations := collection.ContractAuthorizations{
			ContractId: contractID,
		}

		k.iterateContractAuthorizations(ctx, contractID, func(authorization collection.Authorization) (stop bool) {
			contractAuthorizations.Authorizations = append(contractAuthorizations.Authorizations, authorization)
			return false
		})
		if len(contractAuthorizations.Authorizations) != 0 {
			authorizations = append(authorizations, contractAuthorizations)
		}
	}

	return authorizations
}

func (k Keeper) getGrants(ctx sdk.Context, contracts []collection.Contract) []collection.ContractGrants {
	var grants []collection.ContractGrants
	for _, contract := range contracts {
		contractID := contract.ContractId
		contractGrants := collection.ContractGrants{
			ContractId: contractID,
		}

		k.iterateContractGrants(ctx, contractID, func(grant collection.Grant) (stop bool) {
			contractGrants.Grants = append(contractGrants.Grants, grant)
			return false
		})
		if len(contractGrants.Grants) != 0 {
			grants = append(grants, contractGrants)
		}
	}

	return grants
}

func (k Keeper) getSupplies(ctx sdk.Context, contracts []collection.Contract) []collection.ContractStatistics {
	return k.getStatistics(ctx, contracts, k.iterateContractSupplies)
}

func (k Keeper) getBurnts(ctx sdk.Context, contracts []collection.Contract) []collection.ContractStatistics {
	return k.getStatistics(ctx, contracts, k.iterateContractBurnts)
}

func (k Keeper) getStatistics(ctx sdk.Context, contracts []collection.Contract, iterator func(ctx sdk.Context, contractID string, cb func(classID string, amount sdk.Int) (stop bool))) []collection.ContractStatistics {
	var statistics []collection.ContractStatistics
	for _, contract := range contracts {
		contractID := contract.ContractId
		contractStatistics := collection.ContractStatistics{
			ContractId: contractID,
		}

		iterator(ctx, contractID, func(classID string, amount sdk.Int) (stop bool) {
			supply := collection.ClassStatistics{
				ClassId: classID,
				Amount:  amount,
			}
			contractStatistics.Statistics = append(contractStatistics.Statistics, supply)
			return false
		})
		if len(contractStatistics.Statistics) != 0 {
			statistics = append(statistics, contractStatistics)
		}
	}

	return statistics
}
