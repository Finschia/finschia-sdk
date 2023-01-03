package keeper

import (
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/foundation"
)

func (k Keeper) InitGenesis(ctx sdk.Context, data *foundation.GenesisState) error {
	k.SetParams(ctx, data.Params)

	k.SetFoundationInfo(ctx, data.Foundation)

	k.setPreviousProposalID(ctx, data.PreviousProposalId)

	for _, member := range data.Members {
		if err := validateMetadata(member.Metadata, k.config); err != nil {
			return err
		}

		k.SetMember(ctx, member)
	}

	for _, proposal := range data.Proposals {
		if err := validateMetadata(proposal.Metadata, k.config); err != nil {
			return err
		}

		k.setProposal(ctx, proposal)
		k.addProposalToVPEndQueue(ctx, proposal)
	}

	for _, vote := range data.Votes {
		if err := validateMetadata(vote.Metadata, k.config); err != nil {
			return err
		}

		k.setVote(ctx, vote)
	}

	for _, ga := range data.Authorizations {
		grantee := sdk.MustAccAddressFromBech32(ga.Grantee)
		k.setAuthorization(ctx, grantee, ga.GetAuthorization())
	}

	k.SetPool(ctx, data.Pool)

	return nil
}

func (k Keeper) ExportGenesis(ctx sdk.Context) *foundation.GenesisState {
	info := k.GetFoundationInfo(ctx)
	proposals := k.GetProposals(ctx)

	var votes []foundation.Vote
	for _, proposal := range proposals {
		votes = append(votes, k.GetVotes(ctx, proposal.Id)...)
	}

	return &foundation.GenesisState{
		Params:             k.GetParams(ctx),
		Foundation:         info,
		Members:            k.GetMembers(ctx),
		PreviousProposalId: k.getPreviousProposalID(ctx),
		Proposals:          proposals,
		Votes:              votes,
		Authorizations:     k.GetGrants(ctx),
		Pool:               k.GetPool(ctx),
	}
}

func (k Keeper) GetGrants(ctx sdk.Context) []foundation.GrantAuthorization {
	var grantAuthorizations []foundation.GrantAuthorization
	k.iterateAuthorizations(ctx, func(grantee sdk.AccAddress, authorization foundation.Authorization) (stop bool) {
		grantAuthorization := foundation.GrantAuthorization{
			Grantee: grantee.String(),
		}
		if err := grantAuthorization.SetAuthorization(authorization); err != nil {
			panic(err)
		}
		grantAuthorizations = append(grantAuthorizations, grantAuthorization)

		return false
	})
	return grantAuthorizations
}
