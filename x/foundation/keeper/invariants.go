package keeper

import (
	"fmt"

	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/foundation"
)

const (
	moduleAccountInvariant = "module-accounts"
	proposalInvariant      = "proposals"
)

func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	ir.RegisterRoute(foundation.ModuleName, moduleAccountInvariant, ModuleAccountInvariant(k))
	ir.RegisterRoute(foundation.ModuleName, proposalInvariant, ProposalInvariant(k))
}

func ModuleAccountInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		// cache, we don't want to write changes
		ctx, _ = ctx.CacheContext()

		treasuryAcc := k.authKeeper.GetModuleAccount(ctx, foundation.TreasuryName)
		balance := k.bankKeeper.GetAllBalances(ctx, treasuryAcc.GetAddress())

		treasury := k.GetTreasury(ctx)
		msg := fmt.Sprintf("coins in the treasury; expected %s, got %s\n", treasury, balance)
		broken := !treasury.IsEqual(sdk.NewDecCoinsFromCoins(balance...))

		return sdk.FormatInvariant(foundation.ModuleName, moduleAccountInvariant, msg), broken
	}
}

func ProposalInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		// cache, we don't want to write changes
		ctx, _ = ctx.CacheContext()

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
