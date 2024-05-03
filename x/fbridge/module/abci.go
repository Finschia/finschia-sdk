package module

import (
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

		if types.CheckTrustLevelThreshold(k.GetRoleMetadata(ctx).Guardian, voteYes, guardianTrustLevel) || proposal.Proposer == k.GetAuthority() {
			if err := k.UpdateRole(ctx, proposal.Role, sdk.MustAccAddressFromBech32(proposal.Target)); err != nil {
				panic(err)
			}

			k.DeleteRoleProposal(ctx, proposal.Id)
		}
	}
}
