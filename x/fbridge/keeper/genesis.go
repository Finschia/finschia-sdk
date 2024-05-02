package keeper

import (
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/fbridge/types"
)

func (k Keeper) InitGenesis(ctx sdk.Context, gs *types.GenesisState) error {
	k.SetParams(ctx, gs.Params)
	k.setNextSequence(ctx, gs.SendingState.NextSeq)
	k.setNextProposalID(ctx, gs.NextRoleProposalId)
	k.setRoleMetadata(ctx, gs.RoleMetadata)

	for _, proposal := range gs.RoleProposals {
		k.setRoleProposal(ctx, proposal)
	}
	// TODO: we initialize the appropriate genesis parameters whenever the feature is added

	return nil
}

func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	return &types.GenesisState{
		Params: k.GetParams(ctx),
		SendingState: types.SendingState{
			NextSeq: k.GetNextSequence(ctx),
		},
		NextRoleProposalId: k.GetNextProposalID(ctx),
		RoleMetadata:       k.GetRoleMetadata(ctx),
		RoleProposals:      k.GetProposals(ctx),
	}
}
