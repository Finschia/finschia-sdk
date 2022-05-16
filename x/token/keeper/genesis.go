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
			k.setGrant(ctx, contractGrants.ContractId, sdk.AccAddress(grant.Grantee), grant.Permission)
		}
	}

	for _, contractAuthorizations := range data.Authorizations {
		for _, authorization := range contractAuthorizations.Authorizations {
			k.setAuthorization(ctx, contractAuthorizations.ContractId, sdk.AccAddress(authorization.Approver), sdk.AccAddress(authorization.Proxy))
		}
	}

	// TODO: remove it (derive it using mints and burns)
	for _, amount := range data.Supplies {
		if err := k.setSupply(ctx, amount.ContractId, amount.Amount); err != nil {
			panic(err)
		}
	}

	for _, amount := range data.Mints {
		if err := k.setMint(ctx, amount.ContractId, amount.Amount); err != nil {
			panic(err)
		}
	}

	for _, amount := range data.Burns {
		if err := k.setBurn(ctx, amount.ContractId, amount.Amount); err != nil {
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

	balances := make([]token.ContractBalances, len(classes))
	for i, class := range classes {
		id := class.ContractId
		balances[i].ContractId = id
		k.iterateContractBalances(ctx, id, func(contractID string, balance token.Balance) (stop bool) {
			balances[i].Balances = append(balances[i].Balances, balance)
			return false
		})
	}

	var supplies []token.FT
	k.iterateSupplies(ctx, func(amount token.FT) (stop bool) {
		supplies = append(supplies, amount)
		return false
	})

	var mints []token.FT
	k.iterateMints(ctx, func(amount token.FT) (stop bool) {
		mints = append(mints, amount)
		return false
	})

	var burns []token.FT
	k.iterateBurns(ctx, func(amount token.FT) (stop bool) {
		burns = append(burns, amount)
		return false
	})

	grants := make([]token.ContractGrants, len(classes))
	for i, class := range classes {
		id := class.ContractId
		grants[i].ContractId = id
		k.iterateContractGrants(ctx, id, func(contractID string, grant token.Grant) (stop bool) {
			grants[i].Grants = append(grants[i].Grants, grant)
			return false
		})
	}

	authorizations := make([]token.ContractAuthorizations, len(classes))
	for i, class := range classes {
		id := class.ContractId
		authorizations[i].ContractId = id

		k.iterateContractAuthorizations(ctx, id, func(contractID string, authorization token.Authorization) (stop bool) {
			authorizations[i].Authorizations = append(authorizations[i].Authorizations, authorization)
			return false
		})
	}

	return &token.GenesisState{
		ClassState: k.classKeeper.ExportGenesis(ctx),
		Balances:   balances,
		Classes:    classes,
		Grants:     grants,
		Authorizations:   authorizations,
		Supplies:   supplies,
		Mints:      mints,
		Burns:      burns,
	}
}
