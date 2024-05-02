package module

import (
	"fmt"
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/fbridge/keeper"
	"github.com/Finschia/finschia-sdk/x/fbridge/types"
)

func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	guardianTrustLevel := k.GetParams(ctx).GuardianTrustLevel
	proposals := k.GetRoleProposals(ctx)
	for _, proposal := range proposals {
		if ctx.BlockTime().After(proposal.ExpiredAt) {
			k.DeleteRoleProposal(ctx, proposal.Id)
			continue
		}

		votes := k.GetProposalVotes(ctx, proposal.Id)

		voteYes := 0
		for _, vote := range votes {
			if vote.Option == types.OptionYes {
				voteYes++
			}
		}

		var total uint32 = 0
		roleMeta := k.GetRoleMetadata(ctx)
		previousRole := k.GetRole(ctx, sdk.MustAccAddressFromBech32(proposal.Target))
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

		if types.CheckTrustLevelThreshold(uint64(total), uint64(voteYes), guardianTrustLevel) {
			if proposal.Role == types.RoleEmpty {
				k.DeleteRole(ctx, sdk.MustAccAddressFromBech32(proposal.Target))
			} else {
				k.SetRole(ctx, proposal.Role, sdk.MustAccAddressFromBech32(proposal.Target))
			}

			switch proposal.Role {
			case types.RoleGuardian:
				roleMeta.Guardian++
			case types.RoleOperator:
				roleMeta.Operator++
			case types.RoleJudge:
				roleMeta.Judge++
			}
			switch previousRole {
			case types.RoleGuardian:
				roleMeta.Guardian--
			case types.RoleOperator:
				roleMeta.Operator--
			case types.RoleJudge:
				roleMeta.Judge--
			}
			k.SetRoleMetadata(ctx, roleMeta)

			k.DeleteRoleProposal(ctx, proposal.Id)
		}
	}
}
