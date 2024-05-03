package keeper

import (
	"encoding/binary"
	sdk "github.com/Finschia/finschia-sdk/types"
	"github.com/Finschia/finschia-sdk/x/fbridge/types"
)

func (k Keeper) InitGenesis(ctx sdk.Context, gs *types.GenesisState) error {
	k.SetParams(ctx, gs.Params)
	k.setNextSequence(ctx, gs.SendingState.NextSeq)

	k.setNextProposalID(ctx, gs.NextRoleProposalId)
	for _, proposal := range gs.RoleProposals {
		k.setRoleProposal(ctx, proposal)
	}

	for _, vote := range gs.Votes {
		k.setVote(ctx, vote.ProposalId, sdk.MustAccAddressFromBech32(vote.Voter), vote.Option)
	}

	for _, pair := range gs.Roles {
		k.SetRole(ctx, pair.Role, sdk.MustAccAddressFromBech32(pair.Address))
	}
	k.SetRoleMetadata(ctx, gs.RoleMetadata)

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
		RoleProposals:      k.GetRoleProposals(ctx),
		Votes:              k.GetAllVotes(ctx),
		RoleMetadata:       k.GetRoleMetadata(ctx),
		Roles:              k.GetRolePairs(ctx),
	}
}

// IterateVotes iterates over the all the votes for role proposals and performs a callback function
func (k Keeper) IterateVotes(ctx sdk.Context, cb func(proposal types.Vote) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.KeyProposalVotePrefix)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		id, voter := types.SplitVoterVoteKey(iterator.Key())
		opt := types.VoteOption(binary.BigEndian.Uint32(iterator.Value()))
		v := types.Vote{ProposalId: id, Voter: voter.String(), Option: opt}
		if cb(v) {
			break
		}
	}
}

// GetAllVotes returns all the votes from the store
func (k Keeper) GetAllVotes(ctx sdk.Context) (votes []types.Vote) {
	k.IterateVotes(ctx, func(vote types.Vote) bool {
		votes = append(votes, vote)
		return false
	})
	return
}
