package keeper

import (
	sdk "github.com/line/lbm-sdk/types"

	"github.com/line/lbm-sdk/x/foundation"

	stakingtypes "github.com/line/lbm-sdk/x/staking/types"
)

func (k Keeper) InitGenesis(ctx sdk.Context, sk foundation.StakingKeeper, data *foundation.GenesisState) error {
	params := data.Params
	if params == nil {
		params = &foundation.Params{}
	}
	k.SetParams(ctx, params)

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
				Address:       sdk.ValAddress(auth.OperatorAddress).ToAccAddress().String(),
				Participating: true,
				Metadata:      "genesis member",
			}
			members = append(members, member)
		}
	}
	for _, member := range members {
		if err := validateMetadata(member.Metadata, k.config); err != nil {
			return err
		}

		if member.Participating {
			k.setMember(ctx, member)
		}
	}

	info := data.Foundation
	if info == nil {
		info = &foundation.FoundationInfo{
			Version: 1,
		}
	}

	totalWeight := int64(0)
	for _, member := range members {
		if member.Participating {
			totalWeight++
		}
	}
	info.TotalWeight = sdk.NewDec(totalWeight)

	if len(info.Operator) == 0 {
		info.Operator = k.GetAdmin(ctx).String()
	}

	if info.GetDecisionPolicy() == nil ||
		info.GetDecisionPolicy().ValidateBasic() != nil ||
		info.GetDecisionPolicy().Validate(k.config) != nil {
		policy := foundation.DefaultDecisionPolicy(k.config)
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

func (k Keeper) ResetState(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	// TODO: reset validator list too

	// reset foundation
	store.Delete(foundationInfoKey)

	// reset members
	for _, member := range k.GetMembers(ctx) {
		store.Delete(memberKey(sdk.AccAddress(member.Address)))
	}

	// id
	store.Delete(previousProposalIDKey)

	// reset proposals & votes
	for _, proposal := range k.GetProposals(ctx) {
		k.pruneProposal(ctx, proposal)
	}
}
