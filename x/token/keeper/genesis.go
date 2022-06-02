package keeper

import (
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/token"
)

// InitGenesis new token genesis
func (k Keeper) InitGenesis(ctx sdk.Context, data *token.GenesisState) {
	if data.ClassState == nil {
		data.ClassState = token.DefaultClassGenesisState()
	}
	k.classKeeper.InitGenesis(ctx, data.ClassState)

	for _, contractBalances := range data.Balances {
		for _, balance := range contractBalances.Balances {
			if err := k.setBalance(ctx, contractBalances.ContractId, sdk.AccAddress(balance.Address), balance.Amount); err != nil {
				panic(err)
			}
		}
	}

	for _, class := range data.Classes {
		if err := k.setClass(ctx, class); err != nil {
			panic(err)
		}
	}

	for _, contractGrants := range data.Grants {
		for _, grant := range contractGrants.Grants {
			permission := token.Permission(token.Permission_value[grant.Permission])
			k.setGrant(ctx, contractGrants.ContractId, sdk.AccAddress(grant.Grantee), permission)
		}
	}

	for _, contractAuthorizations := range data.Authorizations {
		for _, authorization := range contractAuthorizations.Authorizations {
			k.setAuthorization(ctx, contractAuthorizations.ContractId, sdk.AccAddress(authorization.Holder), sdk.AccAddress(authorization.Operator))
		}
	}

	// TODO: remove it (derive it using mints and burns)
	for _, amount := range data.Supplies {
		if err := k.setSupply(ctx, amount.ContractId, amount.Amount); err != nil {
			panic(err)
		}
	}

	for _, amount := range data.Mints {
		if err := k.setMinted(ctx, amount.ContractId, amount.Amount); err != nil {
			panic(err)
		}
	}

	for _, amount := range data.Burns {
		if err := k.setBurnt(ctx, amount.ContractId, amount.Amount); err != nil {
			panic(err)
		}
	}
}

// ExportGenesis returns a GenesisState for a given context.
func (k Keeper) ExportGenesis(ctx sdk.Context) *token.GenesisState {
	var classes []token.TokenClass
	k.iterateClasses(ctx, func(class token.TokenClass) (stop bool) {
		classes = append(classes, class)
		return false
	})

	var balances []token.ContractBalances
	for _, class := range classes {
		id := class.ContractId
		contractBalances := token.ContractBalances{
			ContractId: id,
		}

		k.iterateContractBalances(ctx, id, func(balance token.Balance) (stop bool) {
			contractBalances.Balances = append(contractBalances.Balances, balance)
			return false
		})
		if len(contractBalances.Balances) != 0 {
			balances = append(balances, contractBalances)
		}
	}

	var supplies []token.ContractCoin
	k.iterateSupplies(ctx, func(contractID string, amount sdk.Int) (stop bool) {
		supply := token.ContractCoin{
			ContractId: contractID,
			Amount:     amount,
		}
		supplies = append(supplies, supply)
		return false
	})

	var mints []token.ContractCoin
	k.iterateMinteds(ctx, func(contractID string, amount sdk.Int) (stop bool) {
		minted := token.ContractCoin{
			ContractId: contractID,
			Amount:     amount,
		}
		mints = append(mints, minted)
		return false
	})

	var burns []token.ContractCoin
	k.iterateBurnts(ctx, func(contractID string, amount sdk.Int) (stop bool) {
		burnt := token.ContractCoin{
			ContractId: contractID,
			Amount:     amount,
		}
		burns = append(burns, burnt)
		return false
	})

	var grants []token.ContractGrants
	for _, class := range classes {
		id := class.ContractId
		contractGrants := token.ContractGrants{
			ContractId: id,
		}

		k.iterateContractGrants(ctx, id, func(grant token.Grant) (stop bool) {
			contractGrants.Grants = append(contractGrants.Grants, grant)
			return false
		})
		if len(contractGrants.Grants) != 0 {
			grants = append(grants, contractGrants)
		}
	}

	var authorizations []token.ContractAuthorizations
	for _, class := range classes {
		id := class.ContractId
		contractAuthorizations := token.ContractAuthorizations{
			ContractId: id,
		}

		k.iterateContractAuthorizations(ctx, id, func(authorization token.Authorization) (stop bool) {
			contractAuthorizations.Authorizations = append(contractAuthorizations.Authorizations, authorization)
			return false
		})
		if len(contractAuthorizations.Authorizations) != 0 {
			authorizations = append(authorizations, contractAuthorizations)
		}
	}

	return &token.GenesisState{
		ClassState:     k.classKeeper.ExportGenesis(ctx),
		Balances:       balances,
		Classes:        classes,
		Grants:         grants,
		Authorizations: authorizations,
		Supplies:       supplies,
		Mints:          mints,
		Burns:          burns,
	}
}
