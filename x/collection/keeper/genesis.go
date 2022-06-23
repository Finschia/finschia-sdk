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
	var contracts []collection.Contract
	k.iterateContracts(ctx, func(contract collection.Contract) (stop bool) {
		contracts = append(contracts, contract)
		return false
	})

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
	}

	var balances []collection.ContractBalances
	for _, contract := range contracts {
		contractID := contract.ContractId
		contractBalances := collection.ContractBalances{
			ContractId: contractID,
			Balances:   k.getContractBalances(ctx, contractID),
		}

		if len(contractBalances.Balances) != 0 {
			balances = append(balances, contractBalances)
		}
	}

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

	var supplies []collection.ContractStatistics
	for _, contract := range contracts {
		contractID := contract.ContractId
		contractSupplies := collection.ContractStatistics{
			ContractId: contractID,
		}

		k.iterateContractSupplies(ctx, contractID, func(classID string, amount sdk.Int) (stop bool) {
			supply := collection.ClassStatistics{
				ClassId: classID,
				Amount:  amount,
			}
			contractSupplies.Statistics = append(contractSupplies.Statistics, supply)
			return false
		})
	}

	var burnts []collection.ContractStatistics
	for _, contract := range contracts {
		contractID := contract.ContractId
		contractBurnts := collection.ContractStatistics{
			ContractId: contractID,
		}

		k.iterateContractBurnts(ctx, contractID, func(classID string, amount sdk.Int) (stop bool) {
			burnt := collection.ClassStatistics{
				ClassId: classID,
				Amount:  amount,
			}
			contractBurnts.Statistics = append(contractBurnts.Statistics, burnt)
			return false
		})
	}

	return &collection.GenesisState{
		Contracts:      contracts,
		Classes:        classes,
		Balances:       balances,
		Parents:        parents,
		Grants:         grants,
		Authorizations: authorizations,
		Supplies:       supplies,
		Burnts:         burnts,
	}
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
