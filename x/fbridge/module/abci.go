package module

import (
	"fmt"
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/fbridge/keeper"
	"github.com/Finschia/finschia-sdk/x/fbridge/types"
)

func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	k.InitMemStore(ctx)

	proposals := k.GetRoleProposals(ctx)
	for _, proposal := range proposals {
		if ctx.BlockTime().After(proposal.ExpiredAt) {
			k.DeleteRoleProposal(ctx, proposal.Id)
		}
	}
}

func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
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

		var total uint64 = 0
		roleMeta := k.GetRoleMetadata(ctx)
		switch proposal.Role {
		case types.RoleGuardian:
			total = roleMeta.Guardian
		case types.RoleOperator:
			total = roleMeta.Operator
		case types.RoleJudge:
			total = roleMeta.Judge
		default:
			panic(fmt.Sprintf("invalid role: %s\n", proposal.Role))
		}

		if types.CheckTrustLevelThreshold(total, voteYes, guardianTrustLevel) {
			if err := k.UpdateRole(ctx, proposal.Role, sdk.MustAccAddressFromBech32(proposal.Target)); err != nil {
				panic(err)
			}

			k.DeleteRoleProposal(ctx, proposal.Id)
		}
	}
}

// RegisterInvariants registers all fbridge invariants
func RegisterInvariants(ir sdk.InvariantRegistry, k keeper.Keeper) {
	ir.RegisterRoute(types.ModuleName, "role-metadata", RoleMeatadataInvariant(k))
}

func RoleMeatadataInvariant(k keeper.Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		actualRoleMeta := k.GetRoleMetadata(ctx)
		expectedRoleMeta := types.RoleMetadata{}
		for _, pair := range k.GetRolePairs(ctx) {
			switch pair.Role {
			case types.RoleGuardian:
				expectedRoleMeta.Guardian++
			case types.RoleOperator:
				expectedRoleMeta.Operator++
			case types.RoleJudge:
				expectedRoleMeta.Judge++
			}
		}

		broken := expectedRoleMeta.Guardian != actualRoleMeta.Guardian ||
			expectedRoleMeta.Operator != actualRoleMeta.Operator ||
			expectedRoleMeta.Judge != actualRoleMeta.Judge

		return sdk.FormatInvariant(types.ModuleName, "registered members and role metadata", fmt.Sprintf(
			"Saved Role Metadata: %+v"+
				"Calculated Role Metadata: %+v",
			actualRoleMeta, expectedRoleMeta)), broken
	}
}
