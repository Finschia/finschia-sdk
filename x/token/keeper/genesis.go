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

	for _, balance := range data.Balances {
		if err := k.addTokens(ctx, sdk.AccAddress(balance.Address), balance.Tokens); err != nil {
			panic(err)
		}
	}

	for _, class := range data.Classes {
		if err := k.setClass(ctx, class); err != nil {
			panic(err)
		}
	}

	for _, grant := range data.Grants {
		k.setGrant(ctx, sdk.AccAddress(grant.Grantee), grant.ClassId, grant.Action, true)
	}

	for _, approve := range data.Approves {
		k.setApprove(ctx, sdk.AccAddress(approve.Approver), sdk.AccAddress(approve.Proxy), approve.ClassId, true)
	}

	for _, amount := range data.Supplies {
		if err := k.setSupply(ctx, amount); err != nil {
			panic(err)
		}
	}

	for _, amount := range data.Mints {
		if err := k.setMint(ctx, amount); err != nil {
			panic(err)
		}
	}

	for _, amount := range data.Burns {
		if err := k.setBurn(ctx, amount); err != nil {
			panic(err)
		}
	}
}

// ExportGenesis returns a GenesisState for a given context.
func (k Keeper) ExportGenesis(ctx sdk.Context) *token.GenesisState {
	var balances []token.Balance
	var addrToIndex map[string]int
	k.iterateBalances(ctx, func(addr sdk.AccAddress, amount token.FT) (stop bool) {
		if index, ok := addrToIndex[addr.String()]; ok {
			balances[index].Tokens = append(balances[index].Tokens, amount)
			return false
		}

		balance := token.Balance{
			Address: addr.String(),
			Tokens: []token.FT{amount},
		}
		balances = append(balances, balance)
		addrToIndex[addr.String()] = len(balances) - 1
		return false
	})

	var classes []token.Token
	k.iterateClasses(ctx, func(class token.Token) (stop bool) {
		classes = append(classes, class)
		return false
	})

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

	var grants []token.Grant
	k.iterateGrants(ctx, func(grant token.Grant) (stop bool) {
		grants = append(grants, grant)
		return false
	})

	var approves []token.Approve
	k.iterateApproves(ctx, func(approve token.Approve) (stop bool) {
		approves = append(approves, approve)
		return false
	})
	
	return &token.GenesisState{
		ClassState: k.classKeeper.ExportGenesis(ctx),
		Balances: balances,
		Classes: classes,
		Grants: grants,
		Approves: approves,
		Supplies: supplies,
		Mints: mints,
		Burns: burns,
	}
}
