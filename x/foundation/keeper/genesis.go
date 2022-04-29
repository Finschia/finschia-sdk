package keeper

import (
	"time"

	sdk "github.com/line/lbm-sdk/types"

	"github.com/line/lbm-sdk/x/foundation"

	stakingtypes "github.com/line/lbm-sdk/x/staking/types"
)

func (k Keeper) InitGenesis(ctx sdk.Context, sk foundation.StakingKeeper, data *foundation.GenesisState) error {
	k.SetParams(ctx, data.Params)

	validatorAuths := data.ValidatorAuths
	if k.GetEnabled(ctx) && len(validatorAuths) == 0 {
		// Allowed validators must exist if the module is enabled,
		// so it should be the very first block of the chain.
		// We gather the information from staking module.
		sk.IterateValidators(ctx, func(_ int64, addr stakingtypes.ValidatorI) (stop bool) {
			auth := foundation.ValidatorAuth{
				OperatorAddress: addr.GetOperator().String(),
				CreationAllowed: true,
			}
			validatorAuths = append(validatorAuths, auth)
			return false
		})
	}

	for _, auth := range validatorAuths {
		if err := k.SetValidatorAuth(ctx, auth); err != nil {
			return err
		}
	}

	members := data.Members
	if len(members) == 0 {
		for _, auth := range validatorAuths {
			member := foundation.Member{
				Address:  sdk.ValAddress(auth.OperatorAddress).ToAccAddress().String(),
				Weight:   sdk.OneDec(),
				Metadata: "genesis member",
			}
			members = append(members, member)
		}
	}
	for _, member := range members {
		k.setMember(ctx, member)
	}

	info := data.Foundation
	if info == nil {
		info = &foundation.FoundationInfo{
			Version: 1,
		}
	}

	totalWeight := sdk.ZeroDec()
	for _, member := range members {
		totalWeight = totalWeight.Add(member.Weight)
	}
	info.TotalWeight = totalWeight

	if err := sdk.ValidateAccAddress(info.Operator); err != nil {
		info.Operator = k.GetAdmin(ctx).String()
	}

	if info.DecisionPolicy == nil ||
		info.GetDecisionPolicy() == nil ||
		info.GetDecisionPolicy().ValidateBasic() != nil ||
		info.GetDecisionPolicy().Validate(k.config) != nil {
		policy := foundation.ThresholdDecisionPolicy{
			Threshold: k.config.MinThreshold,
			Windows: &foundation.DecisionPolicyWindows{
				VotingPeriod: 24 * time.Hour,
			},
		}
		if err := info.SetDecisionPolicy(&policy); err != nil {
			return err
		}
	}

	k.setFoundationInfo(ctx, *info)

	k.setPreviousProposalID(ctx, data.PreviousProposalId)

	for _, proposal := range data.Proposals {
		k.setProposal(ctx, proposal)
		k.addProposalToVPEndQueue(ctx, proposal)
	}

	for _, vote := range data.Votes {
		k.setVote(ctx, vote)
	}

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
		ValidatorAuths:     k.GetValidatorAuths(ctx),
		Foundation:         &info,
		Members:            k.GetMembers(ctx),
		PreviousProposalId: k.getPreviousProposalID(ctx),
		Proposals:          proposals,
		Votes:              votes,
	}
}
