package keeper

import (
	sdk "github.com/line/lbm-sdk/types"
	"github.com/line/lbm-sdk/x/foundation"
	stakingtypes "github.com/line/lbm-sdk/x/staking/types"
	"github.com/line/lbm-sdk/x/stakingplus"
)

func (k Keeper) InitGenesis(ctx sdk.Context, sk foundation.StakingKeeper, data *foundation.GenesisState) error {
	params := data.Params
	if params == nil {
		params = foundation.DefaultParams()
	}
	k.SetParams(ctx, params)

	authorizations := data.Authorizations
	createValidatorGrantees := getCreateValidatorGrantees(authorizations)
	if k.GetEnabled(ctx) && len(createValidatorGrantees) == 0 {
		// Allowed validators must exist if the module is enabled,
		// so it should be the very first block of the chain.
		// We gather the information from staking module.
		sk.IterateValidators(ctx, func(_ int64, addr stakingtypes.ValidatorI) (stop bool) {
			grantee := sdk.AccAddress(addr.GetOperator())
			createValidatorGrantees = append(createValidatorGrantees, grantee)

			// add to authorizations
			authorization := &stakingplus.CreateValidatorAuthorization{
				ValidatorAddress: sdk.ValAddress(grantee).String(),
			}
			ga := foundation.GrantAuthorization{
				Grantee: grantee.String(),
			}
			if err := ga.SetAuthorization(authorization); err != nil {
				panic(err)
			}
			authorizations = append(authorizations, ga)

			return false
		})
	}

	members := data.Members
	if len(members) == 0 {
		for _, grantee := range createValidatorGrantees {
			member := foundation.Member{
				Address:  grantee.String(),
				Metadata: "genesis member",
			}
			members = append(members, member)
		}
	}
	for _, member := range members {
		if err := validateMetadata(member.Metadata, k.config); err != nil {
			return err
		}

		k.setMember(ctx, member)
	}

	info := data.Foundation
	if info == nil {
		info = &foundation.FoundationInfo{
			Version: 1,
		}
	}

	totalWeight := int64(len(members))
	info.TotalWeight = sdk.NewDec(totalWeight)

	if len(info.Operator) == 0 {
		info.Operator = k.GetAdmin(ctx).String()
	}

	if info.GetDecisionPolicy() == nil ||
		info.GetDecisionPolicy().ValidateBasic() != nil ||
		info.GetDecisionPolicy().Validate(*info, k.config) != nil {
		policy := foundation.DefaultDecisionPolicy()
		if err := info.SetDecisionPolicy(policy); err != nil {
			return err
		}
	}

	k.setFoundationInfo(ctx, *info)

	k.setPreviousProposalID(ctx, data.PreviousProposalId)

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

	for _, ga := range authorizations {
		grantee, err := sdk.AccAddressFromBech32(ga.Grantee)
		if err != nil {
			return err
		}
		k.setAuthorization(ctx, grantee, ga.GetAuthorization())
	}

	k.SetPool(ctx, data.Pool)

	k.SetOneTimeMintLeftCount(ctx, data.OneTimeMintLeftCount)

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
		Params:               k.GetParams(ctx),
		Foundation:           &info,
		Members:              k.GetMembers(ctx),
		PreviousProposalId:   k.getPreviousProposalID(ctx),
		Proposals:            proposals,
		Votes:                votes,
		Authorizations:       k.GetGrants(ctx),
		OneTimeMintLeftCount: k.GetOneTimeMintLeftCount(ctx),
	}
}

func (k Keeper) ResetState(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	// TODO: reset validator list too

	// reset foundation
	store.Delete(foundationInfoKey)

	// reset members
	for _, member := range k.GetMembers(ctx) {
		addr, err := sdk.AccAddressFromBech32(member.Address)
		if err != nil {
			panic(err)
		}
		store.Delete(memberKey(addr))
	}

	// id
	store.Delete(previousProposalIDKey)

	// reset proposals & votes
	for _, proposal := range k.GetProposals(ctx) {
		k.pruneProposal(ctx, proposal)
	}

	// reset authorizations
	for _, ga := range k.GetGrants(ctx) {
		grantee, err := sdk.AccAddressFromBech32(ga.Grantee)
		if err != nil {
			panic(err)
		}
		k.deleteAuthorization(ctx, grantee, ga.GetAuthorization().MsgTypeURL())
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

func getCreateValidatorGrantees(authorizations []foundation.GrantAuthorization) []sdk.AccAddress {
	msgTypeURL := stakingplus.CreateValidatorAuthorization{}.MsgTypeURL()
	var grantees []sdk.AccAddress
	for _, ga := range authorizations {
		if ga.GetAuthorization().MsgTypeURL() == msgTypeURL {
			grantee, err := sdk.AccAddressFromBech32(ga.Grantee)
			if err != nil {
				panic(err)
			}
			grantees = append(grantees, grantee)
		}
	}

	return grantees
}
