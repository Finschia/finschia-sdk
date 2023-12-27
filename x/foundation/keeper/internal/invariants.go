package internal

import (
	"fmt"

	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/foundation"
)

const (
	moduleAccountInvariant = "module-accounts"
	totalWeightInvariant   = "total-weight"
	proposalInvariant      = "proposals"
)

func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	for name, invariant := range map[string]func(k Keeper) sdk.Invariant{
		moduleAccountInvariant: ModuleAccountInvariant,
		totalWeightInvariant:   TotalWeightInvariant,
		proposalInvariant:      ProposalInvariant,
	} {
		ir.RegisterRoute(foundation.ModuleName, name, invariant(k))
	}
}

func ModuleAccountInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		treasuryAcc := k.authKeeper.GetModuleAccount(ctx, foundation.TreasuryName)
		balance := k.bankKeeper.GetAllBalances(ctx, treasuryAcc.GetAddress())

		treasury := k.GetTreasury(ctx)
		msg := fmt.Sprintf("coins in the treasury; expected %s, got %s\n", treasury, balance)
		broken := !treasury.IsEqual(sdk.NewDecCoinsFromCoins(balance...))

		return sdk.FormatInvariant(foundation.ModuleName, moduleAccountInvariant, msg), broken
	}
}

func TotalWeightInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		expected := k.GetFoundationInfo(ctx).TotalWeight
		actual := sdk.NewDec(int64(len(k.GetMembers(ctx))))

		msg := fmt.Sprintf("total weight of foundation; expected %s, got %s\n", expected, actual)
		broken := !actual.Equal(expected)

		return sdk.FormatInvariant(foundation.ModuleName, totalWeightInvariant, msg), broken
	}
}

func ProposalInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		version := k.GetFoundationInfo(ctx).Version
		msg := fmt.Sprintf("current foundation version; %d\n", version)
		broken := false

		k.iterateProposals(ctx, func(proposal foundation.Proposal) (stop bool) {
			if proposal.FoundationVersion == version {
				return true
			}

			if proposal.Status == foundation.PROPOSAL_STATUS_SUBMITTED {
				msg += fmt.Sprintf("active old proposal %d found; version %d\n", proposal.Id, proposal.FoundationVersion)
				broken = true
			}

			return false
		})

		return sdk.FormatInvariant(foundation.ModuleName, proposalInvariant, msg), broken
	}
}
