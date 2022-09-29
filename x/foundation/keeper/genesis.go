package keeper

import (
	sdk "github.com/line/lbm-sdk/types"

	"github.com/line/lbm-sdk/x/foundation"

	stakingtypes "github.com/line/lbm-sdk/x/staking/types"
	"github.com/line/lbm-sdk/x/stakingplus"
)

func (k Keeper) InitGenesis(ctx sdk.Context, sk foundation.StakingKeeper, data *foundation.GenesisState) error {
	k.SetParams(ctx, data.Params)

	authorizations := data.Authorizations
	createValidatorGrantees := getCreateValidatorGrantees(authorizations)
	isCreateValidatorCensored := k.IsCensoredMessage(ctx, sdk.MsgTypeURL((*stakingtypes.MsgCreateValidator)(nil)))
	if isCreateValidatorCensored && len(createValidatorGrantees) == 0 {
		// Allowed validators must exist if the `Msg/CreateValidator` is
		// being censored, or no validator would be created.
		// We gather all the operator addresses from staking module,
		// and allow them to create validators.
		sk.IterateValidators(ctx, func(_ int64, addr stakingtypes.ValidatorI) (stop bool) {
			grantee := sdk.AccAddress(addr.GetOperator())
			createValidatorGrantees = append(createValidatorGrantees, grantee)

			// add authorizations
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

	info := data.Foundation
	members := data.Members
	if len(members) == 0 {
		if _, outsourcing := info.GetDecisionPolicy().(*foundation.OutsourcingDecisionPolicy); !outsourcing {
			for _, grantee := range createValidatorGrantees {
				member := foundation.Member{
					Address:  grantee.String(),
					Metadata: "genesis member",
				}
				members = append(members, member)
			}
		}
	}
	for _, member := range members {
		if err := validateMetadata(member.Metadata, k.config); err != nil {
			return err
		}

		k.setMember(ctx, member)
	}

	totalWeight := int64(len(members))
	info.TotalWeight = sdk.NewDec(totalWeight)

	if len(info.Operator) == 0 {
		info.Operator = k.GetDefaultOperator(ctx).String()
	}
	k.setFoundationInfo(ctx, info)

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
	msgTypeURL := sdk.MsgTypeURL((*stakingtypes.MsgCreateValidator)(nil))
	var grantees []sdk.AccAddress
	for _, ga := range authorizations {
		if ga.GetAuthorization().MsgTypeURL() == msgTypeURL {
			grantee := sdk.MustAccAddressFromBech32(ga.Grantee)
			grantees = append(grantees, grantee)
		}
	}

	return grantees
}
