package keeper

import (
	"fmt"

	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/fbridge/types"
)

func (k Keeper) BeginBlocker(ctx sdk.Context) {
	k.InitMemStore(ctx)

	proposals := k.GetRoleProposals(ctx)
	for _, proposal := range proposals {
		if ctx.BlockTime().After(proposal.ExpiredAt) {
			if err := k.deleteRoleProposal(ctx, proposal.Id); err != nil {
				panic(err)
			}
		}
	}
}

func (k Keeper) EndBlocker(ctx sdk.Context) {
	guardianTrustLevel := k.GetParams(ctx).GuardianTrustLevel
	proposals := k.GetRoleProposals(ctx)
	for _, proposal := range proposals {
		votes := k.GetProposalVotes(ctx, proposal.Id)

		var voteYes uint64 = 0
		for _, vote := range votes {
			if vote.Option == types.OptionYes {
				voteYes++
			}
		}

		if types.CheckTrustLevelThreshold(k.GetRoleMetadata(ctx).Guardian, voteYes, guardianTrustLevel) || proposal.Proposer == k.GetAuthority() {
			if err := k.updateRole(ctx, proposal.Role, sdk.MustAccAddressFromBech32(proposal.Target)); err != nil {
				panic(err)
			}

			if err := k.deleteRoleProposal(ctx, proposal.Id); err != nil {
				panic(err)
			}
		}
	}
}

// RegisterInvariants registers the fbridge module invariants
func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	ir.RegisterRoute(types.ModuleName, "guardian-invariant", GuardianInvariant(k))
}

func GuardianInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		numGuardian := 0
		for _, p := range k.GetRolePairs(ctx) {
			if p.Role == types.RoleGuardian {
				numGuardian++
			}
		}

		numBridgeSw := len(k.GetBridgeSwitches(ctx))
		broken := numGuardian != numBridgeSw
		return sdk.FormatInvariant(
			types.ModuleName, "guardian-invariant",
			fmt.Sprintf("number of guardians(%d) != number of bridge switches(%d)", numGuardian, numBridgeSw),
		), broken
	}
}
